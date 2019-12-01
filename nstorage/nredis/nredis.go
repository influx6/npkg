package nredis

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/go-redis/redis"

	"github.com/influx6/npkg/nerror"
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
		return nil, nerror.WrapOnly(err)
	}
	return &red, nil
}

// FromRedisStore returns a new instance of a RedisStore using giving client.
func FromRedisStore(hash string, conn *redis.Client) (*RedisStore, error) {
	if status := conn.Ping(); status.Err() != nil {
		return nil, status.Err()
	}

	var red RedisStore
	red.hashList = hash + "_keys"
	red.hashElem = hash + "_item"
	red.client = conn
	return &red, nil
}

// createConnection attempts to create a new redis connection.
func (rd *RedisStore) createConnection() error {
	client := redis.NewClient(rd.config)
	status := client.Ping()
	if err := status.Err(); err != nil {
		return nerror.WrapOnly(err)
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
		return nil, nerror.WrapOnly(err)
	}
	return nstatus.Val(), nil
}

// Each runs through all elements for giving store, skipping keys
// in redis who have no data or an empty byte slice.
func (rd *RedisStore) Each(fn func([]byte, string) bool) error {
	var nstatus = rd.client.SMembers(rd.hashList)
	if err := nstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}
	for _, item := range nstatus.Val() {
		var gstatus = rd.client.Get(item)
		if err := gstatus.Err(); err == nil {
			if !fn(string2Bytes(gstatus.Val()), item) {
				return nil
			}
		}
	}
	return nil
}

// FindAll returns all match elements for giving function.
func (rd *RedisStore) FindAll(fn func([]byte, string) bool, count int) ([][]byte, error) {
	return rd.FindEach(fn, -1)
}

// Find returns match content for giving function.
func (rd *RedisStore) Find(fn func([]byte, string) bool, count int) ([]byte, error) {
	var res, err = rd.FindEach(fn, 1)
	if err != nil {
		return nil, err
	}
	if len(res) == 1 {
		return res[0], err
	}
	return nil, nil
}

// FindEach returns all matching values within store, if elements found match giving
// count then all values returned.
func (rd *RedisStore) FindEach(fn func([]byte, string) bool, count int) ([][]byte, error) {
	var result [][]byte

	var nstatus = rd.client.SMembers(rd.hashList)
	if err := nstatus.Err(); err != nil {
		return result, nerror.WrapOnly(err)
	}
	for _, item := range nstatus.Val() {
		if count > 0 && count == len(result) {
			return result, nil
		}

		var gstatus = rd.client.Get(item)
		if err := gstatus.Err(); err == nil {
			var data = string2Bytes(gstatus.Val())
			if fn(data, item) {
				result = append(result, data)
				continue
			}
		}
	}
	return result, nil
}

// Exists returns true/false if giving key exists.
func (rd *RedisStore) Exists(key string) (bool, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.SIsMember(rd.hashList, hashKey)
	if err := nstatus.Err(); err != nil {
		return false, nerror.WrapOnly(err)
	}
	return nstatus.Val(), nil
}

// exists returns true/false if giving key is set in redis.
func (rd *RedisStore) exists(key string) (bool, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.Exists(hashKey)
	if err := nstatus.Err(); err != nil {
		return false, nerror.WrapOnly(err)
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
		return nerror.WrapOnly(err)
	}
	return nil
}

// Save adds giving session into storage using redis as underline store.
func (rd *RedisStore) Save(key string, data []byte) error {
	return rd.SaveTTL(key, data, 0)
}

// SaveTTL adds giving session into storage using redis as underline store, with provided
// expiration.
// Duration of 0 means no expiration.
func (rd *RedisStore) SaveTTL(key string, data []byte, expiration time.Duration) error {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.SAdd(rd.hashList, hashKey)
	if err := nstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}

	var nset = rd.client.Set(hashKey, data, expiration)
	if err := nset.Err(); err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

// Update updates giving key with new data slice with 0 duration.
func (rd *RedisStore) Update(key string, data []byte) error {
	return rd.UpdateTTL(key, data, 0)
}

// UpdateTTL updates giving session stored with giving key. It updates
// the underline data.
//
// if expiration is zero then giving value expiration will not be reset but left
// as is.
func (rd *RedisStore) UpdateTTL(key string, data []byte, expiration time.Duration) error {
	var hashKey = rd.getHashKey(key)
	var fstatus = rd.client.SIsMember(rd.hashList, hashKey)
	if err := fstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}
	if !fstatus.Val() {
		return nerror.New("key does not exist")
	}

	if len(data) == 0 {
		return rd.remove(key)
	}

	var newTTL time.Duration
	if expiration > 0 {
		var ttlstatus = rd.client.PTTL(hashKey)
		if err := ttlstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}

		newTTL = ttlstatus.Val() + expiration
		if ttlstatus.Val() <= 0 {
			newTTL = expiration
		}
	}

	var nset = rd.client.Set(hashKey, data, newTTL)
	if err := nset.Err(); err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

// TTL returns current expiration time for giving key.
func (rd *RedisStore) TTL(key string) (time.Duration, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.PTTL(hashKey)
	if err := nstatus.Err(); err != nil {
		return 0, nerror.WrapOnly(err)
	}
	if nstatus.Val() < 0 {
		return 0, nil
	}
	return nstatus.Val(), nil
}

// ExtendTTL extends the expiration of a giving key if it exists, the duration is expected to be
// in milliseconds. If expiration value is zero then we consider that you wish to remove the expiration.
func (rd *RedisStore) ExtendTTL(key string, expiration time.Duration) error {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.PTTL(hashKey)
	if err := nstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}

	if nstatus.Val() < 0 {
		return nil
	}

	if expiration == 0 {
		var exstatus = rd.client.Persist(hashKey)
		return exstatus.Err()
	}

	var newExpiration = expiration + nstatus.Val()
	var exstatus = rd.client.Expire(hashKey, newExpiration)
	return exstatus.Err()
}

// ResetTTL resets giving expiration value to provided duration.
//
// A duration of zero persists the giving key.
func (rd *RedisStore) ResetTTL(key string, expiration time.Duration) error {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.PTTL(hashKey)
	if err := nstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}

	if nstatus.Val() < 0 {
		return nil
	}

	if expiration == 0 {
		var exstatus = rd.client.Persist(hashKey)
		return exstatus.Err()
	}

	var exstatus = rd.client.Expire(hashKey, expiration)
	return exstatus.Err()
}

// Get returns giving session stored with giving key, returning an
// error if not found.
func (rd *RedisStore) Get(key string) ([]byte, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.Get(hashKey)
	if err := nstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}
	return string2Bytes(nstatus.Val()), nil
}

// Remove removes underline key from the redis store after retrieving it and
// returning giving session.
func (rd *RedisStore) Remove(key string) ([]byte, error) {
	var hashKey = rd.getHashKey(key)
	var nstatus = rd.client.Get(hashKey)
	if err := nstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}
	var mstatus = rd.client.SRem(rd.hashList, hashKey)
	if err := mstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}
	var dstatus = rd.client.Del(hashKey)
	if err := dstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}
	return string2Bytes(nstatus.Val()), nil
}

func (rd *RedisStore) remove(key string) error {
	var hashKey = rd.getHashKey(key)
	var dstatus = rd.client.Del(hashKey)
	return dstatus.Err()
}

//*****************************************************
// internal methods
//*****************************************************

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
