//

package nredis_test

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"

	"github.com/gokit/npkg/nstorage/internal/tharness"
	"github.com/gokit/npkg/nstorage/nredis"
)

func TestIntegrationRedisStore(t *testing.T) {
	var ops redis.Options
	ops.Addr = "localhost:6379"
	ops.Network = "tcp"
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	var store, err = nredis.FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestsIntegrationRedisStoreExpirables(t *testing.T) {
	var ops redis.Options
	ops.Addr = "localhost:6379"
	ops.Network = "tcp"
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	var store, err = nredis.FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestExpirableStore(t, store)
}
