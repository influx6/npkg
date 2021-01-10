package nredis

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/influx6/npkg/nstorage/internal/tharness"
	"github.com/stretchr/testify/require"
)

func TestIntegrationRedisStoreGetAnyKeys(t *testing.T) {
	var ops redis.Options
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	if err := redisClient.Ping().Err(); err != nil {
		t.SkipNow()
		return
	}

	var store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreGetAnykeys(t, store)
}

func TestIntegrationRedisStoreGetAllKeys(t *testing.T) {
	var ops redis.Options
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	if err := redisClient.Ping().Err(); err != nil {
		t.SkipNow()
		return
	}

	var store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreGetAllkeys(t, store)
}

func TestIntegrationRedisStoreFindEach(t *testing.T) {
	var ops redis.Options
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	if err := redisClient.Ping().Err(); err != nil {
		t.SkipNow()
		return
	}

	var store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreFindEach(t, store)
}

func TestIntegrationRedisStoreEachPrefixKey(t *testing.T) {
	var ops redis.Options
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	if err := redisClient.Ping().Err(); err != nil {
		t.SkipNow()
		return
	}

	var store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreFindPrefix(t, store)
}

func TestIntegrationRedisStoreFindAll(t *testing.T) {
	var ops redis.Options
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	if err := redisClient.Ping().Err(); err != nil {
		t.SkipNow()
		return
	}

	var store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreFindAll(t, store)
}

func TestIntegrationRedisStore(t *testing.T) {
	var ops redis.Options
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	if err := redisClient.Ping().Err(); err != nil {
		t.SkipNow()
		return
	}

	var store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestIntegrationRedisStoreExpirables(t *testing.T) {
	var ops redis.Options
	require.NotNil(t, &ops)

	var redisClient = redis.NewClient(&ops)
	require.NotNil(t, redisClient)

	if err := redisClient.Ping().Err(); err != nil {
		t.SkipNow()
		return
	}

	var store, err = FromRedisStore("testing_mb", redisClient)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestExpirableStore(t, store)
}
