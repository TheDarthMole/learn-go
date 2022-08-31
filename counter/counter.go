package counter

import "sync"

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}

// NewCounter We want to always initialise a new counter
// to use a pointer rather than to duplicate the value
// when being used as an argument as this would be unsafe
func NewCounter() *Counter {
	return &Counter{}
}
