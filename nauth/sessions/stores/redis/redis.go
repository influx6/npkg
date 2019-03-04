package redis

import (
	"github.com/go-redis/redis"
	"github.com/gokit/npkg/nauth/sessions"
)

// RedisStore implements session management, storage and access using redis as
// underline store.
type RedisStore struct {
	config *redis.Options
	client *redis.Client
}

// NewRedisStore returns a new instance of a redis store.
func NewRedisStore(config redis.Options) (*RedisStore, error) {
	var red RedisStore
	red.config = &config
	if err := red.createConnection(); err != nil {
		return nil, err
	}
	return &red, nil
}

func (rd *RedisStore) createConnection() error {
	client := redis.NewClient(rd.config)
	status := client.Ping()
	if err := status.Err(); err != nil {
		return err
	}
	rd.client = client
	return nil
}

// Save adds giving session into storage using redis as underline store.
func (rd *RedisStore) Save(s sessions.Session) error {
	panic("implement me")
}

// Get returns giving session stored with giving key, returning an
// error if not found.
func (rd *RedisStore) Get(key string) (sessions.Session, error) {
	panic("implement me")
}

// Update updates giving session stored with giving key. It updates
// the underline data.
func (rd *RedisStore) Update(key string, s sessions.Session) error {
	panic("implement me")
}

// Remove removes underline key from the redis store after retrieving it and
// returning giving session.
func (rd *RedisStore) Remove(key string) (sessions.Session, error) {
	panic("implement me")
}
