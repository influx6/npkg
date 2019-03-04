package nredis

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/go-redis/redis"
	"github.com/gokit/npkg/nerror"
)

// RedisStore implements session management, storage and access using redis as
// underline store.
type RedisStore struct {
	hashList string
	hashElem string
	config   *redis.Options
	client   *redis.Client
}

// NewRedisStore returns a new instance of a redis store.
func NewRedisStore(hash string, config redis.Options) (*RedisStore, error) {
	var red RedisStore
	red.hashList = hash + "_keys"
	red.hashElem = hash + "_item"
	red.config = &config
	if err := red.createConnection(); err != nil {
		return nil, err
	}
	return &red, nil
}

// createConnection attempts to create a new redis connection.
func (rd *RedisStore) createConnection() error {
	client := redis.NewClient(rd.config)
	status := client.Ping()
	if err := status.Err(); err != nil {
		return err
	}
	rd.client = client
	return nil
}

// getHashKey returns the key of giving key used to store
// both in map and redis.
func (rd *RedisStore) getHashKey(key string) string {
	return fmt.Sprintf("%s_%s", rd.hashList, key)
}

// Keys returns all giving keys of elements within store.
func (rd *RedisStore) Keys() ([]string, error) {
	var nstatus = rd.client.SMembers(rd.hashList)
	if err := nstatus.Err(); err != nil {
		return nil, err
	}
	return nstatus.Val(), nil
}

// Each runs through all elements for giving store, skipping keys
// in redis who have no data or an empty byte slice.
func (rd *RedisStore) Each(fn func([]byte, string) bool) error {
	var nstatus = rd.client.SMembers(rd.hashList)
	if err := nstatus.Err(); err != nil {
		return err
	}
	for _, item := range nstatus.Val() {
		if data, err := rd.Get(item); err == nil && len(data) > 0 {
			if fn(data, item) {
				continue
			}
		}
	}
	return nil
}

// Exists returns true/false if giving key exists.
func (rd *RedisStore) Exists(key string) (bool, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.SIsMember(rd.hashList, hashKey)
	if err := nstatus.Err(); err != nil {
		return false, err
	}
	return nstatus.Val(), nil
}

// exists returns true/false if giving key is set in redis.
func (rd *RedisStore) exists(key string) (bool, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.Exists(hashKey)
	if err := nstatus.Err(); err != nil {
		return false, err
	}
	return nstatus.Val() == 1, nil
}

// expire expires giving key set from underline hash set.
func (rd *RedisStore) expire(keys []string) error {
	var items = make([]interface{}, len(keys))
	for index, elem := range keys {
		items[index] = elem
	}
	var nstatus = rd.client.SRem(rd.hashList, items...)
	if err := nstatus.Err(); err != nil {
		return err
	}
	return nil
}

// Save adds giving session into storage using redis as underline store.
func (rd *RedisStore) Save(key string, data []byte, expiration time.Duration) error {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.SAdd(rd.hashList, hashKey)
	if err := nstatus.Err(); err != nil {
		return err
	}

	var nset = rd.client.Set(hashKey, data, expiration)
	if err := nset.Err(); err != nil {
		return err
	}
	return nil
}

// ExtendTTL extends the expiration of a giving key if it exists, the duration is expected to be
// in seconds.
func (rd *RedisStore) ExtendTTL(key string, expiration time.Duration) error {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.TTL(hashKey)
	if err := nstatus.Err(); err != nil {
		return err
	}
	if nstatus.Val() < 0 {
		return nil
	}

	var newExpriration = expiration + nstatus.Val()
	var exstatus = rd.client.Expire(hashKey, newExpriration)
	return exstatus.Err()
}

// Get returns giving session stored with giving key, returning an
// error if not found.
func (rd *RedisStore) Get(key string) ([]byte, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.Get(hashKey)
	if err := nstatus.Err(); err != nil {
		return nil, err
	}
	return string2Bytes(nstatus.Val()), nil
}

// Update updates giving session stored with giving key. It updates
// the underline data.
func (rd *RedisStore) Update(key string, data []byte, expiration time.Duration) error {
	var hashKey = rd.getHashKey(key)
	var found, err = rd.Exists(hashKey)
	if err != nil {
		return err
	}
	if !found {
		return nerror.New("key does not exist")
	}

	if len(data) == 0 {
		return rd.remove(key)
	}

	var nset = rd.client.Set(hashKey, data, expiration)
	if err := nset.Err(); err != nil {
		return err
	}
	return nil
}

// Remove removes underline key from the redis store after retrieving it and
// returning giving session.
func (rd *RedisStore) Remove(key string) ([]byte, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.Get(hashKey)
	if err := nstatus.Err(); err != nil {
		return nil, err
	}
	var mstatus = rd.client.SRem(rd.hashList, hashKey)
	if err := mstatus.Err(); err != nil {
		return nil, err
	}
	var dstatus = rd.client.Del(hashKey)
	if err := dstatus.Err(); err != nil {
		return nil, err
	}
	return string2Bytes(nstatus.Val()), nil
}

func (rd *RedisStore) remove(key string) error {
	var hashKey = rd.getHashKey(key)
	var dstatus = rd.client.Del(hashKey)
	return dstatus.Err()
}

//*****************************************************
// unsafe methods
//*****************************************************

// byte2String converts a byte slice into a string.
func bytes2String(bc []byte) string {
	return *(*string)(unsafe.Pointer(&bc))
}

// string2Bytes converts a string into a byte slice.
func string2Bytes(bc string) []byte {
	return *(*[]byte)(unsafe.Pointer(&bc))
}
