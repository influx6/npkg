package nauth

import (
	"bytes"
	"context"

	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nstorage/nredis"
	"github.com/influx6/npkg/ntrace"
	openTracing "github.com/opentracing/opentracing-go"
)

var _ SessionsStorage = (*RedisSessionStore)(nil)

// RedisSessionStore implements a storage type for CRUD operations on
// sessions.
type RedisSessionStore struct {
	Codec SessionCodec
	Store *nredis.RedisStore
}

// NewRedisSessionStore returns a new instance of a RedisSessionStore.
func NewRedisSessionStore(codec SessionCodec, store *nredis.RedisStore) *RedisSessionStore {
	return &RedisSessionStore{
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
func (s *RedisSessionStore) Save(ctx context.Context, se Session) error {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "RedisSessionStore.Save"); span != nil {
		defer span.Finish()
	}

	if err := se.Validate(); err != nil {
		return nerror.Wrap(err, "Session failed validation")
	}

	var content = bytes.NewBuffer(make([]byte, 0, 512))
	if err := s.Codec.Encode(content, se); err != nil {
		return nerror.Wrap(err, "Failed to encode data")
	}

	// Calculate expiration for giving value.
	var expiration = se.Expiring.Sub(se.Created)
	if err := s.Store.SaveTTL(se.ID.String(), content.Bytes(), expiration); err != nil {
		return nerror.Wrap(err, "Failed to save encoded session")
	}
	return nil
}

// Update attempts to update existing session key within store if
// still available.
//
// Update calculates the ttl by subtracting the Session.Updated value from
// the Session.Expiring value.
func (s *RedisSessionStore) Update(ctx context.Context, se Session) error {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "RedisSessionStore.Update"); span != nil {
		defer span.Finish()
	}
	if err := se.Validate(); err != nil {
		return nerror.Wrap(err, "Session failed validation")
	}

	var content = bytes.NewBuffer(make([]byte, 0, 512))
	if err := s.Codec.Encode(content, se); err != nil {
		return nerror.Wrap(err, "Failed to encode data")
	}

	// Calculate expiration for giving value.
	var expiration = se.Expiring.Sub(se.Updated)
	if err := s.Store.UpdateTTL(se.ID.String(), content.Bytes(), expiration); err != nil {
		return nerror.Wrap(err, "Failed to update encoded session")
	}
	return nil
}

// GetAll returns all sessions stored within store.
func (s *RedisSessionStore) GetAll(ctx context.Context) ([]Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "RedisSessionStore.Update"); span != nil {
		defer span.Finish()
	}

	var records []Session
	var err = s.Store.Find((content []byte, id string) bool {
		return true
	})
	if err != nil {
		return session, nerror.WrapOnly(err)
	}
	return records, nil
}

// GetByUser retrieves giving session from store based on the provided
// session user value.
func (s *RedisSessionStore) GetByUser(ctx context.Context, key string) (Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "RedisSessionStore.Get"); span != nil {
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

// GetByID retrieves giving session from store based on the provided
// session ID value.
func (s *RedisSessionStore) GetByID(ctx context.Context, key string) (Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "RedisSessionStore.Get"); span != nil {
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
func (s *RedisSessionStore) Remove(ctx context.Context, key string) (Session, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewSpanFromContext(ctx, "RedisSessionStore.Remove"); span != nil {
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
