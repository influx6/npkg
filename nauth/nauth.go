package nauth

import (
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/influx6/npkg/nxid"
)

const (
	// AuthorizationHeaderName defines the giving header name for retrieving
	// a authorization token for a authentication user.
	AuthorizationHeaderName = "Authorization"
)

// ErrNoCredentials is returned when giving claim fails to provide
// a credential.
var ErrNoCredentials = errors.New("Claim has no attached credentail")

// Roles exposes an interface to retrieve roles information for
// a giving identity id.
type Roles interface {
	Get(id nxid.ID) ([]string, error)
}

// Credential defines what we expect from a custom implemented
// credential.
type Credential interface {
	User() string
	Validate() error
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
	// Method sets the defined credential authentication type being used.
	Method string // email-password, phone-number, token,..etc

	// Provider defines the provider of authentication, providing adequate information
	// as towards the source.
	Provider string // google, in-house, phone, facebook, we-chat, github, ...etc

	// Attached IP during reception of claim, possibly attached by the handler receiving
	// giving request which maybe able to get IP details.
	IP net.IP

	// Attached Agent during reception of claim, possibly attached by the handler receiving
	// giving request which maybe able to get Agent details.
	Agent string

	// Credentials contains the deserialized data delivered by the user for authentication
	// which must match the method and provider expected data type.
	Credentials Credential
}

// Valid returns an error if giving credentials could not be validated
// or if giving Claim has no attached credential.
func (c Claim) Valid() error {
	if c.Credentials != nil {
		return c.Credentials.Validate()
	}
	return ErrNoCredentials
}

// VerifiedClaim represents the response received back from the
// Authenticator as to a giving authenticated claim with associated
// session data.
type VerifiedClaim struct {
	User         nxid.ID
	BrowserAgent string
	IP           net.IP
	Method       string      // email-password, phone-number, token,..etc
	Provider     string      // google, in-house, phone, facebook, we-chat, github, ...etc
	Roles        []string    // Roles of verified claim.
	Data         interface{} // Extra Data to be attached to session for user.
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

// Authenticator defines what we expect from a Authenticator of
// claims. It exposes the underline method used for verifying
// an authentication claim.
type Authenticator interface {
	// VerifyClaim exposes the underline function within Authenticator.Authenticate
	// used to authenticate the request claim and the returned verified claim. It
	// allows testing and also
	VerifyClaim(Claim) (VerifiedClaim, error)
}

// AuthenticationProvider defines what the Authentication should be as,
// it both exposes the the method from Authenticator and the provides
// the Initiate and Authenticate methods which are the underline
// handlers of the initiation and finalization of requests to authenticate.
//
// Exposes such a final form allows us to swap in, any form of authentication
// be it email, facebook, google or oauth based without much work.
type AuthenticationProvider interface {
	Authenticator

	// Initiate handles the initial response to a request to initiate/begin
	// a authentication procedure e.g to redirect to
	// a page for user-name and password login or google oauth page with
	// a secure token.
	Initiate(res http.ResponseWriter, req *http.Request)

	// Authenticate finalizes the response to initiation of authentication
	// with the call to AuthenticationProvider.Initiate.
	//
	// It handles the process which finalizes and verifies the authentication data sent
	// back after the initiation, with a response as dictated by provider.
	//
	// The authenticate process can be the authentication of a new login
	// or the authentication of an existing login. The provider implementation
	// should decide for it'self as it sees fit to match on how this two should
	// be managed.
	Authenticate(res http.ResponseWriter, req *http.Request)

	// Verify exposes to others by the provider a means of getting a verified
	// claim from a incoming request after it has being authenticated in some previous step.
	//
	// It exists to let you handle cases of already authenticated users whoes session is yet
	// to expire and are making new request for resources.
	//
	// This lets others step into the middle of the Authentication procedure
	// to retrieve the verified request claim as dictated by provider, which
	// can be used for other uses.
	Verify(req *http.Request) (VerifiedClaim, error)

	// Refresh handles the refreshing of an authentication session, useful
	// for protocols that require and provide refresh token as a means of
	// updating their access token expiry timeline.
	// This is based on protocols and a protocol may not implement it
	// and hence return a 501 (NOT Implemented) status
	Refresh(res http.ResponseWriter, req *http.Request)
}

// ParseAuthorization returns the scheme and token of the Authorization string
// if it's valid.
func ParseAuthorization(val string) (authType string, token string, err error) {
	authSplit := strings.SplitN(val, " ", 2)
	if len(authSplit) != 2 {
		err = errors.New("invalid authorization: Expected content: `AuthType Token`")
		return
	}

	authType = strings.TrimSpace(authSplit[0])
	token = strings.TrimSpace(authSplit[1])
	return
}

// ParseTokens parses the base64 encoded token sent as part of the Authorization string,
// It expects all parts of string to be seperated with ':', returning splitted slice.
func ParseTokens(val string) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(decoded), ":"), nil
}
