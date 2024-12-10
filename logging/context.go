package logging

import (
	"context"
	"fmt"
	"slices"
)

// ContextKey =
type ContextKey struct{}

// ContextInformation that will store all the context data to display in log
type ContextInformation struct {
	m map[string]string
}

// NewContextInformation creates a new ContextInformation object
func NewContextInformation() *ContextInformation {
	return &ContextInformation{m: map[string]string{}}
}

// Add adds a key value pair to the context
func (c *ContextInformation) Add(key, value string) {
	c.m[key] = value
}

// Get returns the value of the key
func (c *ContextInformation) Get(key string) string {
	return c.m[key]
}

// String returns the string representation of the context
func (c *ContextInformation) String() string {
	keys := make([]string, 0, len(c.m))
	for k := range c.m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	toReturn := "["
	for _, k := range keys {
		toReturn += fmt.Sprintf("%s: %s, ", k, c.m[k])
	}
	toReturn = toReturn[:len(toReturn)-2]
	toReturn += "]"
	return toReturn
}

// Merge merges two context information objects
func Merge(one *ContextInformation, two *ContextInformation) *ContextInformation {
	if one == nil && two == nil {
		return nil
	}
	merged := NewContextInformation()
	if one != nil {
		for k, v := range one.m {
			merged.Add(k, v)
		}
	}
	if two != nil {
		for k, v := range two.m {
			merged.Add(k, v)
		}
	}

	return merged
}

// GetContext returns the context information object from the context
func GetContext(ctx context.Context) *ContextInformation {
	lc := ctx.Value(ContextKey{})
	if lc == nil {
		return nil
	}
	ci, ok := lc.(*ContextInformation)
	if !ok {
		return nil
	}
	return ci
}
