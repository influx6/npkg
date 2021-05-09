package nthen

type Future struct {
	err error
	val interface{}
	resolved chan interface{}
}

func NewFuture() *Future {
	return &Future{
		resolved:   make(chan interface{}, 0),
	}
}

func Fn(fn func(ft *Future)) *Future {
	var ft = NewFuture()
	go fn(ft)
	return ft
}

func From(fn func() (interface{}, error)) *Future {
	var ft = NewFuture()
	go func() {
		var val, err = fn()
		if err != nil {
			ft.WithError(err)
			return
		}
		ft.WithValue(val)
	}()
	return ft
}


// CollectFor collects all resolved results (value or error)
// as a list which is used to resolve the returned future.
func CollectFor(fts ...*Future) *Future {
	var ft = NewFuture()
	go func(targets []*Future) {
		var results = make([]interface{}, len(targets))
		for index, work := range targets {
			work.Wait()
			results[index] = work.Value()
		}
		 ft.WithValue(results)
	}(fts)
	return ft
}

// WaitFor for collects all resolved successfully values
// and returns as a list of values but resolves the
// future with an error and a partial list if
// any of the futures resolved with an error.
func WaitFor(fts ...*Future) *Future {
	var ft = NewFuture()
	go func(targets []*Future) {
		var results = make([]interface{}, 0, len(targets))
		for _, work := range targets {
			work.Wait()
			var value, err = work.Get()
			if value != nil {
				results = append(results, value)
			}
			if err != nil {
				 ft.WithValueAndError(results, err)
				break
			}
		}
		 ft.WithValue(results)
	}(fts)
	return ft
}

// Then will async-ly wait in a goroutine for the target
// future to be resolved at which the provided next Future will
// be resolved as well.
func (f *Future) Then(next *Future)  {
	go func(){
		f.WaitThen(next)
	}()
}

// WaitThen will block till this future resolves, at which it
// resolves the next future provided as an argument.
func (f *Future) WaitThen(next *Future)   {
	f.Wait()
	if f.val != nil {
		next.WithValue(f.val)
	}
	next.WithError(f.err)
}

// Wait blocks till future is resolved.
func (f *Future) Wait()  {
	<-f.resolved
}

// IsResolved returns true/false if future is resolved either as error or value.
func (f *Future) IsResolved() bool {
	select {
	case <-f.resolved:
		return true
	default:
		return false
	}
}

// Get returns result of future as a tuple of value and error.
func (f *Future) Get()  (interface{}, error) {
	f.Wait()
	return f.val, f.err
}

// Value returns the value of the future, it blocks if the future is yet
// resolved and will wait till it's resolved.
func (f *Future) Value()  interface{} {
	f.Wait()
	return f.val
}

// Err returns the error of the future, it blocks if the future is yet
// resolved and will wait till it's resolved.
func (f *Future) Err() error {
	f.Wait()
	return f.err
}

// WithError resolves this Future as a failed operation with provided
// error.
func (f *Future) WithError(v error)  {
	select {
	case <-f.resolved:
		return
	default:
		f.err = v
		close(f.resolved)
	}
	return
}

// WithValueAndError resolves this Future as a with a value for
// both result and error. Useful for operations where a
// value is returned but so was an error.
func (f *Future) WithValueAndError(v interface{}, err error)  {
	select {
	case <-f.resolved:
		return
	default:
		f.val = v
		f.err = err
		close(f.resolved)
	}
	return
}


// WithValue resolves this Future as a completed operation with provided
// value.
func (f *Future) WithValue(v interface{})  {
	select {
	case <-f.resolved:
		return
	default:
		f.val = v
		close(f.resolved)
	}
	return
}

