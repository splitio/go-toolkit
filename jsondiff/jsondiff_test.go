package jsondiff

import (
	"testing"
)

func opsEquals(expected []Operation, actual []Operation) bool {
	if len(expected) != len(actual) {
		return false
	}

	for _, e := range expected {
		found := false
		for _, a := range actual {
			if a.Op == e.Op && a.Path == e.Path && a.Value == e.Value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestSimpleObject(t *testing.T) {
	json1 := `
    {
	"propA": 1,
	"propB": 2
    }`

	json2 := `
    {
	"propB": 3,
	"propC": 4
    }`
	res, err := Calculate([]byte(json1), []byte(json2))
	if err != nil {
		t.Error(err)
		return
	}

	expected := []Operation{
		{Op: "remove", Path: "/propA"},
		{Op: "replace", Path: "/propB", Value: float64(3)},
		{Op: "add", Path: "/propC", Value: float64(4)},
	}

	if !opsEquals(expected, res) {
		t.Error("Incorrect operations")
		t.Errorf("Got: %v", res)
		t.Errorf("Expected %v", expected)
	}
}

func TestSimpleArray(t *testing.T) {
	json1 := `[1,2,3,4]`
	json2 := `[1,8,3,4,7]`
	res, err := Calculate([]byte(json1), []byte(json2))
	if err != nil {
		t.Error(err)
		return
	}

	expected := []Operation{
		{Op: "replace", Path: "/1", Value: float64(8)},
		{Op: "add", Path: "/4", Value: float64(7)},
	}

	if !opsEquals(expected, res) {
		t.Error("Incorrect operations")
		t.Errorf("Got: %v", res)
		t.Errorf("Expected %v", expected)
	}
}

func TestNestedObject(t *testing.T) {
	json1 := `
    {
	"Name": "Someone",
	"Address": {
	    "Street": "Some Street",
	    "Number": 123
	},
	"Children": ["child1", "child2"]
    }`

	json2 := `
    {
	"Name": "Someone",
	"Address": {
	    "Street": "Some Other Street",
	    "Number": 456
	},
	"Children": ["child1", "child2", "child3"]
    }`
	res, err := Calculate([]byte(json1), []byte(json2))
	if err != nil {
		t.Error(err)
		return
	}

	expected := []Operation{
		{Op: "replace", Path: "/Address/Street", Value: "Some Other Street"},
		{Op: "replace", Path: "/Address/Number", Value: float64(456)},
		{Op: "add", Path: "/Children/2", Value: "child3"},
	}

	if !opsEquals(expected, res) {
		t.Error("Incorrect operations")
		t.Errorf("Got: %v", res)
		t.Errorf("Expected %v", expected)
	}
}
