package providers

import (
	"context"
	"github.com/influx6/npkg/nauth"
	"github.com/influx6/npkg/nauth/sessions"
	"github.com/influx6/npkg/nxid"
	"net/http"
	"time"
)

// AuthBehaviour defines a int type used to define what type of behaviour a
// auth provider implementation should follow.
type Operation int

const (
	// API defines the operation of a authentication provider should be off
	// the API form which probably responds with JSON and has no concept of
	// sessions (but this is not a hard fact).
	API Operation = iota + 1

	// EgeSite refers to the application itself used by the user, or CMS which
	// does have the concept of sessions sent as cookies, where a user interacts
	// with the service directly or where such a service is not a API.
	EdgeSite
)

type HTTPSession interface {

	// Get retrieves the underline session from request, retrieving
	// underline session from the store from the information retrieved
	// from the request.
	Get(req *http.Request) (sessions.Session, error)

	// Create creates new session information from verified claim
	// attaching claim data.
	Create(ctx context.Context, claim nauth.VerifiedClaim) (sessions.Session, error)

	// GetByID attempts to retrieve an existing session by it's unique
	// nxid ID.
	GetByID(ctx context.Context, id nxid.ID) (sessions.Session, error)

	// DeleteByUserId removes giving session from underline store.
	DeleteByUserId(ctx context.Context, userId nxid.ID) error

	// DeleteBySid removes giving session from underline store.
	DeleteBySid(ctx context.Context, sid nxid.ID) error

	// Extend extends giving session underline lifetime to
	// extend giving session time.
	Extend(ctx context.Context, id nxid.ID, lifeTime time.Duration) (sessions.Session, error)
}
