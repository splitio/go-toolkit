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
