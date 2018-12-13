package helpers

import (
	"fmt"
	"testing"
)

func TestLen0(t *testing.T) {
	type OriginChild struct {
		two   string
		three int
	}

	originChild := OriginChild{two: "Test", three: 1}

	test := make(map[string]interface{})

	err := ValidateConfiguration(originChild, test)
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "no configuration provided" {
		t.Error("Wrong message")
	}
}

func TestSame(t *testing.T) {
	type OriginChild struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	originChild := OriginChild{Two: "Test", Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	testChild := map[string]interface{}{"two": "test", "three": 10}
	test := map[string]interface{}{"one": 10, "originChild": testChild}

	err := ValidateConfiguration(origin, test)
	if err != nil {
		t.Error("Should not inform error")
	}
}

func TestDifferentPropertyParent(t *testing.T) {
	type OriginChild struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	originChild := OriginChild{Two: "Test", Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	testChild := map[string]interface{}{"two": "test", "three": 10}
	test := map[string]interface{}{"four": 10, "originChild": testChild}

	err := ValidateConfiguration(origin, test)
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"four\" is not a valid section or property in configuration" {
		t.Error("Wrong message")
	}
}

func TestDifferentPropertyChild(t *testing.T) {
	type OriginChild struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	originChild := OriginChild{Two: "Test", Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	testChild := map[string]interface{}{"two": "test", "four": 10}
	test := map[string]interface{}{"one": 10, "originChild": testChild}

	err := ValidateConfiguration(origin, test)
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"four\" in section \"originChild\" is not valid configuration" {
		t.Error("Wrong message")
	}
}

func TestDifferentParentAndChild(t *testing.T) {
	type OriginChild struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	originChild := OriginChild{Two: "Test", Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	testChild := map[string]interface{}{"two": "test", "three": 10}
	test := map[string]interface{}{"one": 10, "testChild": testChild}

	err := ValidateConfiguration(origin, test)
	fmt.Println(err.Error())
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"testChild\" is not a valid section or property in configuration" {
		t.Error("Wrong message, it should inform parent")
	}
}

func TestDifferentPropertyInChild(t *testing.T) {
	type OriginChild struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	originChild := OriginChild{Two: "Test", Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	testChild := map[string]interface{}{"two": "test", "three": 10, "four": 10}
	test := map[string]interface{}{"one": 10, "originChild": testChild}

	err := ValidateConfiguration(origin, test)
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"four\" in section \"originChild\" is not valid configuration" {
		t.Error("Wrong message=")
	}
}

func TestDifferentPropertyInChildBool(t *testing.T) {
	type OriginChild struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	originChild := OriginChild{Two: "Test", Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	testChild := map[string]interface{}{"two": "test", "three": 10, "four": true}
	test := map[string]interface{}{"one": 10, "originChild": testChild}

	err := ValidateConfiguration(origin, test)
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"four\" in section \"originChild\" is not valid configuration" {
		t.Error("Wrong message=")
	}
}

func TestDifferentPropertyInChildNumber(t *testing.T) {
	type OriginChild struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	originChild := OriginChild{Two: "Test", Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	testChild := map[string]interface{}{"two": "test", "three": 10, "four": 10}
	test := map[string]interface{}{"one": 10, "originChild": testChild}

	err := ValidateConfiguration(origin, test)
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"four\" in section \"originChild\" is not valid configuration" {
		t.Error("Wrong message=")
	}
}
