package sessions

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gokit/npkg/njson"

	"github.com/gokit/npkg"

	"github.com/gokit/npkg/nauth"
	"github.com/gokit/npkg/nxid"
)

// Session embodies a current accessible session for a user
// over a underline service.
type Session struct {
	Provider  string
	Method    string
	ID        nxid.ID
	User      nxid.ID
	Expiring  time.Time
	ClaimData map[string]interface{}
	Attached  map[string]interface{}
}

// EncodeToCookie returns a http.Cookie with session encoded into
// it.
func (s *Session) EncodeToCookie() (*http.Cookie, error) {
	var sessionJSON = njson.Object()
	if err := s.EncodeObject(sessionJSON); err != nil {
		return nil, err
	}

	var encodedSession = bytes.NewBuffer(make([]byte, 0, len(sessionJSON.Buf())))
	if _, err := sessionJSON.WriteTo(encodedSession); err != nil {
		return nil, err
	}

	var cookie http.Cookie
	cookie.Name = "_auth_session"
	cookie.Value = base64.StdEncoding.EncodeToString(encodedSession.Bytes())
	return &cookie, nil
}

// EncodeObject implements the npkg.EncodableObject interface.
func (s *Session) EncodeObject(encoder npkg.ObjectEncoder) error {
	if err := encoder.String("method", s.Method); err != nil {
		return err
	}
	if err := encoder.String("id", s.ID.String()); err != nil {
		return err
	}
	if err := encoder.String("user", s.User.String()); err != nil {
		return err
	}
	if err := encoder.String("provider", s.Provider); err != nil {
		return err
	}
	if err := encoder.Int64("expiring", s.Expiring.Unix()); err != nil {
		return err
	}
	if err := encoder.Int64("expiring_nano", s.Expiring.UnixNano()); err != nil {
		return err
	}
	if err := encoder.Object("claim_data", npkg.EncodableMap(s.ClaimData)); err != nil {
		return err
	}
	if err := encoder.Object("attached", npkg.EncodableMap(s.ClaimData)); err != nil {
		return err
	}
	return nil
}

// Sessions embodies what we expect from a session store or provider
// which handles the underline storing and management of sessions.
type Sessions interface {
	// Get retrieves the underline session from request, retrieving
	// underline session from the store from the information retrieved
	// from the request.
	Get(req *http.Request) (Session, error)

	// Delete removes giving session from underline store.
	Delete(session Session) error

	// Extend extends giving session underline lifetime to
	// extend giving session time.
	Extend(session Session) error

	// Create creates new session information for verified claim
	// attaching claim data.
	Create(verifiedClaim nauth.VerifiedClaim) (Session, error)
}
