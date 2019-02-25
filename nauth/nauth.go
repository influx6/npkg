package nauth

import (
	"errors"
)

// ErrNoCredentials is returned when giving claim fails to provide
// a credential.
var ErrNoCredentials = errors.New("Claim has no attached credentail")

// Credential defines what we expect from a custom implemented
// credential.
type Credential interface {
	Type() string
	Validate() error
}

// ClaimProvider defines what we expect from a Claim provider.
type ClaimProvider interface {
	EncodeClaim(Claim) ([]byte, error)
	ParseClaim(claim []byte) (Claim, error)
}

// Claims define what we expect from a Claim implementation.
type Claims interface {
	Valid() error
	HasRoles(...string) bool
	HasAnyRoles(...string) bool
}

// Claim defines authentication claims parsed from underline
// data provide to authenticator.
type Claim struct {
	Method string
	Roles  []string
	Cred   Credential
}

// HasRoles returns true/false if giving claim as all roles.
func (c Claim) HasRoles(roles ...string) bool {
	for _, role := range roles {
		if c.checkRole(role) {
			continue
		}
		return false
	}
	return true
}

// HasAnyRoles returns true if giving claim as at least one roles.
func (c Claim) HasAnyRoles(roles ...string) bool {
	for _, role := range roles {
		if c.checkRole(role) {
			return true
		}
	}
	return false
}

// checkRole checks if any roles of Claim match provided.
func (c Claim) checkRole(role string) bool {
	for _, myrole := range c.Roles {
		if myrole == role {
			return true
		}
	}
	return false
}

// Valid returns an error if giving credentials could not be validated
// or if giving Claim has no attached credential.
func (c Claim) Valid() error {
	if c.Cred != nil {
		return c.Cred.Validate()
	}
	return ErrNoCredentials
}

// Authenticator defines what we expect from a Authenticator of
// claims/
type Authenticator interface {
	Authenticate(Claim) error
}
