package datautils

import (
    "testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressDecompressError(t *testing.T) {
	data := "compression"

	_, err := Compress([]byte(data), 4)
    assert.ErrorContains(t, err, "compression type not found")

	_, err = Decompress([]byte("err"), 4)
    assert.ErrorContains(t, err, "compression type not found")
}

func TestCompressDecompressGZip(t *testing.T) {
	data := "compression gzip"

	compressed, err := Compress([]byte(data), GZip)
    assert.Nil(t, err)

	decompressed, err := Decompress(compressed, GZip)
    assert.Nil(t, err)

    assert.Equal(t, data, string(decompressed))
}

func TestCompressDecompressZLib(t *testing.T) {
	data := "compression zlib"

	compressed, err := Compress([]byte(data), Zlib)
    assert.Nil(t, err)

	decompressed, err := Decompress(compressed, Zlib)
    assert.Nil(t, err)

    assert.Equal(t, data, string(decompressed))
}
