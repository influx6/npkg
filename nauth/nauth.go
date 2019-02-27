package nauth

import (
	"errors"
	"net/http"

	"github.com/gokit/npkg/nxid"
)

// ErrNoCredentials is returned when giving claim fails to provide
// a credential.
var ErrNoCredentials = errors.New("Claim has no attached credentail")

// Credential defines what we expect from a custom implemented
// credential.
type Credential interface {
	Validate() error

	Type() string
	User() string
	Provider() string // google, email, phone, facebook, wechat, github, ...
}

// ClaimProvider defines what we expect from a Claim provider.
type ClaimProvider interface {
	EncodeClaim(Claim) ([]byte, error)
	ParseClaim(claim []byte) (Claim, error)
	ExtractClaim(r *http.Request) (Claim, error)
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
	Method string // jwt, user-password, oauth, ...
	Cred   Credential
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
// claims.
type Authenticator interface {
	Authenticate(Claim) (VerifiedClaim, error)
}

// VerifiedClaim represents the response received back from the
// Authenticator as to a giving authenticated claim with associated
// session data.
type VerifiedClaim struct {
	User  nxid.ID
	Roles []string               // Roles of verified claim.
	Data  map[string]interface{} // Extra Data to be attached to session for user.
}

// Valid returns an error if giving credentials could not be validated
// or if giving Claim has no attached credential.
func (c VerifiedClaim) Valid() error {
	return nil
}

// HasRoles returns true/false if giving claim as all roles.
func (c VerifiedClaim) HasRoles(roles ...string) bool {
	for _, role := range roles {
		if c.checkRole(role) {
			continue
		}
		return false
	}
	return true
}

// HasAnyRoles returns true if giving claim as at least one roles.
func (c VerifiedClaim) HasAnyRoles(roles ...string) bool {
	for _, role := range roles {
		if c.checkRole(role) {
			return true
		}
	}
	return false
}

// checkRole checks if any roles of Claim match provided.
func (c VerifiedClaim) checkRole(role string) bool {
	for _, myrole := range c.Roles {
		if myrole == role {
			return true
		}
	}
	return false
}
