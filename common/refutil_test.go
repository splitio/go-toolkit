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
