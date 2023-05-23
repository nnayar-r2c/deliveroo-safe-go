package concurrent

import (
	"fmt"
	"testing"

	"github.com/deliveroo/assert-go"
	"github.com/deliveroo/safe-go"
)

type testerrgroup struct {
	result error
}

func (t *testerrgroup) Go(fn func() error) {
	// make the call synchronous for tests
	_ = fn()
}

func TestFetch(t *testing.T) {
	t.Run("demonstration", func(t *testing.T) {
		// initialize a waitgroup so that the client can wait for all operation results
		var wg = &safe.Group{}

		// this fetch represents an expensive call (e.g. an external request for a string)
		stringFuture := Fetch(wg, func() (string, error) {
			// this function will be called asynchronously
			return "hello there", nil
		})

		// this fetch represents another expensive call
		intFuture := Fetch(wg, func() (int, error) {
			return 3, nil
		})

		// the Futures can be passed to different threads
		go func() {
			// this call will block until the result is available
			str, err := stringFuture.Get()
			assert.Nil(t, err)
			assert.Equal(t, "hello there", str)
		}()

		// client can call wg.Wait to ensure all operations complete before continuing
		err := wg.Wait()
		assert.Nil(t, err)

		// it is now safe to call MustGet assuming no errors were returned by the errgroup
		// otherwise MustGet can panic
		val := intFuture.MustGet()
		assert.Equal(t, val, 3)
	})

	t.Run("nil value", func(t *testing.T) {
		var wg testerrgroup
		p := Fetch(&wg, func() (*string, error) {
			return nil, nil
		})
		assert.NotNil(t, p)
		assert.Nil(t, p.value)
		assert.Nil(t, p.err)

		val, err := p.Get()
		assert.Nil(t, val)
		assert.Nil(t, err)

		val = p.MustGet()
		assert.Nil(t, val)
		assert.Nil(t, wg.result)
	})

	t.Run("not nil value", func(t *testing.T) {
		var wg testerrgroup
		var testval = "hello there"
		p := Fetch(&wg, func() (*string, error) {
			return &testval, nil
		})
		assert.NotNil(t, p)
		assert.Equal(t, &testval, p.value)
		assert.Nil(t, p.err)

		val, err := p.Get()
		assert.Equal(t, &testval, val)
		assert.Nil(t, err)

		val = p.MustGet()
		assert.Equal(t, &testval, val)
		assert.Nil(t, wg.result)
	})

	t.Run("not nil value", func(t *testing.T) {
		var wg testerrgroup
		var testerr = fmt.Errorf("unknown error")
		p := Fetch(&wg, func() (*string, error) {
			return nil, testerr
		})
		assert.NotNil(t, p)
		assert.Nil(t, p.value)
		assert.Equal(t, testerr, p.err)

		val, err := p.Get()
		assert.Nil(t, val)
		assert.Equal(t, testerr, err)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected the code to panic")
			}
		}()
		val = p.MustGet()
		assert.Nil(t, val)
		assert.Equal(t, testerr, wg.result)
	})
}
