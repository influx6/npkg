package emailauth

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/influx6/npkg/nhttp"

	"github.com/gorilla/securecookie"

	"github.com/influx6/npkg/ntrace"
	openTracing "github.com/opentracing/opentracing-go"

	"github.com/influx6/npkg"
	"github.com/influx6/npkg/njson"

	"github.com/influx6/npkg/nauth/providers"
	"github.com/influx6/npkg/nauth/sessions"

	"github.com/influx6/npkg/nauth"
	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nxid"
)

const (
	// CLAIM_TYPE is the type of claim authentication implemented by emailauth.
	CLAIM_TYPE = "email_username_password"
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

var _ nauth.AuthenticationProvider = (*EmailAuth)(nil)

// VerifiedEmail contains the verified token delivered to by a verified client login.
type VerifiedEmail struct {
	Token    string `json:"token"`
	PublicID string `json:"public_id"`
}

type UserBaseData struct {
	Email    string
	Username string
	PublicID nxid.ID
}

func (u UserBaseData) Type() string {
	return "EmailAuth_UserData"
}

func (u UserBaseData) WebSafe() map[string]string {
	return map[string]string{
		"email":     u.Email,
		"username":  u.Username,
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

// EmailAuth provides an implementation of a AuthenticationProvider,
// which provides email and password authentication and authorization.
type EmailAuth struct {
	SessionDuration time.Duration

	AuthInitiator nhttp.Handler
	AuthFinalizer nhttp.Handler

	UserStore     UserStore
	UserValidator UserValidator
	Logs          nauth.Logs
	Sessions      providers.HTTPSession
	CookieSigner  *securecookie.SecureCookie
}

// Verify implements the nauth.Authenticator interface.
func (eu EmailAuth) Verify(cm nauth.Claim) (nauth.VerifiedClaim, error) {
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

func (eu EmailAuth) GetSession(req *http.Request) (sessions.Session, error) {
	var ctx = req.Context()
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "EmailAuth.GetSession"); span != nil {
		defer span.Finish()
	}
	var verified sessions.Session
	var session, err = eu.Sessions.Get(req)
	if err != nil {
		return verified, nerror.Wrap(err, "http.Request has no existing auth session")
	}
	return session, nil
}

func (eu EmailAuth) GetSessionClaim(req *http.Request) (nauth.VerifiedClaim, error) {
	var ctx = req.Context()
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "EmailAuth.GetSessionClaim"); span != nil {
		defer span.Finish()
	}

	var verified nauth.VerifiedClaim

	// Retrieve user session from request.
	var session, err = eu.GetSession(req)
	if err != nil {
		return verified, nerror.Forward(err)
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
func (eu EmailAuth) Refresh(res http.ResponseWriter, req *http.Request) error {
	var ctx = req.Context()
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "EmailAuth.Refresh"); span != nil {
		defer span.Finish()
	}
	var userSession, getSessionErr = eu.GetSession(req)
	if getSessionErr != nil {
		return nerror.Forward(getSessionErr)
	}

	var timeLeft = userSession.Expiring.Sub(time.Now())

	// has this expired?
	if timeLeft < 0 {
		if deleteSessionErr := eu.Sessions.DeleteBySid(req.Context(), userSession.ID); deleteSessionErr != nil {
			eu.Logs.Write(njson.MJSON("failed to delete session", func(event npkg.Encoder) {
				event.String("session_id", userSession.ID.String())
				event.String("session_user_id", userSession.User.String())
			}))
		}
		return nerror.New("Expired user session")
	}

	var _, extendLifetimeErr = eu.Sessions.Extend(ctx, userSession.ID, eu.SessionDuration)
	if extendLifetimeErr != nil {
		return nerror.Wrap(extendLifetimeErr, "Failed to extend user session lifetim")
	}

	return nil
}

// Initiate implements the nauth.AuthenticationProvider interface.
//
// Initiate runs the provided User http.handler which handles the rendering of
// necessary login form for the user to input defined credentials for authentication.
func (eu EmailAuth) Initiate(res http.ResponseWriter, req *http.Request) error {
	var ctx = req.Context()
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "EmailAuth.Initiate"); span != nil {
		defer span.Finish()
	}
	if eu.AuthInitiator != nil {
		if initErr := eu.AuthInitiator.Handle(res, req.WithContext(ctx)); initErr != nil {
			return nerror.Wrap(initErr, "Failed to initialize authentication")
		}
	}
	res.WriteHeader(http.StatusNotImplemented)
	return nil
}

type LoginCredentialDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l LoginCredentialDTO) User() string {
	if l.Email != "" {
		return l.Email
	}
	return l.Username
}

func (l LoginCredentialDTO) Validate() error {
	if l.Email == "" && l.Username == "" {
		return nerror.New("invalid login credentials")
	}
	if l.Password == "" {
		return nerror.New("password not provided")
	}
	return nil
}

const (
	USER_CLAIME_VALIDATION_ERROR_MESSAGE = "username/email or password is incorrect"
	USER_SESSION_CREATION_ERROR_MESSAGE  = "Failed to create sessions for user validated claim"
)

// Authenticate implements the nauth.AuthenticationProvider interface.
func (eu EmailAuth) Authenticate(res http.ResponseWriter, req *http.Request) error {
	var ctx = req.Context()
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "EmailAuth.Authenticate"); span != nil {
		defer span.Finish()
	}

	var credentials LoginCredentialDTO
	if decodedErr := json.NewDecoder(req.Body).Decode(&credentials); decodedErr != nil {
		return nerror.Wrap(decodedErr, "Failed to decode request body into LoginCredentialDTO")
	}

	ntrace.WithTag(span, "user", credentials.User())

	var claim nauth.Claim
	claim.Method = CLAIM_TYPE
	claim.Provider = "InHouse"
	claim.Credentials = credentials
	claim.IP = net.ParseIP(req.RemoteAddr)
	claim.Agent = req.Header.Get("User-Agent")

	var validatedClaim, validatedClaimErr = eu.Verify(claim)
	if validatedClaimErr != nil {
		return nerror.Wrap(validatedClaimErr, USER_CLAIME_VALIDATION_ERROR_MESSAGE)
	}

	ntrace.WithTag(span, "validated_data", validatedClaim)

	var newUserSession, newSessionErr = eu.Sessions.Create(ctx, validatedClaim)
	if newSessionErr != nil {
		return nerror.Wrap(newSessionErr, USER_SESSION_CREATION_ERROR_MESSAGE)
	}

	var userSessionCookie, userSessionCookieErr = newUserSession.EncodeToCookie(eu.CookieSigner)
	if userSessionCookieErr != nil {
		return nerror.Wrap(newSessionErr, "failed to write cookie to request")
	}

	http.SetCookie(res, userSessionCookie)
	return nil
}

func (eu EmailAuth) Finalize(res http.ResponseWriter, req *http.Request) error {
	var ctx = req.Context()
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "EmailAuth.Finalize"); span != nil {
		defer span.Finish()
	}

	return eu.AuthFinalizer.ServeHTTP(res, req.WithContext(ctx))
}
