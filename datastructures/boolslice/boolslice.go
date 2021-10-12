package boolslice

import (
	"errors"
)

// ErrorOutOfBounds err
var ErrorOutOfBounds error = errors.New("out of bounds")

// BoolSlice struct
type BoolSlice struct {
	raw  []byte
	size int
}

// NewBoolSlice builds new empty BoolSlice
func NewBoolSlice(size int) (*BoolSlice, error) {
	if size%8 != 0 {
		return nil, errors.New("size must be a multiple of 8")
	}
	return &BoolSlice{
		raw:  make([]byte, size/8),
		size: size,
	}, nil
}

// Rebuild generates new BoolSlice from data
func Rebuild(size int, data []byte) (*BoolSlice, error) {
	if size%8 != 0 {
		return nil, errors.New("size must be a multiple of 8")
	}
	if data == nil {
		return nil, errors.New("data cannot be empty")
	}
	return &BoolSlice{
		raw:  data,
		size: size,
	}, nil
}

// Set sets a bit from index passed
func (b *BoolSlice) Set(index int) error {
	if index > b.size {
		return ErrorOutOfBounds
	}

	internal := index / 8
	offset := index % 8
	b.raw[internal] |= (1 << offset)
	return nil
}

// Clear clears a bit from index passed
func (b *BoolSlice) Clear(index int) error {
	if index > b.size {
		return ErrorOutOfBounds
	}

	internal := index / 8
	offset := index % 8
	b.raw[internal] &= ((1 << offset) ^ 0xFF)
	return nil
}

// Get gets value from index passed
func (b *BoolSlice) Get(index int) (bool, error) {
	if index > b.size {
		return false, ErrorOutOfBounds
	}

	internal := index / 8
	offset := index % 8
	return (b.raw[internal] & (1 << offset)) != 0, nil
}

// Bytes gets bitmap
func (b *BoolSlice) Bytes() []byte {
	return b.raw
}
