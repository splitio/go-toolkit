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

func TestStringRefOrNil(t *testing.T) {
	a := "someString"
	if a != *StringRefOrNil(a) {
		t.Error("Wrong string reference")
	}

	if StringRefOrNil("") != nil {
		t.Error("Should be nil")
	}
}

func TestInt64RefOrNil(t *testing.T) {
	a := int64(3)
	if a != *Int64RefOrNil(a) {
		t.Error("Wrong int64 reference")
	}
	if Int64RefOrNil(0) != nil {
		t.Error("Should be nil")
	}
}

func TestIntRefOrNil(t *testing.T) {
	a := int(43)
	if a != *IntRefOrNil(a) {
		t.Error("Wrong int reference")
	}
	if IntRefOrNil(0) != nil {
		t.Error("Should be nil")
	}
}

func TestAsInt64OrNil(t *testing.T) {
	a := AsInt64OrNil(int64(123456789))
	if a == nil {
		t.Error("Wrong int reference")
	}
	b := AsInt64OrNil("some")
	if b != nil {
		t.Error("It should be nil")
	}
}

func TestAsStringOrNil(t *testing.T) {
	a := AsStringOrNil("some")
	if a == nil {
		t.Error("Wrong string reference")
	}
	b := AsStringOrNil(int64(123456789))
	if b != nil {
		t.Error("It should be nil")
	}
}
