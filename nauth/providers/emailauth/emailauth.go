package emailauth

import (
	"github.com/influx6/npkg/nauth/providers"
	"net/http"

	"github.com/influx6/npkg/nauth"
	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nxid"
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

var _ nauth.AuthenticationProvider = (*InhouseEmailAuth)(nil)

// VerifiedEmail contains the verified token delivered to by a verified client login.
type VerifiedEmail struct {
	Token    string `json:"token"`
	PublicID string `json:"public_id"`
}

type UserBaseData struct {
	Email          string
	Username       string
	PublicID       nxid.ID
}

func (u UserBaseData) Type() string {
	return "EmailAuth_UserData"
}

func (u UserBaseData) WebSafe() map[string]string {
	return map[string]string{
		"email": u.Email,
		"username": u.Username,
		"public_id": u.PublicID.String(),
	}
}

func (u *UserBaseData) FromMap(data map[string]string) error {
	var userName, hasUserName = data["username"]
	var userEmail, hasUserEmail = data["email"]
	if !hasUserEmail && !hasUserName {
		return nerror.New("data has no username or email attribute")
	}
	if userEmail == "" && userName == "" {
		return nerror.New("data has no username and email attribute is empty")
	}

	u.Username = userName
	u.Email = userEmail

	var publicId, hasPublicID = data["public_id"]
	if !hasPublicID {
		return nerror.New("data has no public_id attribute")
	}
	if publicId == "" {
		return nerror.New("public_id value in data is empty")
	}

	var id, err = nxid.FromString(publicId)
	if err != nil {
		return nerror.Wrap(err, "Failed to convert public_id, its invalid")
	}

	u.PublicID = id
	return nil
}

type UserData struct {
	UserBaseData

	PrivateSalt    string
	HashedPassword string
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
	Verify(credential EmailCredential, data UserData) error
}

// InhouseEmailAuth provides an implementation of a AuthenticationProvider,
// which provides email and password authentication and authorization.
type InhouseEmailAuth struct {
	UserStore     UserStore
	UserValidator UserValidator
	AuthInitiator http.Handler
	Sessions      providers.HTTPSession
}

func (eu InhouseEmailAuth) Finalize(req *http.Request) error {
	panic("implement me")
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

	if credentialErr := cm.Credentials.Validate(); credentialErr != nil {
		return verified, nerror.Wrap(credentialErr, "claim credentials are invalid")
	}

	var credential, ok = cm.Credentials.(EmailCredential)
	if !ok {
		return verified, nerror.New("claim has unsupported/invalid credentials")
	}

	var userIdentifier string
	if credential.Email != "" {
		userIdentifier = credential.Email
	}
	if credential.Username != "" {
		userIdentifier = credential.Username
	}

	var userdata, getUserErr = eu.UserStore.GetEmail(userIdentifier)
	if getUserErr != nil {
		return verified, nerror.Wrap(getUserErr, "Failed to find user")
	}

	var err = eu.UserValidator.Verify(credential, userdata)
	if err != nil {
		return verified, err
	}

	verified.Attached = userdata
	verified.Roles = userdata.Roles
	verified.User = userdata.PublicID
	return verified, nil
}

// GetVerifiedClaim implements the nauth.AuthenticationProvider interface.
//
// Verify exists for the purpose of verifying  an authenticated session with
// an existing bearer token.
func (eu InhouseEmailAuth) GetVerifiedClaim(req *http.Request) (nauth.VerifiedClaim, error) {
	var verified nauth.VerifiedClaim

	// Retrieve user session from request.
	var session, err = eu.Sessions.Get(req)
	if err != nil {
		return verified, nerror.Wrap(err, "http.Request has no existing auth session")
	}

	var userData UserBaseData
	if userErr := userData.FromMap(session.Data); userErr != nil {
		return verified, nerror.Wrap(userErr, "Failed to convert user session data to UserBaseData")
	}

	verified.Attached = userData
	verified.Method = session.Method
	verified.Provider = session.Provider
	return verified, nil
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
