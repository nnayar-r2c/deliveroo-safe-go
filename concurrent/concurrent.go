package concurrent

import (
	"sync"
)

// errgroup abstracts errgroup-like structs such as `errgroup.Group` or `safe.Group`
type errgroup interface {
	Go(func() error)
}

// Future is a struct holding a generic future value.
// The value can only be retrieved using Get, or MustGet functions.
type Future[T any] struct {
	value T
	err   error
	mu    sync.RWMutex
}

// FetchFn is a function that the client wishes to be called asynchronously.
// It should return an arbitrary type and an error if the data fetch failed.
type FetchFn[T any] func() (T, error)

// Fetch accepts an errgroup and a function that returns an arbitrary type.
// The function will be called asynchronously and a Future of the result will be returned immediately.
// The Future holds a future value, that the clients can retrieve using a blocking Get function.
// The function will be executed using the provided errgroup, so that the clients can Wait for a group of Futures.
// See unit tests for an example of usage.
func Fetch[T any](wg errgroup, fn FetchFn[T]) *Future[T] {
	var p Future[T]
	p.mu.Lock()

	wg.Go(func() error {
		defer p.mu.Unlock()

		value, err := fn()
		p.value = value
		p.err = err

		return err
	})

	return &p
}

// Get blocks until the future result is available
// and returns the result.
func (f *Future[T]) Get() (T, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.value, f.err
}

// MustGet is the same as Get but panics on errors
// It should be safe to call if Wait on the provided errgroup returned without errors
func (f *Future[T]) MustGet() T {
	val, err := f.Get()
	if err != nil {
		panic(err)
	}

	return val
}
