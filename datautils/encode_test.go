package datautils

import "testing"

func TestError(t *testing.T) {
	_, err := Encode([]byte("err"), 4)
	if err == nil || err.Error() != "encode type not found" {
		t.Error("It should return err")
	}

	_, err = Decode("err", 4)
	if err == nil || err.Error() != "encode type not found" {
		t.Error("It should return err")
	}
}

func TestB64EncodeDecode(t *testing.T) {
	data := "encode b64"
	encoded, err := Encode([]byte(data), Base64)
	if err != nil {
		t.Error("It should not return err")
	}

	decoded, err := Decode(encoded, Base64)
	if err != nil {
		t.Error("It should not return err")
	}

	if data != string(decoded) {
		t.Error("It should be equal")
	}
}
