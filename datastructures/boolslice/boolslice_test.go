package boolslice

import (
	"math"
	"testing"
)

func TestBoolSlice(t *testing.T) {
	_, err := NewBoolSlice(12)
	if err == nil {
		t.Error("It should return err")
	}

	b, err := NewBoolSlice(int(math.Pow(2, 15)))
	if err != nil {
		t.Error("It should not return err", err)
	}

	i1 := 12
	i2 := 20
	i3 := 123
	i4 := 2000
	i5 := 8192

	if err := b.Set(int(math.Pow(2, 15)) + 1); err == nil {
		t.Error("It should return err")
	}
	if err := b.Set(i1); err != nil {
		t.Error("It should not return err")
	}
	if err := b.Set(i2); err != nil {
		t.Error("It should not return err")
	}
	if err := b.Set(i3); err != nil {
		t.Error("It should not return err")
	}
	if err := b.Set(i4); err != nil {
		t.Error("It should not return err")
	}
	if err := b.Set(i5); err != nil {
		t.Error("It should not return err")
	}

	if _, err := b.Get(int(math.Pow(2, 15)) + 1); err == nil {
		t.Error("It should return err")
	}
	if v, _ := b.Get(i1); !v {
		t.Error("It should match", i1)
	}
	if v, _ := b.Get(i2); !v {
		t.Error("It should match", i2)
	}
	if v, _ := b.Get(i3); !v {
		t.Error("It should match", i3)
	}
	if v, _ := b.Get(i4); !v {
		t.Error("It should match", i4)
	}
	if v, _ := b.Get(i5); !v {
		t.Error("It should match", i5)
	}
	if v, _ := b.Get(200); v {
		t.Error("It should not match 200")
	}
	if v, _ := b.Get(5000); v {
		t.Error("It should not match 5000")
	}

	if len(b.Bytes()) != int(math.Pow(2, 15)/8) {
		t.Error("Len should be 4096")
	}

	if err := b.Clear(int(math.Pow(2, 15)) + 1); err == nil {
		t.Error("It should return err")
	}
	if err := b.Clear(i1); err != nil {
		t.Error("It should not return err")
	}

	if v, _ := b.Get(i1); v {
		t.Error("It should not match after cleared", i1)
	}

	if _, err := Rebuild(1, nil); err.Error() != "size must be a multiple of 8" {
		t.Error("It should return err")
	}

	if _, err := Rebuild(8, nil); err.Error() != "data cannot be empty" {
		t.Error("It should return err")
	}

	rebuilt, err := Rebuild(int(math.Pow(2, 15)), b.Bytes())
	if err != nil {
		t.Error("It should not return err")
	}
	if v, _ := rebuilt.Get(i2); !v {
		t.Error("It should match", i2)
	}
	if v, _ := rebuilt.Get(i3); !v {
		t.Error("It should match", i3)
	}
	if v, _ := rebuilt.Get(i4); !v {
		t.Error("It should match", i4)
	}
	if v, _ := rebuilt.Get(i5); !v {
		t.Error("It should match", i5)
	}
	if v, _ := rebuilt.Get(i1); v {
		t.Error("It should not match 12")
	}
	if v, _ := rebuilt.Get(200); v {
		t.Error("It should not match 200")
	}
	if v, _ := rebuilt.Get(5000); v {
		t.Error("It should not match 5000")
	}
}
