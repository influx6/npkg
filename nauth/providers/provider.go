package providers

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
