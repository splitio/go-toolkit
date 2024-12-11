package logging

import (
	"context"
	"fmt"
	"slices"
)

// ContextKey =
type ContextKey struct{}

// ContextData that will store all the context data to display in log
type ContextData struct {
	m map[string]string
}

// NewContext creates a new ContextData object
func NewContext() *ContextData {
	return &ContextData{m: map[string]string{}}
}

// WithTag adds a key value pair to the context
func (c *ContextData) WithTag(key string, value string) *ContextData {
	c.m[key] = value
	return c
}

// Get returns the value of the key
func (c *ContextData) Get(key string) string {
	return c.m[key]
}

// String returns the string representation of the context
func (c *ContextData) String() string {
	keys := make([]string, 0, len(c.m))
	for k := range c.m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	toReturn := "["
	for _, k := range keys {
		toReturn += fmt.Sprintf("%s=%s, ", k, c.m[k])
	}
	toReturn = toReturn[:len(toReturn)-2]
	toReturn += "]"
	return toReturn
}

// Merge merges two context information objects
func Merge(one *ContextData, two *ContextData) *ContextData {
	if one == nil && two == nil {
		return nil
	}
	merged := NewContext()
	if one != nil {
		for k, v := range one.m {
			merged = merged.WithTag(k, v)
		}
	}
	if two != nil {
		for k, v := range two.m {
			merged = merged.WithTag(k, v)
		}
	}

	return merged
}

// GetContext returns the context information object from the context
func GetContext(ctx context.Context) *ContextData {
	ci, ok := ctx.Value(ContextKey{}).(*ContextData)
	if !ok {
		return nil
	}
	return ci
}
