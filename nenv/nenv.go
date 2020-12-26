package nenv

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/influx6/npkg/npair"
)

type Provider interface {
	Add(key string, value interface{})
}

// Loader defines a loader function which returns all
// associated environment functions as a slice of
// key=value pairs.
type Loader interface {
	Register(Provider) error
}

var _ Loader = (*EnvironmentLoader)(nil)

type ListLoader []string

func (el ListLoader) Register(p Provider) error {
	for _, item := range el {
		var eqIndex = strings.Index(item, "=")
		if eqIndex == -1 {
			continue
		}
		var key = strings.ToLower(item[:eqIndex])
		var value = item[eqIndex+1:]
		p.Add(key, value)
	}
	return nil
}

type EnvironmentLoader struct{}

func (el *EnvironmentLoader) Register(p Provider) error {
	var envMap = os.Environ()
	for _, item := range envMap {
		var eqIndex = strings.Index(item, "=")
		if eqIndex == -1 {
			continue
		}
		var key = strings.ToLower(item[:eqIndex])
		var value = item[eqIndex+1:]
		p.Add(key, value)
	}
	return nil
}

func LoadEnvironment(prefix string) (*EnvStore, error) {
	var loader EnvironmentLoader
	return LoadFrom(prefix, &loader)
}

func LoadFrom(prefix string, loader Loader) (*EnvStore, error) {
	var store = New(prefix)
	return store, loader.Register(store)
}

// EnvStore defines an implementation which during initial
// loading stores all environment keys and values.
//
// All environment keys are lowercased to avoid to ensure
// case insensitivity, so be aware when using such for keys.
//
// It provides a means of reloading giving values
// storing them efficient for retrieval and change.
type EnvStore struct {
	prefix string
	fl     sync.RWMutex
	cache  npair.Field
}

// New returns a new map which contains all environment variables.
//
// If Prefix is provided then only environment variables with giving prefix
// are loaded. Prefixed are expected to be in the format: PREFIX_ENV_VAR.
func New(prefix string) *EnvStore {
	var env EnvStore
	env.cache = npair.Field{}
	if len(prefix) > 0 {
		env.prefix = strings.ToLower(strings.TrimSuffix(prefix, "_"))
	}
	return &env
}

func (env *EnvStore) Add(key string, value interface{}) {
	key = strings.ToLower(key)
	if len(env.prefix) > 0 && !strings.HasPrefix(key, env.prefix+"_") {
		return
	}
	env.fl.Lock()
	env.cache[key] = value
	env.fl.Unlock()
}

// Has returns a true/false if giving key value exists.
func (env *EnvStore) Has(key string) bool {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.Has(getHash(env.prefix, key))
}

// Set sets provided key within store and sets the giving environment
// variable in to the environment.
func (env *EnvStore) Set(key string, value string) {
	env.Add(getHash(env.prefix, key), value)
}

// GetBool collects the string value of a key if it exists.
func (env *EnvStore) GetBool(key string) (bool, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetBool(getHash(env.prefix, key))
}

// GetFloat64 collects the string value of a key if it exists.
func (env *EnvStore) GetFloat64(key string) (float64, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetFloat64(getHash(env.prefix, key))
}

// GetFloat32 collects the string value of a key if it exists.
func (env *EnvStore) GetFloat32(key string) (float32, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetFloat32(getHash(env.prefix, key))
}

// GetInt8 collects the string value of a key if it exists.
func (env *EnvStore) GetInt8(key string) (int8, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetInt8(getHash(env.prefix, key))
}

// GetInt16 collects the string value of a key if it exists.
func (env *EnvStore) GetInt16(key string) (int16, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetInt16(getHash(env.prefix, key))
}

// GetInt64 collects the string value of a key if it exists.
func (env *EnvStore) GetInt64(key string) (int64, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetInt64(getHash(env.prefix, key))
}

// GetInt32 collects the string value of a key if it exists.
func (env *EnvStore) GetInt32(key string) (int32, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetInt32(getHash(env.prefix, key))
}

// GetInt collects the string value of a key if it exists.
func (env *EnvStore) GetInt(key string) (int, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetInt(getHash(env.prefix, key))
}

// KeyFor returns the expected name of a giving key from the environment
// store.
func (env *EnvStore) KeyFor(key string) string {
	return getHash(env.prefix, key)
}

// GetString collects the string value of a key if it exists.
func (env *EnvStore) GetString(key string) (string, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetString(getHash(env.prefix, key))
}

// GetDuration collects the string value of a key if it exists.
func (env *EnvStore) GetDuration(key string) (time.Duration, bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.GetDuration(getHash(env.prefix, key))
}

// Get collects the value of a key if it exists.
func (env *EnvStore) Get(key string) (value interface{}, found bool) {
	env.fl.RLock()
	defer env.fl.RUnlock()
	return env.cache.Get(getHash(env.prefix, key))
}

// getHash returns giving format of key.
func getHash(prefix string, key string) string {
	if len(prefix) == 0 {
		return strings.ToLower(key)
	}
	key = strings.ToLower(key)
	key = strings.TrimSuffix(key, "_")
	return strings.ToLower(prefix + "_" + key)
}
