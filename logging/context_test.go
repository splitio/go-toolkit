package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextInformation(t *testing.T) {
	one := NewContextInformation()
	one.Add("key", "value")
	assert.Equal(t, one.Get("key"), "value")
	assert.Empty(t, one.Get("key2"))
	assert.Equal(t, one.String(), "[key: value]")
	assert.Nil(t, Merge(nil, nil))
	assert.Equal(t, Merge(one, nil).String(), "[key: value]")
	two := NewContextInformation()
	two.Add("key2", "value2")
	assert.Equal(t, Merge(nil, two).String(), "[key2: value2]")
	assert.Equal(t, Merge(one, two).String(), "[key: value, key2: value2]")
}
