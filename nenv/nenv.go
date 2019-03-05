package nenv

import (
	"sync"
	"sync/atomic"
)

// EnvStore defines an implementation which during initial
// loading stores all environment keys and values.
//
// It provides a means of reloading giving values
// storing them efficient for retrieval and change.
type EnvStore struct {
	lock      sync.Mutex
	cache     atomic.Value
	userCache atomic.Value
}

func copyMap(src map[string]string) map[string]string {
	var dest = make(map[string]string, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

func copyISMap(src map[string]interface{}) map[string]interface{} {
	var dest = make(map[string]interface{}, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

func copyIMap(src map[interface{}]interface{}) map[interface{}]interface{} {
	var dest = make(map[interface{}]interface{}, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}
