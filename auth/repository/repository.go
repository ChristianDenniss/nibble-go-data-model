package repository

import (
	"context"
	"time"

	"github.com/ChristianDenniss/go-data-model/auth/entity"
)

type Repository interface {
	// CreateAccount inserts the account (and its identity, when set) atomically.
	// Returns entity.ErrEmailTaken when the email is already registered.
	CreateAccount(ctx context.Context, account entity.NewAccount) error
	// CredentialsByEmail matches email case-insensitively. Returns entity.ErrAccountNotFound.
	CredentialsByEmail(ctx context.Context, email string) (entity.Credentials, error)
	// IdentityBySubject returns entity.ErrIdentityNotFound.
	IdentityBySubject(ctx context.Context, provider, subject string) (entity.Identity, error)
	LinkIdentity(ctx context.Context, identity entity.Identity) error
	CreateSession(ctx context.Context, session entity.Session) error
	// SessionByTokenHash returns entity.ErrSessionNotFound when missing or expired at now.
	SessionByTokenHash(ctx context.Context, tokenHash string, now time.Time) (entity.Session, error)
	DeleteSession(ctx context.Context, tokenHash string) error
}
