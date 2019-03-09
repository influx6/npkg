package emailauth

import (
	"net/http"

	"github.com/gokit/npkg/nauth"
	"github.com/gokit/npkg/nauth/sessions"
	"github.com/gokit/npkg/nerror"
	"github.com/gokit/npkg/nxid"
)

//***************************************************************************
// In-House Email Credential
//***************************************************************************

// EmailCredential defines the underline data expected for
// the authentication of a user session.
type EmailCredential struct {
	Email    string `json:"email"`
	Password string `json:"Password"`
	Username string `json:"username" optional:"false"`
}

// Validate returns an error if credential is invalid.
func (e EmailCredential) Validate() error {
	if e.Email == "" && e.Username != "" {
		return nerror.New("EmailCredential must have Email or Username")
	}
	if e.Password == "" {
		return nerror.New("EmailCredential.Password can not be empty")
	}
	return nil
}

// User returns giving user credential, either the email if
// available, else the username.
func (e EmailCredential) User() string {
	if e.Email != "" {
		return e.Email
	}
	return e.Username
}

//***************************************************************************
// In-House Email Authentication Provider
//***************************************************************************

const (
	// ProviderName sets the provider name for giving email authentication package.
	ProviderName = "in-house"

	// MethodName sets the expected method for authentication.
	MethodName = "email-password"
)

// VerifiedEmail contains the verified token delivered to by a verified client login.
type VerifiedEmail struct {
	Token    string `json:"token"`
	PublicID string `json:"public_id"`
}

// UserData contains data retrieved from underline UserStore containing
// appropriate data useful for authentication and authorization.
type UserData struct {
	Email          string
	Username       string
	PrivateSalt    string
	HashedPassword string
	PublicID       nxid.ID
	Roles          []string
}

// UserStore defines what we expect from a user store implementation,
// it exposes the methods necessary to retrieve roles for a giving
// user.
type UserStore interface {
	GetEmail(email string) (UserData, error)
	GetPublicID(id string) (UserData, error)
	GetUsername(username string) (UserData, error)
}

// UserValidator embodies what we expect to use for verifying user data and
// credentials. It returns some data it want attached to the verified user data.
type UserValidator interface {
	Verify(credential EmailCredential, data UserData) (interface{}, error)
}

var _ nauth.AuthenticationProvider = (*InhouseEmailAuth)(nil)

// InhouseEmailAuth provides an implementation of a AuthenticationProvider,
// which provides email and password authentication and authorization.
type InhouseEmailAuth struct {
	UserStore     UserStore
	UserValidator UserValidator
	AuthInitiator http.Handler
	Sessions      sessions.Sessions
}

// Initiate implements the nauth.AuthenticationProvider interface.
//
// Initiate runs the provided User http.handler which handles the rendering of
// necessary login form for the user to input defined credentials for authentication.
func (eu InhouseEmailAuth) Initiate(res http.ResponseWriter, req *http.Request) {
	if eu.AuthInitiator != nil {
		eu.AuthInitiator.ServeHTTP(res, req)
		return
	}
	res.WriteHeader(http.StatusNotImplemented)
}

// VerifyClaim implements the nauth.Authenticator interface.
func (eu InhouseEmailAuth) VerifyClaim(cm nauth.Claim) (nauth.VerifiedClaim, error) {
	var verified nauth.VerifiedClaim
	if cm.Method != MethodName {
		return verified, nerror.New("claim.Method does not matched supported")
	}
	if cm.Provider != ProviderName {
		return verified, nerror.New("claim.Provider does not matched supported")
	}

	verified.Method = cm.Method
	verified.Provider = cm.Provider

	var err error
	if err = cm.Credentials.Validate(); err != nil {
		return verified, nerror.Wrap(err, "claim credentials are invalid")
	}

	var credential, ok = cm.Credentials.(EmailCredential)
	if !ok {
		return verified, nerror.New("claim has unsupported/invalid credentials")
	}

	var userdata UserData
	if credential.Email != "" {
		userdata, err = eu.UserStore.GetEmail(credential.Email)
	}
	if credential.Username != "" {
		userdata, err = eu.UserStore.GetUsername(credential.Username)
	}

	if err != nil {
		return verified, err
	}

	var attached interface{}
	attached, err = eu.UserValidator.Verify(credential, userdata)
	if err != nil {
		return verified, err
	}

	verified.Data = attached
	verified.Roles = userdata.Roles
	verified.User = userdata.PublicID
	return verified, nil
}

// Verify implements the nauth.AuthenticationProvider interface.
//
// Verify exists for the purpose of verifying  an authenticated session with
// an existing bearer token.
func (eu InhouseEmailAuth) Verify(req *http.Request) (nauth.VerifiedClaim, error) {
	var verified nauth.VerifiedClaim

	// Retrieve Authorization Header from request.
	var authorizationHeader = req.Header.Get(nauth.AuthorizationHeaderName)
	var authType, authToken, err = nauth.ParseAuthorization(authorizationHeader)
	if err != nil {
		return verified, err
	}

	_ = authorizationHeader
	_ = authType
	_ = authToken
	panic("implement me")
}

// Refresh implements the nauth.AuthenticationProvider interface.
// It refreshes the authentication session, by updating the session expiring
// for giving user. This allows you to set a default session deadline and have
// your UI implement some inactivity checker which calls this to update session
// expiry and allow user have longer access to your site.
func (eu InhouseEmailAuth) Refresh(res http.ResponseWriter, req *http.Request) {
	panic("implement me")
}

// Authenticate implements the nauth.AuthenticationProvider interface.
func (eu InhouseEmailAuth) Authenticate(res http.ResponseWriter, req *http.Request) {
	panic("implement me")
}
