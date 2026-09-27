package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/ChristianDenniss/go-data-model/auth/entity"
	"github.com/ChristianDenniss/go-data-model/auth/repository"
)

const (
	SessionTTL        = 30 * 24 * time.Hour
	MinPasswordLength = 8
	// bcrypt ignores everything past 72 bytes, so longer passwords are rejected
	// rather than silently truncated.
	MaxPasswordBytes = 72
)

// PasswordHasher is supplied by the transport so the domain stays dependency-free.
type PasswordHasher interface {
	Hash(password string) (string, error)
	// Compare returns nil only when password matches hash.
	Compare(hash, password string) error
}

type Service struct {
	repo   repository.Repository
	hasher PasswordHasher
	now    func() time.Time

	dummyOnce sync.Once
	dummyHash string
}

func New(repo repository.Repository, hasher PasswordHasher) *Service {
	return &Service{repo: repo, hasher: hasher, now: time.Now}
}

func (s *Service) SignUp(ctx context.Context, name, email, password string) (entity.SignIn, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return entity.SignIn{}, entity.ErrNameRequired
	}
	email, err := normalizeEmail(email)
	if err != nil {
		return entity.SignIn{}, err
	}
	if err := validatePassword(password); err != nil {
		return entity.SignIn{}, err
	}
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return entity.SignIn{}, err
	}
	accountID, err := newAccountID()
	if err != nil {
		return entity.SignIn{}, err
	}
	if err := s.repo.CreateAccount(ctx, entity.NewAccount{
		ID:           accountID,
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}); err != nil {
		return entity.SignIn{}, err
	}
	return s.issueSession(ctx, accountID)
}

func (s *Service) Login(ctx context.Context, email, password string) (entity.SignIn, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return entity.SignIn{}, entity.ErrInvalidCredentials
	}
	creds, err := s.repo.CredentialsByEmail(ctx, email)
	if errors.Is(err, entity.ErrAccountNotFound) {
		// Spend the same time as a real comparison so response timing does not
		// reveal which emails are registered.
		_ = s.hasher.Compare(s.fallbackHash(), password)
		return entity.SignIn{}, entity.ErrInvalidCredentials
	}
	if err != nil {
		return entity.SignIn{}, err
	}
	if creds.PasswordHash == "" {
		_ = s.hasher.Compare(s.fallbackHash(), password)
		return entity.SignIn{}, entity.ErrInvalidCredentials
	}
	if err := s.hasher.Compare(creds.PasswordHash, password); err != nil {
		return entity.SignIn{}, entity.ErrInvalidCredentials
	}
	return s.issueSession(ctx, creds.AccountID)
}

// SignInWithProvider finds the account linked to the provider subject, links an
// existing account with the same verified email, or creates a new account.
func (s *Service) SignInWithProvider(ctx context.Context, profile entity.ExternalProfile) (entity.SignIn, error) {
	if profile.Provider == "" || profile.Subject == "" {
		return entity.SignIn{}, entity.ErrProviderRequired
	}

	identity, err := s.repo.IdentityBySubject(ctx, profile.Provider, profile.Subject)
	if err == nil {
		return s.issueSession(ctx, identity.AccountID)
	}
	if !errors.Is(err, entity.ErrIdentityNotFound) {
		return entity.SignIn{}, err
	}

	email, err := normalizeEmail(profile.Email)
	if err != nil {
		return entity.SignIn{}, err
	}
	if !profile.EmailVerified {
		return entity.SignIn{}, entity.ErrEmailUnverified
	}
	link := entity.Identity{Provider: profile.Provider, Subject: profile.Subject, Email: email}

	creds, err := s.repo.CredentialsByEmail(ctx, email)
	switch {
	case err == nil:
		link.AccountID = creds.AccountID
		if err := s.repo.LinkIdentity(ctx, link); err != nil {
			return entity.SignIn{}, err
		}
		return s.issueSession(ctx, creds.AccountID)
	case !errors.Is(err, entity.ErrAccountNotFound):
		return entity.SignIn{}, err
	}

	accountID, err := newAccountID()
	if err != nil {
		return entity.SignIn{}, err
	}
	link.AccountID = accountID
	name := strings.TrimSpace(profile.Name)
	if name == "" {
		name = strings.SplitN(email, "@", 2)[0]
	}
	if err := s.repo.CreateAccount(ctx, entity.NewAccount{
		ID:       accountID,
		Name:     name,
		Email:    email,
		Identity: &link,
	}); err != nil {
		return entity.SignIn{}, err
	}
	return s.issueSession(ctx, accountID)
}

// Authenticate resolves a raw session token to its account.
func (s *Service) Authenticate(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", entity.ErrSessionNotFound
	}
	session, err := s.repo.SessionByTokenHash(ctx, hashToken(token), s.now())
	if err != nil {
		return "", err
	}
	return session.AccountID, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.repo.DeleteSession(ctx, hashToken(token))
}

func (s *Service) issueSession(ctx context.Context, accountID string) (entity.SignIn, error) {
	token, err := randomToken(32)
	if err != nil {
		return entity.SignIn{}, err
	}
	expiresAt := s.now().Add(SessionTTL)
	if err := s.repo.CreateSession(ctx, entity.Session{
		TokenHash: hashToken(token),
		AccountID: accountID,
		ExpiresAt: expiresAt,
	}); err != nil {
		return entity.SignIn{}, err
	}
	return entity.SignIn{AccountID: accountID, Token: token, ExpiresAt: expiresAt}, nil
}

func (s *Service) fallbackHash() string {
	s.dummyOnce.Do(func() {
		s.dummyHash, _ = s.hasher.Hash("nibble-timing-equalizer")
	})
	return s.dummyHash
}

func normalizeEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" {
		return "", entity.ErrEmailRequired
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", entity.ErrEmailInvalid
	}
	return email, nil
}

func validatePassword(password string) error {
	if len([]rune(password)) < MinPasswordLength {
		return entity.ErrPasswordTooShort
	}
	if len(password) > MaxPasswordBytes {
		return entity.ErrPasswordTooLong
	}
	return nil
}

func newAccountID() (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "acct_" + hex.EncodeToString(buf), nil
}

func randomToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
