package common

import (
	"testing"
)

func TestStringValueOrDefault(t *testing.T) {
	if StringValueOrDefault("abc", "def") != "abc" {
		t.Error("Should have returned original value")
	}

	if StringValueOrDefault("", "def") != "def" {
		t.Error("Should have returned default value")
	}
}

func TestStringFromRef(t *testing.T) {
	s1 := "string1"
	var s2 *string = nil
	if StringFromRef(&s1) != "string1" {
		t.Error("Should have returned original value")
	}

	if StringFromRef(s2) != "" {
		t.Error("Should have returned empty string")
	}
}
