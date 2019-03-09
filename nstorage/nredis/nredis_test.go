// +build unit

package nredis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"

	"github.com/gokit/npkg/nstorage/internal/tharness"
)

func TestStoreWithRedisStore(t *testing.T) {
	var server, err = miniredis.Run()
	require.NoError(t, err)
	require.NotNil(t, server)

	defer server.Close()

	var ops redis.Options
	ops.Addr = server.Addr()
	ops.Network = "tcp"
	require.NotNil(t, &ops)
	require.Equal(t, server.Addr(), ops.Addr)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	var store *RedisStore
	store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestsRedisExpiryStore(t *testing.T) {
	var server, err = miniredis.Run()
	require.NoError(t, err)
	require.NotNil(t, server)

	defer server.Close()

	var ops redis.Options
	ops.Addr = server.Addr()
	ops.Network = "tcp"
	require.NotNil(t, &ops)
	require.Equal(t, server.Addr(), ops.Addr)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	var store *RedisStore
	store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestExpirableStore(t, store)
}
