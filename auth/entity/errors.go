package entity

import "errors"

var (
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrEmailInvalid       = errors.New("email is not valid")
	ErrEmailUnverified    = errors.New("provider did not verify this email")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong    = errors.New("password must be at most 72 bytes")
	ErrEmailTaken         = errors.New("an account with this email already exists")
	ErrInvalidCredentials = errors.New("email or password is incorrect")
	ErrAccountNotFound    = errors.New("account not found")
	ErrIdentityNotFound   = errors.New("identity not found")
	ErrProviderRequired   = errors.New("provider and subject are required")
	ErrSessionNotFound    = errors.New("session not found or expired")
)
