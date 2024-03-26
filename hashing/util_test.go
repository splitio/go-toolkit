package hashing

import "testing"

func TestEncode(t *testing.T) {
	hash, err := Encode(nil, "something")
	if hash != "" {
		t.Error("Unexpected result")
	}
	if err == nil || err.Error() != "Hasher could not be nil" {
		t.Error("Unexpected error message")
	}

	hasher := NewMurmur332Hasher(0)
	hash2, err := Encode(hasher, "something")
	if err != nil {
		t.Error("It should not return error")
	}
	if hash2 != "NDE0MTg0MjI2MQ==" {
		t.Error("Unexpected result")
	}
}
