package datautils

import "testing"

func TestCompressDecompressError(t *testing.T) {
	data := "compression"

	_, err := Compress([]byte(data), 4)
	if err == nil || err.Error() != "compression type not found" {
		t.Error("It should return err")
	}

	_, err = Decompress([]byte("err"), 4)
	if err == nil || err.Error() != "compression type not found" {
		t.Error("It should return err")
	}
}

func TestCompressDecompressGZip(t *testing.T) {
	data := "compression gzip"

	compressed, err := Compress([]byte(data), GZip)
	if err != nil {
		t.Error("err should be nil")
	}

	decompressed, err := Decompress(compressed, GZip)
	if err != nil {
		t.Error("err should be nil")
	}

	if string(decompressed) != data {
		t.Error("It should be equal")
	}
}

func TestCompressDecompressZLib(t *testing.T) {
	data := "compression zlib"

	compressed, err := Compress([]byte(data), Zlib)
	if err != nil {
		t.Error("err should be nil")
	}

	decompressed, err := Decompress(compressed, Zlib)
	if err != nil {
		t.Error("err should be nil")
	}

	if string(decompressed) != data {
		t.Error("It should be equal")
	}
}
