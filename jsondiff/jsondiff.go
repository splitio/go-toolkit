package jsondiff

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/splitio/go-toolkit/datastructures/set"
)

// Operation struct encapsulates an RFC6902 encoded jsonpatch operation
type Operation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func compareMap(src interface{}, dst interface{}, path string, ops []Operation) ([]Operation, error) {
	taSrc, okSrc := src.(map[string]interface{})
	taDst, okDst := dst.(map[string]interface{})
	if !okSrc || !okDst {
		return ops, errors.New("Incorrect cast")
	}

	srcKeys := set.NewSet()
	for k := range taSrc {
		srcKeys.Add(k)
	}

	dstKeys := set.NewSet()
	for k := range taDst {
		dstKeys.Add(k)
	}

	// Create "add" operations for new keys
	keysAdded := set.Difference(dstKeys, srcKeys)
	for _, added := range keysAdded.List() {
		addedStr, ok := added.(string)
		if !ok {
			return nil, errors.New("Error casting map key to string")
		}
		ops = append(ops, Operation{
			Op:    "add",
			Path:  path + "/" + addedStr,
			Value: taDst[addedStr],
		})
	}

	// Create "remove" operations for removed keys
	keysRemoved := set.Difference(srcKeys, dstKeys)
	for _, removed := range keysRemoved.List() {
		removedStr, ok := removed.(string)
		if !ok {
			return nil, errors.New("Error casting map key to string")
		}
		ops = append(ops, Operation{
			Op:   "remove",
			Path: path + "/" + removedStr,
		})
	}

	// Traverse keys appering in both objects
	keysPossiblyUpdated := set.Intersection(srcKeys, dstKeys)
	for _, key := range keysPossiblyUpdated.List() {
		keyStr, ok := key.(string)
		if !ok {
			return nil, errors.New("Error casting map key to string")
		}
		var err error
		ops, err = compare(taSrc[keyStr], taDst[keyStr], path+"/"+keyStr, ops)
		if err != nil {
			return nil, err
		}
	}

	return ops, nil
}

func compareArray(src interface{}, dst interface{}, path string, ops []Operation) ([]Operation, error) {
	taSrc, okSrc := src.([]interface{})
	taDst, okDst := dst.([]interface{})

	if !okSrc || !okDst {
		return ops, errors.New("Incorrect cast")
	}

	lengthSource := len(taSrc)
	lengthDestination := len(taDst)

	for index := 0; float64(index) < math.Max(float64(lengthSource), float64(lengthDestination)); index++ {
		if index < lengthSource && index < lengthDestination {
			// If the index exists in both arrays, compare the item
			var err error
			ops, err = compare(taSrc[index], taDst[index], fmt.Sprintf("%s/%d", path, index), ops)
			if err != nil {
				return nil, err
			}
		} else if index < lengthSource {
			// If the index only exists in the source array, remove the item
			ops = append(ops, Operation{
				Op:   "remove",
				Path: fmt.Sprintf("%s/%d", path, index),
			})
		} else if index < lengthDestination {
			// If the index only exists in the destination array, add the item
			ops = append(ops, Operation{
				Op:    "add",
				Path:  fmt.Sprintf("%s/%d", path, index),
				Value: taDst[index],
			})
		}
	}

	return ops, nil
}

func compare(src interface{}, dst interface{}, path string, ops []Operation) ([]Operation, error) {
	srcType := reflect.TypeOf(src).Kind()
	dstType := reflect.TypeOf(dst).Kind()
	var err error

	if srcType == reflect.Map && dstType == reflect.Map {
		ops, err = compareMap(src, dst, path, ops)
	} else if (srcType == reflect.Array || srcType == reflect.Slice) &&
		(dstType == reflect.Array || dstType == reflect.Slice) {
		ops, err = compareArray(src, dst, path, ops)
	} else if src != dst {
		ops = append(ops, Operation{
			Op:    "replace",
			Path:  path,
			Value: dst,
		})
	}

	return ops, err
}

// Calculate the difference between two json.RawMessage variables
func Calculate(src json.RawMessage, dst json.RawMessage) ([]Operation, error) {
	var sourceObj interface{}
	var destinationObj interface{}

	errSrc := json.Unmarshal(src, &sourceObj)
	if errSrc != nil {
		return nil, errSrc
	}

	errDst := json.Unmarshal(dst, &destinationObj)
	if errDst != nil {
		return nil, errSrc
	}

	ops := make([]Operation, 0)
	var err error
	ops, err = compare(sourceObj, destinationObj, "", ops)
	return ops, err

}
