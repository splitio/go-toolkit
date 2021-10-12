package datautils

import (
	"encoding/base64"
	"fmt"
)

const (
	Base64 = iota
)

func Encode(data []byte, encodeType int) (string, error) {
	switch encodeType {
	case Base64:
		return base64.StdEncoding.EncodeToString(data), nil
	}
	return "", fmt.Errorf("encode type not found")

}

func Decode(data string, decodeType int) ([]byte, error) {
	switch decodeType {
	case Base64:
		return base64.StdEncoding.DecodeString(data)
	}
	return []byte{}, fmt.Errorf("encode type not found")
}
