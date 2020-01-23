package common

import (
	"testing"
)

func TestStringRef(t *testing.T) {
	a := "someString"
	if a != *StringRef(a) {
		t.Error("Wrong string reference")
	}
}

func TestInt64Ref(t *testing.T) {
	a := int64(3)
	if a != *Int64Ref(a) {
		t.Error("Wrong int64 reference")
	}
}

func TestIntRef(t *testing.T) {
	a := int(43)
	if a != *IntRef(a) {
		t.Error("Wrong int reference")
	}
}

func TestInt64Value(t *testing.T) {
	a := int64(3)
	if Int64Value(&a) != 3 {
		t.Error("Should be 3")
	}

	if Int64Value(nil) != 0 {
		t.Error("Should be 0")
	}
}
