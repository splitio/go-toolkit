package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextData(t *testing.T) {
	one := NewContext().WithTag("key", "value")
	assert.Equal(t, one.Get("key"), "value")
	assert.Empty(t, one.Get("key2"))
	assert.Equal(t, one.String(), "[key=value]")
	assert.Nil(t, Merge(nil, nil))
	assert.Equal(t, Merge(one, nil).String(), "[key=value]")
	two := NewContext().WithTag("key2", "value2")
	assert.Equal(t, Merge(nil, two).String(), "[key2=value2]")
	assert.Equal(t, Merge(one, two).String(), "[key=value, key2=value2]")
}
