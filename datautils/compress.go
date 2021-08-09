package datautils

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io/ioutil"
)

const (
	None = iota
	GZip
	Zlib
)

func Compress(data []byte, compressType int) ([]byte, error) {
	var b bytes.Buffer
	switch compressType {
	case GZip:
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write(data); err != nil {
			return nil, err
		}
		if err := gz.Close(); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case Zlib:
		zl := zlib.NewWriter(&b)
		if _, err := zl.Write(data); err != nil {
			return nil, err
		}
		if err := zl.Close(); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}
	return nil, fmt.Errorf("compression type not found")
}

func Decompress(data []byte, compressType int) ([]byte, error) {
	b := bytes.NewReader(data)
	switch compressType {
	case GZip:
		gz, err := gzip.NewReader(b)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		raw, err := ioutil.ReadAll(gz)
		if err != nil {
			return nil, err
		}
		return raw, nil
	case Zlib:
		zl, err := zlib.NewReader(b)
		if err != nil {
			return nil, err
		}
		defer zl.Close()
		raw, err := ioutil.ReadAll(zl)
		if err != nil {
			return nil, err
		}
		return raw, nil
	}
	return nil, fmt.Errorf("compression type not found")
}
