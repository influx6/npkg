package nredis

import (
	"context"
	regexp2 "regexp"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"

	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nstorage"
	"github.com/influx6/npkg/nunsafe"
)

var _ nstorage.ExpirableStore = (*RedisStore)(nil)

// RedisStore implements session management, storage and access using redis as
// underline store.
type RedisStore struct {
	ctx       context.Context
	tableName string
	hashList  string
	hashZList string
	hashElem  string
	Config    *redis.Options
	Client    *redis.Client
}

// NewRedisStore returns a new instance of a redis store.
func NewRedisStore(ctx context.Context, tableName string, config redis.Options) (*RedisStore, error) {
	var red RedisStore
	red.ctx = ctx
	red.tableName = tableName
	red.hashList = tableName + "_keys"
	red.hashElem = tableName + "_item"
	red.Config = &config
	if err := red.createConnection(); err != nil {
		return nil, nerror.WrapOnly(err)
	}
	return &red, nil
}

// FromRedisStore returns a new instance of a RedisStore using giving client.
func FromRedisStore(ctx context.Context, tableName string, conn *redis.Client) (*RedisStore, error) {
	if status := conn.Ping(ctx); status.Err() != nil {
		return nil, status.Err()
	}

	var red RedisStore
	red.ctx = ctx
	red.tableName = tableName
	red.hashList = tableName + "_keys"
	red.hashElem = tableName + "_item"
	red.hashZList = tableName + "_zset"
	red.Client = conn
	return &red, nil
}

// createConnection attempts to create a new redis connection.
func (rd *RedisStore) createConnection() error {
	client := redis.NewClient(rd.Config)
	status := client.Ping(rd.ctx)
	if err := status.Err(); err != nil {
		return nerror.WrapOnly(err)
	}
	rd.Client = client
	return nil
}

func (rd *RedisStore) Close() error {
	return rd.Client.Close()
}

// doHashKey returns formatted for unique form towards using creating
// efficient hashmaps to contain list of keys.
func (rd *RedisStore) doHashKey(key string) string {
	return strings.Join([]string{rd.hashElem, key}, "_")
}

func (rd *RedisStore) unHashKey(key string) string {
	return strings.TrimPrefix(key, rd.hashElem+"_")
}

func (rd *RedisStore) unHashKeyList(keys []string) []string {
	for index, key := range keys {
		keys[index] = rd.unHashKey(key)
	}
	return keys
}

// Keys returns all giving keys of elements within store.
func (rd *RedisStore) Keys() ([]string, error) {
	var nstatus = rd.Client.SMembers(rd.ctx, rd.hashList)
	if err := nstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}

	return rd.unHashKeyList(nstatus.Val()), nil
}

// Each runs through all elements for giving store, skipping keys
// in redis who have no data or an empty byte slice.
func (rd *RedisStore) Each(fn nstorage.EachItem) error {
	var nstatus = rd.Client.SMembers(rd.ctx, rd.hashList)
	if err := nstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}

	var keys = nstatus.Val()
	var pipeliner = rd.Client.Pipeline()

	var values = make([]*redis.StringCmd, len(keys))
	for index, key := range keys {
		var result = pipeliner.Get(rd.ctx, key)
		values[index] = result
	}

	var _, err = pipeliner.Exec(rd.ctx)
	if err != nil && err != redis.Nil {
		return nerror.WrapOnly(err)
	}

	for index, item := range values {
		if item.Err() != nil {
			continue
		}
		var key = keys[index]
		var data = nunsafe.String2Bytes(item.Val())
		if doErr := fn(data, key); doErr != nil {
			if nerror.IsAny(doErr, nstorage.ErrJustStop) {
				return nil
			}
			return doErr
		}
	}
	return nil
}

// EachKeyPrefix returns all matching values within store, if elements found match giving
// count then all values returned.
//
// if an error occurs, the partially collected list of keys and error is returned.
//
// Return nstorage.ErrJustStop if you want to just stop iterating.
func (rd *RedisStore) EachKeyMatch(regexp string) ([]string, error) {
	return rd.FindPrefixFor(100, regexp)
}

// ScanMatche uses underline redis scan methods for a hashmap, relying on the lastIndex
// as a way to track the last cursor point on the store. Note that due to the way redis works
// works, the count is not guaranteed to stay as such, it can be ignored and more may be returned
// or less/none.
//
// With scan the order is not guaranteed.
func (rd *RedisStore) ScanMatch(count int64, lastIndex int64, _ string, regexp string) (nstorage.ScanResult, error) {
	if len(regexp) == 0 {
		regexp = ".+"
	}

	var rs nstorage.ScanResult
	var regx, rgErr = regexp2.Compile(regexp)
	if rgErr != nil {
		return rs, nerror.WrapOnly(rgErr)
	}

	var scanned = rd.Client.ZRange(rd.ctx, rd.hashZList, lastIndex, lastIndex+count-1)
	var ky, err = scanned.Result()
	if err != nil {
		return rs, nerror.WrapOnly(err)
	}

	var keys = make([]string, 0, len(ky))
	for _, item := range ky {
		var ritem = rd.unHashKey(item)
		if !regx.MatchString(ritem) {
			continue
		}
		keys = append(keys, ritem)
	}

	// rs.Finished = cursor == 0
	var lastKey string
	if keysCount := len(keys); keysCount > 0 {
		lastKey = keys[count-1]
	}

	// var isFinished = cursor == 0
	return nstorage.ScanResult{
		Finished:  false,
		Keys:      keys,
		LastIndex: lastIndex + count,
		LastKey:   lastKey,
	}, nil
}

// Count returns the total count of element in the store.
func (rd *RedisStore) Count() (int64, error) {
	var command = rd.Client.HLen(rd.ctx, rd.hashElem)

	var err = command.Err()
	if err != nil {
		return -1, nerror.WrapOnly(err)
	}

	var count = command.Val()
	return count, nil
}

// FindPrefixFor returns all matching values within store, if elements found match giving
// count then all values returned.
//
// if an error occurs, the partially collected list of keys and error is returned.
func (rd *RedisStore) FindPrefixFor(count int64, regexp string) ([]string, error) {
	if len(regexp) == 0 {
		regexp = ".+"
	}

	var regx, rgErr = regexp2.Compile(regexp)
	if rgErr != nil {
		return nil, nerror.WrapOnly(rgErr)
	}

	var cursor uint64
	var keys = make([]string, 0, count)
	var err error
	for {
		var scanned = rd.Client.SScan(rd.ctx, rd.hashList, cursor, "*", count)
		var ky, cursor, err = scanned.Result()
		if err != nil {
			return keys, nerror.WrapOnly(err)
		}

		for _, item := range ky {
			if !regx.MatchString(item) {
				continue
			}
			keys = append(keys, item)
		}

		if cursor == 0 {
			break
		}
	}

	if err != nil {
		return nil, nerror.WrapOnly(err)
	}
	return keys, nil
}

// Exists returns true/false if giving key exists.
func (rd *RedisStore) Exists(key string) (bool, error) {
	var hashKey = rd.doHashKey(key)
	var nstatus = rd.Client.SIsMember(rd.ctx, rd.hashList, hashKey)
	if err := nstatus.Err(); err != nil {
		return false, nerror.WrapOnly(err)
	}
	return nstatus.Val(), nil
}

// exists returns true/false if giving key is set in redis.
func (rd *RedisStore) exists(key string) (bool, error) {
	var hashKey = rd.doHashKey(key)
	var nstatus = rd.Client.Exists(rd.ctx, hashKey)
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
	var _, err = rd.Client.TxPipelined(rd.ctx, func(pipeliner redis.Pipeliner) error {
		var zstatus = pipeliner.ZRem(rd.ctx, rd.hashZList, items...)
		if err := zstatus.Err(); err != nil {
			return err
		}
		var mstatus = pipeliner.SRem(rd.ctx, rd.hashList, items...)
		if err := mstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		var dstatus = pipeliner.Del(rd.ctx, keys...)
		if err := dstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		return nil
	})
	if err != nil {
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
	var hashKey = rd.doHashKey(key)
	var _, pipeErr = rd.Client.TxPipelined(rd.ctx, func(pipeliner redis.Pipeliner) error {
		var nstatus = pipeliner.SAdd(rd.ctx, rd.hashList, hashKey)
		if err := nstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}

		var zs redis.Z
		zs.Score = 0
		zs.Member = hashKey

		var zstatus = pipeliner.ZAdd(rd.ctx, rd.hashZList, &zs)
		if err := zstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}

		var nset = pipeliner.Set(rd.ctx, hashKey, data, expiration)
		if err := nset.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		return nil
	})

	if err := pipeErr; err != nil {
		return nerror.WrapOnly(pipeErr)
	}

	return nil
}

// Update updates giving key with new data slice with 0 duration.
func (rd *RedisStore) Update(key string, data []byte) error {
	return rd.UpdateTTL(key, data, 0)
}

// UpdateTTL updates giving session stored with giving key. It updates
// the underline data and increases the expiration with provided value.
//
// if expiration is zero then giving value expiration will not be reset but left
// as is.
func (rd *RedisStore) UpdateTTL(key string, data []byte, expiration time.Duration) error {
	var hashKey = rd.doHashKey(key)
	var fstatus = rd.Client.SIsMember(rd.ctx, rd.hashList, hashKey)
	if err := fstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}
	if !fstatus.Val() {
		return nerror.New("key does not exist")
	}

	var _, pipeErr = rd.Client.TxPipelined(rd.ctx, func(cl redis.Pipeliner) error {
		if len(data) == 0 {
			var dstatus = cl.Del(rd.ctx, hashKey)
			if err := dstatus.Err(); err != nil {
				return err
			}

			var zs redis.Z
			zs.Score = 0
			zs.Member = hashKey
			var zstatus = cl.ZRem(rd.ctx, rd.hashZList, zs)
			if err := zstatus.Err(); err != nil {
				return err
			}

			return nil
		}

		var zs redis.Z
		zs.Score = 0
		zs.Member = hashKey

		var zstatus = cl.ZAdd(rd.ctx, rd.hashZList, &zs)
		if err := zstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}

		var nset = cl.Set(rd.ctx, hashKey, data, expiration)
		if err := nset.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		return nil
	})

	if err := pipeErr; err != nil {
		return nerror.WrapOnly(pipeErr)
	}
	return nil
}

// TTL returns current expiration time for giving key.
func (rd *RedisStore) TTL(key string) (time.Duration, error) {
	var hashKey = rd.doHashKey(key)
	var nstatus = rd.Client.PTTL(rd.ctx, hashKey)
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
	var hashKey = rd.doHashKey(key)
	var nstatus = rd.Client.PTTL(rd.ctx, hashKey)
	if err := nstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}

	if nstatus.Val() < 0 {
		return nil
	}

	var newExpiration = expiration + nstatus.Val()
	var _, pipeErr = rd.Client.TxPipelined(rd.ctx, func(cl redis.Pipeliner) error {
		if expiration == 0 {
			var exstatus = cl.Persist(rd.ctx, hashKey)
			return exstatus.Err()
		}

		var exstatus = cl.Expire(rd.ctx, hashKey, newExpiration)
		return exstatus.Err()
	})

	if err := pipeErr; err != nil {
		return nerror.WrapOnly(pipeErr)
	}

	return nil
}

// ResetTTL resets giving expiration value to provided duration.
//
// A duration of zero persists the giving key.
func (rd *RedisStore) ResetTTL(key string, expiration time.Duration) error {
	var hashKey = rd.doHashKey(key)
	var nstatus = rd.Client.PTTL(rd.ctx, hashKey)
	if err := nstatus.Err(); err != nil {
		return nerror.WrapOnly(err)
	}

	if nstatus.Val() < 0 {
		return nil
	}

	var _, pipeErr = rd.Client.TxPipelined(rd.ctx, func(cl redis.Pipeliner) error {
		if expiration == 0 {
			var exstatus = cl.Persist(rd.ctx, hashKey)
			return exstatus.Err()
		}

		var exstatus = cl.Expire(rd.ctx, hashKey, expiration)
		return exstatus.Err()
	})
	if err := pipeErr; err != nil {
		return nerror.WrapOnly(pipeErr)
	}

	return nil
}

// GetAnyKeys returns a list of values for any of the key's found.
// Unless a specific error occurred retrieving the value of a key, if a
// key is not found then it is ignored and a nil is set in it's place.
func (rd *RedisStore) GetAnyKeys(keys ...string) ([][]byte, error) {
	var modifiedKeys = make([]string, len(keys))
	for index, key := range keys {
		modifiedKeys[index] = rd.doHashKey(key)
	}

	var nstatus = rd.Client.MGet(rd.ctx, modifiedKeys...)
	if err := nstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}

	var values = make([][]byte, len(keys))
	var contentList = nstatus.Val()
	for index, val := range contentList {
		switch mv := val.(type) {
		case string:
			values[index] = nunsafe.String2Bytes(mv)
		case []byte:
			values[index] = mv
		default:
			values[index] = nil
		}
	}
	return values, nil
}

// GetAllKeys returns a list of values for any of the key's found.
// if the value of a key is not found then we stop immediately, returning
// an error and the current set of items retreived.
func (rd *RedisStore) GetAllKeys(keys ...string) ([][]byte, error) {
	var modifiedKeys = make([]string, len(keys))
	for index, key := range keys {
		modifiedKeys[index] = rd.doHashKey(key)
	}

	var nstatus = rd.Client.MGet(rd.ctx, modifiedKeys...)
	if err := nstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}

	var values = make([][]byte, len(keys))
	var contentList = nstatus.Val()
	for index, val := range contentList {
		switch mv := val.(type) {
		case string:
			values[index] = nunsafe.String2Bytes(mv)
		case []byte:
			values[index] = mv
		default:
			return values, nerror.New("value with type %T has value %#v but is not bytes or string for key %q", mv, mv, keys[index])
		}
	}
	return values, nil
}

// Get returns giving session stored with giving key, returning an
// error if not found.
func (rd *RedisStore) Get(key string) ([]byte, error) {
	var hashKey = rd.doHashKey(key)
	var nstatus = rd.Client.Get(rd.ctx, hashKey)
	if err := nstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}
	return nunsafe.String2Bytes(nstatus.Val()), nil
}

// RemoveKeys removes underline key from the redis store after retrieving it and
// returning giving session.
func (rd *RedisStore) RemoveKeys(keys ...string) error {
	var modifiedKeys = make([]string, len(keys))
	var modifiedIKeys = make([]interface{}, len(keys))

	for index, key := range keys {
		var mod = rd.doHashKey(key)
		modifiedKeys[index] = mod
		modifiedIKeys[index] = mod
	}

	var _, err = rd.Client.TxPipelined(rd.ctx, func(pipeliner redis.Pipeliner) error {
		var zstatus = pipeliner.ZRem(rd.ctx, rd.hashZList, modifiedIKeys...)
		if err := zstatus.Err(); err != nil {
			return err
		}
		var mstatus = pipeliner.SRem(rd.ctx, rd.hashList, modifiedIKeys...)
		if err := mstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		var dstatus = pipeliner.Del(rd.ctx, modifiedKeys...)
		if err := dstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		return nil
	})
	if err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

// Remove removes underline key from the redis store after retrieving it and
// returning giving session.
func (rd *RedisStore) Remove(key string) ([]byte, error) {
	var hashKey = rd.doHashKey(key)
	var nstatus = rd.Client.Get(rd.ctx, hashKey)
	if err := nstatus.Err(); err != nil {
		return nil, nerror.WrapOnly(err)
	}

	var _, err = rd.Client.TxPipelined(rd.ctx, func(pipeliner redis.Pipeliner) error {
		var zstatus = pipeliner.ZRem(rd.ctx, rd.hashZList, hashKey)
		if err := zstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		var mstatus = pipeliner.SRem(rd.ctx, rd.hashList, hashKey)
		if err := mstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		var dstatus = pipeliner.Del(rd.ctx, hashKey)
		if err := dstatus.Err(); err != nil {
			return nerror.WrapOnly(err)
		}
		return nil
	})
	if err != nil {
		return nil, nerror.WrapOnly(err)
	}
	return nunsafe.String2Bytes(nstatus.Val()), nil
}
