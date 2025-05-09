package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	hash, err := Encode(nil, "something")
    assert.ErrorContains(t, err, "Hasher could not be nil")
    assert.Equal(t, "", hash)
    
	hasher := NewMurmur332Hasher(0)
	hash2, err := Encode(hasher, "something")
    assert.Nil(t, err)
    assert.Equal(t, "NDE0MTg0MjI2MQ==", hash2)
}
