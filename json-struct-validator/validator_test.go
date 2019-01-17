package validator

import (
	"testing"
)

func TestLen0(t *testing.T) {
	type OriginChild struct {
		two   string
		three int
	}

	originChild := OriginChild{two: "Test", three: 1}

	err := ValidateConfiguration(originChild, nil)
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

	err := ValidateConfiguration(origin, []byte("{\"one\": 10, \"originChild\": {\"two\": \"test\", \"three\": 10}}"))
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

	err := ValidateConfiguration(origin, []byte("{\"four\": 10, \"originChild\": {\"two\": \"test\", \"three\": 10}}"))
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"four\" is not a valid property in configuration" {
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

	err := ValidateConfiguration(origin, []byte("{\"one\": 10, \"originChild\": {\"two\": \"test\", \"four\": 10}}"))
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"originChild.four\" is not a valid property in configuration" {
		t.Error("Wrong message", err.Error())
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

	err := ValidateConfiguration(origin, []byte("{\"one\": 10, \"testChild\": {\"two\": \"test\", \"three\": 10}}"))
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"testChild\" is not a valid property in configuration" {
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

	err := ValidateConfiguration(origin, []byte("{\"one\": 10, \"originChild\": {\"two\": \"test\", \"three\": 10, \"four\": 10}}"))
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"originChild.four\" is not a valid property in configuration" {
		t.Error("Wrong message=", err.Error())
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

	err := ValidateConfiguration(origin, []byte("{\"one\": 10, \"originChild\": {\"two\": \"test\", \"three\": 10, \"four\": true}}"))
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"originChild.four\" is not a valid property in configuration" {
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

	err := ValidateConfiguration(origin, []byte("{\"one\": 10, \"originChild\": {\"two\": \"test\", \"three\": 10, \"four\": 10}}"))
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"originChild.four\" is not a valid property in configuration" {
		t.Error("Wrong message=")
	}
}

func TestSameThirdLevel(t *testing.T) {
	type Child struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type OriginChild struct {
		Child Child `json:"child"`
		Three int   `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	child := Child{Two: "Test", Three: 1}
	originChild := OriginChild{Child: child, Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	err := ValidateConfiguration(origin, []byte("{\"one\": 10, \"originChild\": {\"child\": {\"two\": \"test\", \"three\": 10}, \"three\": 10}}"))
	if err != nil {
		t.Error(err.Error())

		t.Error("Should not inform error")
	}
}

func TestDifferenthirdLevel(t *testing.T) {
	type Child struct {
		Two   string `json:"two"`
		Three int    `json:"three"`
	}

	type OriginChild struct {
		Child Child `json:"child"`
		Three int   `json:"three"`
	}

	type Origin struct {
		OriginChild OriginChild `json:"originChild"`
		One         int         `json:"one"`
	}

	child := Child{Two: "Test", Three: 1}
	originChild := OriginChild{Child: child, Three: 1}
	origin := Origin{OriginChild: originChild, One: 1}

	err := ValidateConfiguration(origin, []byte("{\"one\": 10, \"originChild\": {\"child\": {\"t\": \"test\", \"three\": 10}, \"three\": 10}}"))
	if err == nil {
		t.Error("Should inform error")
	}
	if err.Error() != "\"originChild.child.t\" is not a valid property in configuration" {
		t.Error("Wrong message", err.Error())
	}
}
