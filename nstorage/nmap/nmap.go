package nmap

import (
	"sync"
	"sync/atomic"
)

// MapStore defines an implementation which during initial
// loading stores all key and value pairs.
//
// It provides a safe, concurrently usable implementation with
// blazing read and write speed.
type MapStore struct {
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
