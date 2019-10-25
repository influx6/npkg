package nredis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"

	"github.com/influx6/npkg/nstorage/internal/tharness"
)

func TestRedisStore(t *testing.T) {
	var server = miniredis.NewMiniRedis()
	require.NotNil(t, server)

	var err = server.StartAddr("localhost:0")
	require.NoError(t, err)

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
