package jsonvalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen0(t *testing.T) {
	type OriginChild struct {
		two   string
		three int
	}

	originChild := OriginChild{two: "Test", three: 1}

	err := ValidateConfiguration(originChild, nil)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "no configuration provided")
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
	assert.Nil(t, err)
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
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "\"four\" is not a valid property in configuration")
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
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "\"originChild.four\" is not a valid property in configuration")
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
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "\"testChild\" is not a valid property in configuration")
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
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "\"originChild.four\" is not a valid property in configuration")
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
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "\"originChild.four\" is not a valid property in configuration")
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
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "\"originChild.four\" is not a valid property in configuration")
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
	assert.Nil(t, err)
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
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "\"originChild.child.t\" is not a valid property in configuration")
}
