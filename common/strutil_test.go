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
