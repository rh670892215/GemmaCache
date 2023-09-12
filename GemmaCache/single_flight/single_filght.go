package single_flight

import (
	"sync"
)

type Caller struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type ConcurrencyLimiter struct {
	mutex   sync.Mutex
	callers map[string]*Caller
}

func (c *ConcurrencyLimiter) Do(key string, f func() (interface{}, error)) (interface{}, error) {
	c.mutex.Lock()
	if c.callers == nil {
		c.callers = make(map[string]*Caller)
	}

	caller, ok := c.callers[key]
	if ok {
		c.mutex.Unlock()
		caller.wg.Wait()
		return caller.val, caller.err
	}

	caller = &Caller{}
	caller.wg.Add(1)
	c.callers[key] = caller
	c.mutex.Unlock()

	caller.val, caller.err = f()
	caller.wg.Done()

	c.mutex.Lock()
	delete(c.callers, key)
	c.mutex.Unlock()
	return caller.val, caller.err
}
