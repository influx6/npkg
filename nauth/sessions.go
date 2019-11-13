package nauth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/influx6/npkg/nstorage"

	"github.com/gorilla/securecookie"
	"github.com/influx6/npkg"
	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/njson"
	"github.com/influx6/npkg/ntrace"
	"github.com/influx6/npkg/nunsafe"
	"github.com/influx6/npkg/nxid"
)

const (
	// SessionCookieName defines the name used for the Session cookie.
	SessionCookieName = "_auth_session"

	// SessionUserDataKeyName defines the name used for the user data
	// attached to a session.
	SessionUserDataKeyName = "_auth_session_user_data"

	// CookieHeaderName defines the name of the cookie header.
	CookieHeaderName = "Set-Cookie"
)

// Session embodies a current accessible session for a user
// over a underline service.
type Session struct {
	Provider string    `json:"provider"`
	Method   string    `json:"method"`
	Browser  string    `json:"browser"`
	IP       net.IP    `json:"ip"`
	ID       nxid.ID   `json:"id"`
	User     nxid.ID   `json:"user"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Expiring time.Time `json:"expiring"`
}

// EncodeToCookie returns a http.Cookie with session encoded into
// it. If the signer is provided, then an attempt is made to sign the
// the value cookie using provided signer which uses gob underneath.
//
// If no signer is provided then, the provided cookie is json encoded,
// transformed into base64, then set as cookie value.
func (s *Session) EncodeToCookie(signer *securecookie.SecureCookie) (*http.Cookie, error) {
	var cookie http.Cookie
	cookie.HttpOnly = true
	cookie.Name = SessionCookieName

	var sessionJSON = njson.Object()
	if err := s.EncodeForCookie(sessionJSON); err != nil {
		return nil, err
	}

	if signer == nil {
		cookie.Value = base64.StdEncoding.EncodeToString(sessionJSON.Buf())
		return &cookie, nil
	}

	cookie.Secure = true
	encrypted, err := signer.Encode(cookie.Name, sessionJSON.Message())
	if err != nil {
		return nil, err
	}
	cookie.Value = encrypted
	return &cookie, nil
}

// Validate returns an error if giving session was invalid.
func (s *Session) Validate() error {
	if s.Created.IsZero() {
		return nerror.New("session.Created has no created time stamp")
	}
	if s.Updated.IsZero() {
		return nerror.New("session.Updated has no updated time stamp")
	}
	if s.Expiring.IsZero() {
		return nerror.New("session.Expiring has no expiration time stamp")
	}
	if len(s.ID) == 0 {
		return nerror.New("session.ID must have a valid value")
	}
	if len(s.User) == 0 {
		return nerror.New("session.User must have a valid value")
	}
	if len(s.Provider) == 0 {
		return nerror.New("session.Provider must have a valid value")
	}
	if len(s.Method) == 0 {
		return nerror.New("session.Method must have a valid value")
	}
	return nil
}

// Writes giving session as a cookie into the provided http.ResponseWriter.
func (s *Session) Write(signer *securecookie.SecureCookie, w http.ResponseWriter, mods ...func(*http.Cookie)) error {
	var cookie, err = s.EncodeToCookie(signer)
	if err != nil {
		return err
	}

	for _, mod := range mods {
		mod(cookie)
	}

	w.Header().Add(CookieHeaderName, cookie.String())
	return nil
}

// EncodeForCookie encodes giving session to npkg.ObjectEncoder for
// delivery to the client.
func (s *Session) EncodeForCookie(encoder npkg.ObjectEncoder) error {
	if err := s.Validate(); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("id", s.ID.String()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Object("attached", npkg.EncodableMap(s.Attached)); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("method", s.Method); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("provider", s.Provider); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("user", s.User.String()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Int64("created", s.Created.Unix()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Int64("updated", s.Updated.Unix()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Int64("expiring", s.Expiring.Unix()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Int64("expiring_nano", s.Expiring.UnixNano()); err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

// EncodeObject implements the npkg.EncodableObject interface.
func (s *Session) EncodeObject(encoder npkg.ObjectEncoder) error {
	if err := s.Validate(); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("browser", s.Browser); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("method", s.Method); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("id", s.ID.String()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("user", s.User.String()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.String("provider", s.Provider); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Int64("created", s.Created.Unix()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Int64("updated", s.Updated.Unix()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Int64("expiring", s.Expiring.Unix()); err != nil {
		return nerror.WrapOnly(err)
	}
	if err := encoder.Int64("expiring_nano", s.Expiring.UnixNano()); err != nil {
		return nerror.WrapOnly(err)
	}
	if len(s.IP) != 0 {
		if err := encoder.String("method", s.IP.String()); err != nil {
			return nerror.WrapOnly(err)
		}
	}
	return nil
}

//**********************************************
// Session Codec
//**********************************************

// SessionEncoder defines what we expect from a encoder of Session type.
type SessionEncoder interface {
	Encode(w io.Writer, s Session) error
}

// SessionDecoder defines what we expect from a encoder of Session type.
type SessionDecoder interface {
	Decode(r io.Reader, s *Session) error
}

// SessionCodec exposes an interface which combines the SessionEncoder and
// SessionDecoder interfaces.
type SessionCodec interface {
	SessionEncoder
	SessionDecoder
}

// GobSessionCodec implements the SessionCodec interface for using
// the gob codec.
type GobSessionCodec struct{}

// Encode encodes giving session using the internal gob format.
// Returning provided data.
func (gb *GobSessionCodec) Encode(w io.Writer, s Session) error {
	if err := gob.NewEncoder(w).Encode(s); err != nil {
		return nerror.Wrap(err, "Failed to encode giving session")
	}
	return nil
}

// Decode decodes giving data into provided session instance.
func (gb *GobSessionCodec) Decode(r io.Reader, s *Session) error {
	if err := gob.NewDecoder(r).Decode(s); err != nil {
		return nerror.Wrap(err, "Failed to decode bytes as gob into nauth.Session")
	}
	return nil
}

//**********************************************
// Session Storage
//**********************************************

var (
	_          SessionsStorage = (*SessionStorage)(nil)
	readerPool                 = sync.Pool{
		New: func() interface{} {
			return bytes.NewReader(nil)
		},
	}
)

// SessionsStorage defines a underline store for a giving session by key.
type SessionsStorage interface {
	Save(context.Context, Session) error
	Update(context.Context, Session) error
	GetAll(context.Context) ([]Session, error)
	Remove(context.Context, string) (Session, error)
	GetByID(context.Context, string) (Session, error)
	GetByUser(context.Context, string) (Session, error)
	GetAllByUser(context.Context, string) ([]Session, error)
}

// SessionStorage implements a storage type for CRUD operations on
// sessions.
type SessionStorage struct {
	Codec SessionCodec
	Store nstorage.ExpirableStorage
}

// NewSessionStorage returns a new instance of a SessionStorage.
func NewSessionStorage(codec SessionCodec, store nstorage.ExpirableStorage) *SessionStorage {
	return &SessionStorage{
		Codec: codec,
		Store: store,
	}
}

// Save adds giving session into underline store.
//
// It sets the session to expire within the storage based on
// the giving session's expiration duration.
//
// Save calculates the ttl by subtracting the Session.Created value from
// the Session.Expiring value.
func (s *SessionStorage) Save(ctx context.Context, s Session) error {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionStorage.Save"); span != nil {
		defer span.Finish()
	}

	if err := s.Validate(); err != nil {
		return nerror.Wrap(err, "Session failed validation")
	}

	var content = bytes.NewBuffer(make([]byte, 0, 512))
	if err := s.Codec.Encode(content, s); err != nil {
		return nerror.Wrap(err, "Failed to encode data")
	}

	// Calculate expiration for giving value.
	var expiration = s.Expiring.Sub(s.Created)
	if err := s.Store.SaveTTL(s.ID.String(), content.Bytes(), expiration); err != nil {
		return nerror.Wrap(err, "Failed to save encoded session")
	}
	return nil
}

// Update attempts to update existing session key within store if
// still available.
//
// Update calculates the ttl by subtracting the Session.Updated value from
// the Session.Expiring value.
func (s *SessionStorage) Update(ctx context.Context, s Session) error {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionStorage.Update"); span != nil {
		defer span.Finish()
	}
	if err := s.Validate(); err != nil {
		return nerror.Wrap(err, "Session failed validation")
	}

	var content = bytes.NewBuffer(make([]byte, 0, 512))
	if err := s.Codec.Encode(content, s); err != nil {
		return nerror.Wrap(err, "Failed to encode data")
	}

	// Calculate expiration for giving value.
	var expiration = s.Expiring.Sub(s.Updated)
	if err := s.Store.UpdateTTL(s.ID.String(), content.Bytes(), expiration); err != nil {
		return nerror.Wrap(err, "Failed to update encoded session")
	}
	return nil
}

// GetByID retrieves giving session from store based on the provided
// session ID value.
func (s *SessionStorage) GetByID(ctx context.Context, key string) (Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionStorage.Get"); span != nil {
		defer span.Finish()
	}
	var session Session
	var sessionBytes, err = s.Store.Get(key)
	if err != nil {
		return session, nerror.WrapOnly(err)
	}

	var reader = readerPool.Get().(*bytes.Reader)
	defer readerPool.Put(reader)

	reader.Reset(sessionBytes)
	defer reader.Reset(nil)
	if err := s.Codec.Decode(reader, &session); err != nil {
		return session, nerror.WrapOnly(err)
	}
	return session, nil
}

// Remove removes underline session if still present from underline store.
func (s *SessionStorage) Remove(ctx context.Context, key string) (Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionStorage.Remove"); span != nil {
		defer span.Finish()
	}
	var session Session
	var sessionBytes, err = s.Store.Remove(key)
	if err != nil {
		return session, nerror.WrapOnly(err)
	}

	var reader = readerPool.Get().(*bytes.Reader)
	defer readerPool.Put(reader)

	reader.Reset(sessionBytes)
	defer reader.Reset(nil)
	if err := s.Codec.Decode(reader, &session); err != nil {
		return session, nerror.WrapOnly(err)
	}
	return session, nil
}

//**********************************************
// Sessions
//**********************************************

const (
	// SessionKey defines the key used to save a session instance in a
	// request object.
	SessionKey = sessionKey("nauth-session")
)

type sessionKey string

// GetSessionFromContext returns a Session instance attached to a
// context. It returns true, if found as second value or false if
// not found.
func GetSessionFromContext(ctx context.Context) (Session, bool) {
	var ok bool
	var session Session
	if session, ok = ctx.Value(SessionKey).(Session); ok {
		return session, true
	}
	return session, false
}

// GetUserDataFromSession returns possible user data attached to
// giving session.
func GetUserDataFromSession(s Session) (interface{}, error) {
	if userData, ok := s.Attached[SessionUserDataKeyName]; ok {
		return userData, nil
	}
	return nil, nerror.New("no user data attached")
}

// AddSessionToContext adds session into provided context, returning new context.
func AddSessionToContext(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, SessionKey, session)
}

// Sessions embodies what we expect from a session store or provider
// which handles the underline storing and management of sessions.
type Sessions interface {
	// Get retrieves the underline session from request, retrieving
	// underline session from the store from the information retrieved
	// from the request.
	Get(req *http.Request) (Session, error)

	// Create creates new session information for verified claim
	// attaching claim data.
	Create(ctx context.Context, claim VerifiedClaim) (Session, error)

	// GetByID attempts to retrieve an existing session by it's unique
	// nxid ID.
	GetByID(ctx context.Context, id nxid.ID) (Session, error)

	// Delete removes giving session from underline store.
	Delete(ctx context.Context, id nxid.ID) (Session, error)

	// Extend extends giving session underline lifetime to
	// extend giving session time.
	Extend(ctx context.Context, id nxid.ID) (Session, error)
}

// SessionConfig defines the configuration which are the values for giving
// Sessions provider.
type SessionConfig struct {
	// Lifetime is the default time-to-live (ttl) for a new session created.
	Lifetime time.Duration

	// Extension is the extension ttl to be used when extending a existing
	// session lifetime.
	Extension time.Duration

	// Storage defines the core session storage to be used by provided session
	// provider for use.
	Storage SessionStorage

	// Signer sets the signer to be used to sign giving cookie.
	Signer *securecookie.SecureCookie
}

// Validate validates giving session config is valid.
func (s *SessionConfig) Validate() error {
	if s.Storage == nil {
		return nerror.New("SessionConfig.Storage requires a SessionStorage")
	}
	if s.Extension <= 0 {
		return nerror.New("SessionConfig.Extension can't be zero or below it")
	}
	if s.Lifetime <= 0 {
		return nerror.New("SessionConfig.Lifetime can't be zero or below it")
	}
	return nil
}

// SessionImpl implements the Sessions interface, providing the
// necessary decorator that uses a SessionStorage for the managing
// of a session with a http request.
type SessionImpl struct {
	Config SessionConfig
}

// NewSessionImpl returns a new instance of a SessionImpl which implements
// the Sessions interface.
func NewSessionImpl(config SessionConfig) (*SessionImpl, error) {
	if err := config.Validate(); err != nil {
		return nil, nerror.WrapOnly(err)
	}
	var impl SessionImpl
	impl.Config = config
	return &impl, nil
}

// Get returns a Session parsed out of the request or already attached to the
// request context.
//
// The function attempts to retrieve an existing session from the underlying request
// context if already existing, else tries to get the cookie using the SessionCookieName
// decoding the value and returning the Session object.
func (s *SessionImpl) Get(req *http.Request) (Session, error) {
	var span openTracing.Span
	var ctx = req.Context()
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionImpl.Get"); span != nil {
		defer span.Finish()
	}

	var found bool
	var session Session
	if session, found = GetSessionFromContext(ctx); found {
		return session, nil
	}

	var sessionCookie *http.Cookie
	var cookies = req.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == SessionCookieName {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		return session, nerror.New("No session cookie found in request")
	}

	var err error
	var content string

	if s.Config.Signer != nil {
		if err = s.Config.Signer.Decode(SessionCookieName, sessionCookie.Value, &content); err != nil {
			return session, nerror.WrapOnly(err)
		}

		var tmp Session
		if err = json.Umarshal(nunsafe.String2Bytes(content), &tmp); err != nil {
			return session, nerror.WrapOnly(err)
		}

		session, err = s.Config.Storage.Get(ctx, tmp.ID.String())
		if err != nil {
			return session, nerror.WrapOnly(err)
		}

		return session, nil
	}

	content, err = base64.StdEncoding.DecodeString(sessionCookie.Value)
	if err != nil {
		return session, nerror.WrapOnly(err)
	}

	var tmp Session
	if err = json.Umarshal(nunsafe.String2Bytes(content), &tmp); err != nil {
		return session, nerror.WrapOnly(err)
	}

	session, err = s.Config.Storage.Get(ctx, tmp.ID.String())
	if err != nil {
		return session, nerror.WrapOnly(err)
	}
	return session, nil
}

// GetByID retrieves a giving Session from the underline SessionStorage.
func (s *SessionImpl) GetByID(ctx context.Context, id nxid.ID) (Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionImpl.GetByID"); span != nil {
		defer span.Finish()
	}

	session, err = s.Config.Storage.Get(req.Context(), id.String())
	if err != nil {
		return session, nerror.WrapOnly(err)
	}
	return session, nil
}

// Create creates a new session for giving user if non currently exists
// within underline storage, it returns a new session representing
// said user with associated information to be included within
// such session.
func (s *SessionImpl) Create(ctx context.Context, claim VerifiedClaim) (Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionImpl.Create"); span != nil {
		defer span.Finish()
	}

	var session Session
	session.ID = nxid.New()
	session.User = claim.User
	session.Method = claim.Method
	session.Provider = claim.Provider
	session.Created = time.Now()
	session.Updated = session.Created
	session.Expiring = session.Created.Add(s.Config.Lifetime)
	if claim.Data != nil {
		session.Attached = map[string]interface{}{
			SessionUserDataKeyName: claim.Data,
		}
	}

	if err := s.Config.Storage.Save(ctx, session); err != nil {
		return session, nerror.WrapOnly(err)
	}
	return session, nil
}

// Delete removes a giving Session from the underline SessionStorage.
func (s *SessionImpl) Delete(ctx context.Context, id nxid.ID) (Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionImpl.Delete"); span != nil {
		defer span.Finish()
	}
	return s.Config.Storage.Remove(ctx, id.String())
}

// Extend extends a giving Session with default extension duration.
func (s *SessionImpl) Extend(ctx context.Context, id nxid.ID) error {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "SessionImpl.Extend"); span != nil {
		defer span.Finish()
	}

	var session, err = s.Config.Storage.Get(ctx, id.String())
	if err != nil {
		return nerror.WrapOnly(err)
	}

	session.Updated = time.Now()
	session.Expiring = session.Updated.Add(s.Config.Extension)
	if err := s.Config.Storage.Update(ctx, session); err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

//**********************************************
// init
//**********************************************

func init() {
	gob.Register((*Session)(nil))
}
