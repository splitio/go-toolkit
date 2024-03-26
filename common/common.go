package common

import "cmp"

// New helpers to be used when with newer go versions.
// the rest of the common package should be removed in v3, and consumers of the lib
// should only rely on these functions

// Ref creates a copy of `x` in heap and returns a pointer to it
func Ref[T any](x T) *T {
	return &x
}

// RefOrNil returns a pointer to the value supplied if it's not the default value, nil otherwise
func RefOrNil[T comparable](x T) *T {
    var t T
    if x == t {
        return nil
    }
    return &x
}

// PointerOf performs a type-assertion to T and returns a pointer if successful, nil otherwise.
func PointerOf[T any](x interface{}) *T {
	if x == nil {
		return nil
	}

	ta, ok := x.(T)
	if !ok {
		return nil
	}

	return &ta
}

// PartitionSliceByLength partitions a slice into multiple slices of up to `maxItems` size
func PartitionSliceByLength[T comparable](items []T, maxItems int) [][]T {
	var splitted [][]T
	for i := 0; i < len(items); i += maxItems {
		end := i + maxItems
		if end > len(items) {
			end = len(items)
		}
		splitted = append(splitted, items[i:end])
	}
	return splitted
}

// DedupeInNewSlice creates a new slice from `items` without duplicate elements
func UnorderedDedupedCopy[T comparable](items []T) []T {
	present := make(map[T]struct{}, len(items))
	for idx := range items {
		present[items[idx]] = struct{}{}
	}

	ret := make([]T, 0, len(present))
	for key := range present {
		ret = append(ret, key)
	}

	return ret
}

// ValueOr returns the supplied value if it has something other than the default value
// for type T. Returns `fallback` otherwise
func ValueOr[T comparable](in T, fallback T) T {
	var t T
	if in == t {
		return fallback
	}
	return in
}

// Max returns the greatest item of all supplied
func Max[T cmp.Ordered](i1 T, rest ...T) T {
	max := i1
	for idx := range rest {
		if rest[idx] > max {
			max = rest[idx]
		}
	}
	return max
}

// Min returns the minimum item of all supplied
func Min[T cmp.Ordered](i1 T, rest ...T) T {
	min := i1
	for idx := range rest {
		if rest[idx] < min {
			min = rest[idx]
		}
	}
	return min
}
