package nenv

import (
	"os"
	"strings"

	"github.com/influx6/npkg/nstorage/nmap"
)

// Loader defines a loader function which returns all
// associated environment functions as a slice of
// key=value pairs.
type Loader func() ([]string, error)

// EnvStore defines an implementation which during initial
// loading stores all environment keys and values.
//
// All environment keys are lowercased to avoid to ensure
// case insensitivity, so be aware when using such for keys.
//
// It provides a means of reloading giving values
// storing them efficient for retrieval and change.
type EnvStore struct {
	nprefix string
	prefix  string
	cache   *nmap.StringMap
}

// New returns a new map which contains all environment variables.
//
// If Prefix is provided then only environment variables with giving prefix
// are loaded. Prefixed are expected to be in the format: PREFIX_ENV_VAR.
func New(prefix string, loader Loader) (*EnvStore, error) {
	var environs, err = loader()
	if err != nil {
		return nil, err
	}

	var env EnvStore
	env.prefix = strings.ToLower(strings.TrimSuffix(prefix, "_"))
	env.nprefix = env.prefix + "_"
	env.cache = nmap.NewStringMap(uint(len(environs)))
	env.cache.SetMany(func(values map[string]string) {
		for _, item := range environs {
			var eqIndex = strings.Index(item, "=")
			if eqIndex == -1 {
				continue
			}
			var key = strings.ToLower(item[:eqIndex])
			var value = item[eqIndex+1:]
			if prefix != "" && !strings.HasPrefix(key, env.nprefix) {
				continue
			}
			values[key] = value
		}
	})
	return &env, nil
}

// Has returns a true/false if giving key value exists.
func (env *EnvStore) Has(key string) bool {
	return env.cache.Has(env.getHash(key))
}

// Set sets provided key within store and sets the giving environment
// variable in to the environment.
func (env *EnvStore) Set(key string, value string) {
	os.Setenv(env.getNHash(key), value)
	env.cache.Set(env.getHash(key), value)
}

// Get returns a giving key value.
func (env *EnvStore) Get(key string) string {
	return env.cache.Get(env.getHash(key))
}

// getHash returns giving format of key.
func (env *EnvStore) getHash(key string) string {
	key = strings.ToLower(key)
	key = strings.TrimSuffix(key, "_")
	if env.prefix != "" {
		key = env.nprefix + key
	}
	return key
}

// getNHash returns giving format of key as should be set in environment.
func (env *EnvStore) getNHash(key string) string {
	key = strings.ToUpper(key)
	key = strings.TrimSuffix(key, "_")
	if env.prefix != "" {
		key = env.nprefix + key
	}
	return key
}
