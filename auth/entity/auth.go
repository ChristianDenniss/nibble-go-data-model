package entity

import "time"

const (
	ProviderGoogle = "google"
	ProviderApple  = "apple"
)

// Credentials is the password login record for an account. PasswordHash is
// empty for accounts that only ever signed in through a provider.
type Credentials struct {
	AccountID    string
	PasswordHash string
}

// Identity links a provider subject (Google `sub`, Apple `sub`) to an account.
type Identity struct {
	Provider  string
	Subject   string
	AccountID string
	Email     string
}

// Session is one signed-in browser. Only the SHA-256 of the cookie token is stored.
type Session struct {
	TokenHash string
	AccountID string
	ExpiresAt time.Time
}

// NewAccount is the profile captured when an account is created. Identity is
// set when the account is created by a provider sign-in.
type NewAccount struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Identity     *Identity
}

// ExternalProfile is what a provider vouched for after a successful OAuth exchange.
type ExternalProfile struct {
	Provider      string
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
}

// SignIn is a freshly issued session. Token is the raw secret for the cookie
// and is never persisted.
type SignIn struct {
	AccountID string
	Token     string
	ExpiresAt time.Time
}
