package datautils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	_, err := Encode([]byte("err"), 4)
	assert.ErrorContains(t, err, "encode type not found")

	_, err = Decode("err", 4)
	assert.ErrorContains(t, err, "encode type not found")
}

func TestB64EncodeDecode(t *testing.T) {
	data := "encode b64"
	encoded, err := Encode([]byte(data), Base64)
	assert.Nil(t, err)

	decoded, err := Decode(encoded, Base64)
	assert.Nil(t, err)
	assert.Equal(t, data, string(decoded))
}
