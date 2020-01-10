package cache

import (
	"fmt"
	"time"
)

// Miss is a special type of error indicating that a key was not found
type Miss struct {
	Where string
	Key   string
}

func (c *Miss) Error() string {
	return fmt.Sprintf("key %s not found in cache %s", c.Key, c.Where)
}

// Expired is a special type of error indicating that a key is no longer valid (value is still attached in the error)
type Expired struct {
	Key   string
	When  time.Time
	Value interface{}
}

func (e *Expired) Error() string {
	return fmt.Sprintf("Key %s is expired.", e.Key)
}
