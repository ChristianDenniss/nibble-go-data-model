package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ChristianDenniss/go-data-model/auth/entity"
)

type plainHasher struct{}

func (plainHasher) Hash(password string) (string, error) { return "h:" + password, nil }

func (plainHasher) Compare(hash, password string) error {
	if hash != "h:"+password {
		return errors.New("mismatch")
	}
	return nil
}

type memRepo struct {
	accounts   map[string]entity.NewAccount
	identities map[string]entity.Identity
	sessions   map[string]entity.Session
}

func newMemRepo() *memRepo {
	return &memRepo{
		accounts:   map[string]entity.NewAccount{},
		identities: map[string]entity.Identity{},
		sessions:   map[string]entity.Session{},
	}
}

func (m *memRepo) CreateAccount(_ context.Context, a entity.NewAccount) error {
	for _, existing := range m.accounts {
		if existing.Email == a.Email {
			return entity.ErrEmailTaken
		}
	}
	m.accounts[a.ID] = a
	if a.Identity != nil {
		m.identities[a.Identity.Provider+"/"+a.Identity.Subject] = *a.Identity
	}
	return nil
}

func (m *memRepo) CredentialsByEmail(_ context.Context, email string) (entity.Credentials, error) {
	for _, a := range m.accounts {
		if a.Email == email {
			return entity.Credentials{AccountID: a.ID, PasswordHash: a.PasswordHash}, nil
		}
	}
	return entity.Credentials{}, entity.ErrAccountNotFound
}

func (m *memRepo) IdentityBySubject(_ context.Context, provider, subject string) (entity.Identity, error) {
	identity, ok := m.identities[provider+"/"+subject]
	if !ok {
		return entity.Identity{}, entity.ErrIdentityNotFound
	}
	return identity, nil
}

func (m *memRepo) LinkIdentity(_ context.Context, identity entity.Identity) error {
	m.identities[identity.Provider+"/"+identity.Subject] = identity
	return nil
}

func (m *memRepo) CreateSession(_ context.Context, session entity.Session) error {
	m.sessions[session.TokenHash] = session
	return nil
}

func (m *memRepo) SessionByTokenHash(_ context.Context, hash string, now time.Time) (entity.Session, error) {
	session, ok := m.sessions[hash]
	if !ok || !session.ExpiresAt.After(now) {
		return entity.Session{}, entity.ErrSessionNotFound
	}
	return session, nil
}

func (m *memRepo) DeleteSession(_ context.Context, hash string) error {
	delete(m.sessions, hash)
	return nil
}

func TestSignUpLoginLogout(t *testing.T) {
	ctx := context.Background()
	svc := New(newMemRepo(), plainHasher{})

	signed, err := svc.SignUp(ctx, " Alex ", "Alex@Example.com ", "correct-horse")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	if got, err := svc.Authenticate(ctx, signed.Token); err != nil || got != signed.AccountID {
		t.Fatalf("authenticate after signup = %q, %v", got, err)
	}

	if _, err := svc.SignUp(ctx, "Other", "alex@example.com", "another-pass"); !errors.Is(err, entity.ErrEmailTaken) {
		t.Fatalf("duplicate signup err = %v, want ErrEmailTaken", err)
	}
	if _, err := svc.Login(ctx, "alex@example.com", "wrong-pass"); !errors.Is(err, entity.ErrInvalidCredentials) {
		t.Fatalf("wrong password err = %v", err)
	}
	if _, err := svc.Login(ctx, "nobody@example.com", "correct-horse"); !errors.Is(err, entity.ErrInvalidCredentials) {
		t.Fatalf("unknown email err = %v", err)
	}

	loggedIn, err := svc.Login(ctx, "ALEX@example.com", "correct-horse")
	if err != nil || loggedIn.AccountID != signed.AccountID {
		t.Fatalf("login = %+v, %v", loggedIn, err)
	}
	if err := svc.Logout(ctx, loggedIn.Token); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if _, err := svc.Authenticate(ctx, loggedIn.Token); !errors.Is(err, entity.ErrSessionNotFound) {
		t.Fatalf("authenticate after logout err = %v", err)
	}
}

func TestSignUpValidation(t *testing.T) {
	svc := New(newMemRepo(), plainHasher{})
	cases := []struct {
		name, email, password string
		want                  error
	}{
		{"", "a@b.co", "longenough", entity.ErrNameRequired},
		{"A", "", "longenough", entity.ErrEmailRequired},
		{"A", "not-an-email", "longenough", entity.ErrEmailInvalid},
		{"A", "a@b.co", "short", entity.ErrPasswordTooShort},
	}
	for _, tc := range cases {
		if _, err := svc.SignUp(context.Background(), tc.name, tc.email, tc.password); !errors.Is(err, tc.want) {
			t.Errorf("SignUp(%q, %q) err = %v, want %v", tc.name, tc.email, err, tc.want)
		}
	}
}

func TestSessionExpires(t *testing.T) {
	ctx := context.Background()
	svc := New(newMemRepo(), plainHasher{})
	signed, err := svc.SignUp(ctx, "Alex", "alex@example.com", "correct-horse")
	if err != nil {
		t.Fatal(err)
	}
	svc.now = func() time.Time { return time.Now().Add(SessionTTL + time.Minute) }
	if _, err := svc.Authenticate(ctx, signed.Token); !errors.Is(err, entity.ErrSessionNotFound) {
		t.Fatalf("expired session err = %v", err)
	}
}

func TestSignInWithProvider(t *testing.T) {
	ctx := context.Background()
	repo := newMemRepo()
	svc := New(repo, plainHasher{})

	existing, err := svc.SignUp(ctx, "Alex", "alex@example.com", "correct-horse")
	if err != nil {
		t.Fatal(err)
	}

	google := entity.ExternalProfile{
		Provider: entity.ProviderGoogle, Subject: "g-1", Email: "Alex@example.com", EmailVerified: true,
	}
	linked, err := svc.SignInWithProvider(ctx, google)
	if err != nil || linked.AccountID != existing.AccountID {
		t.Fatalf("link by verified email = %+v, %v", linked, err)
	}
	again, err := svc.SignInWithProvider(ctx, entity.ExternalProfile{Provider: entity.ProviderGoogle, Subject: "g-1"})
	if err != nil || again.AccountID != existing.AccountID {
		t.Fatalf("returning identity = %+v, %v", again, err)
	}

	fresh, err := svc.SignInWithProvider(ctx, entity.ExternalProfile{
		Provider: entity.ProviderApple, Subject: "a-1", Email: "relay@privaterelay.appleid.com", EmailVerified: true,
	})
	if err != nil || fresh.AccountID == existing.AccountID {
		t.Fatalf("new provider account = %+v, %v", fresh, err)
	}
	if got := repo.accounts[fresh.AccountID].Name; got != "relay" {
		t.Fatalf("fallback name = %q, want relay", got)
	}
	if _, err := svc.Login(ctx, "relay@privaterelay.appleid.com", "anything-long"); !errors.Is(err, entity.ErrInvalidCredentials) {
		t.Fatalf("password login on provider-only account err = %v", err)
	}

	_, err = svc.SignInWithProvider(ctx, entity.ExternalProfile{
		Provider: entity.ProviderGoogle, Subject: "g-2", Email: "alex@example.com", EmailVerified: false,
	})
	if !errors.Is(err, entity.ErrEmailUnverified) {
		t.Fatalf("unverified email err = %v, want ErrEmailUnverified", err)
	}
}
