package nstream

// Publisher represents a giving stream producer, which has giving
// elements it can produce to all giving subscribers. Usually
// it is advised to have a Publisher only accept a single subscriber
// to avoid management details, but this will not be the case for
// Producers who may fan-out elements to multiple subscribers.
//
// Producers can also be called Publishers.
type Publisher interface {
	// Subscribe takes a giving Subscriber which has the intention of
	// listening elements produced by said Publisher. Subscribers will
	// be provided a Subscription object which can be used to request
	// elements from producers, this allows back pressure mitigation
	// techniques and Subscriber based pulling of data which allows
	// efficient management of resources.
	Subscribe(Subscriber) error
}

// BytePublisher represents a giving byte stream producer, which has giving
// elements it can produce to all giving subscribers. Usually
// it is advised to have a Publisher only accept a single subscriber
// to avoid management details, but this will not be the case for
// Producers who may fan-out elements to multiple subscribers.
//
// Producers can also be called Publishers.
type BytePublisher interface {
	Subscribe(ByteSubscriber) error
}

// Subscriber defines a process interested within a giving stream,
// it receives a subscription once and then continuously calls
// for elements until completion based on it's pace.
type Subscriber interface {
	// OnError is called when a an unrecoverable error occurs during
	// the delivery of giving stream, this can usually mean the tear
	// down of giving subscription by either the subscriber or producer.
	// Usually such a choice depends on implementation details.
	OnError(error)

	// OnCompletion is called when Publisher has completed sending
	// all data elements to subscriber.
	// It may optionally take a giving completion value.
	OnCompletion(interface{})

	// OnNext is called with the next received element requested  by
	// subscriber arrives.
	OnNext(interface{}) error

	// OnSubscription is only ever called once with provided subscription.
	// The subscriber will use said subscription for reading nstream to completion
	// or error.
	OnSubscription(Subscription)
}

// ByteSubscriber defines a process interested within a giving stream,
// it receives a subscription once and then continuously calls
// for elements until completion based on it's pace.
type ByteSubscriber interface {
	// OnError is called when a an unrecoverable error occurs during
	// the delivery of giving stream, this can usually mean the tear
	// down of giving subscription by either the subscriber or producer.
	// Usually such a choice depends on implementation details.
	OnError(error)
	
	// OnCompletion is called when Publisher has completed sending
	// all data elements to subscriber.
	// It may optionally take a giving completion value.
	OnCompletion([]byte)
	
	// OnNext is called with the next received element requested  by
	// subscriber arrives.
	OnNext([]byte)
	
	// OnSubscription is only ever called once with provided subscription.
	// The subscriber will use said subscription for reading nstream to completion
	// or error.
	OnSubscription(Subscription)
}

// Subscription represents a agreed subscription between a producer
// and Subscriber. It is used by subscriber to request more elements
// from giving Publisher.
type Subscription interface {
	// Next request more n elements form underline Publisher, it indicates
	// to producer the desire of Subscriber to be able to take n giving elements
	// and is the only means of flow.
	//
	// It returns an error if it fails to deliver next batch request to publisher.
	Next(int) error

	// Stop ends giving Subscription, restoring all resources used up by it's
	// implementation. An optional error can be supplied to indicate why.
	Stop(error)
}

// Duplexer are the both subscribers and publishers of streams.
type Duplexer interface {
	Publisher
	Subscriber
}

// ByteDuplexer are the both subscribers and publishers of byte streams.
type ByteDuplexer interface {
	BytePublisher
	ByteSubscriber
}
