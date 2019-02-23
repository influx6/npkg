package nstream

import "sync"

// NextSource source provides what a OnePublisher expects
// as a possible source for providing subscription to.
type NextSource interface {
	More() bool
	Next() (interface{}, error)
}

// OnePublisher implements the Publisher interface and provides
// subscription for one single subscriber for delivery of reactive
// streams of data.
type OnePublisher struct {
	source     NextSource
	subm       sync.Mutex
	subscriber Subscriber
}
