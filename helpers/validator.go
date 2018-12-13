package helpers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/splitio/go-toolkit/datastructures/set"
)

func getFieldsForStructRecursive(prefix string, structType reflect.Type) []string {
	var prefixToUse string
	if prefix != "" {
		prefixToUse = prefix + "."
	} else {
		prefixToUse = ""
	}

	currentNames := make([]string, 0)
	for i := 0; i < structType.NumField(); i++ {
		name := structType.Field(i).Tag.Get("json")
		if strings.TrimSpace(name) == "" {
			name = structType.Field(i).Name
		}
		varType := structType.Field(i).Type
		currentNames = append(currentNames, prefixToUse+name)
		if varType.Kind() == reflect.Struct {
			nestedNames := getFieldsForStructRecursive(name, varType)
			currentNames = append(currentNames, nestedNames...)
		}
	}
	return currentNames
}

func getFieldsForStruct(s interface{}) []string {
	t := reflect.TypeOf(s)
	return getFieldsForStructRecursive("", t)
}

func getFieldsForMapRecursive(prefix string, s map[string]interface{}) []string {
	var prefixToUse string
	if prefix != "" {
		prefixToUse = prefix + "."
	} else {
		prefixToUse = ""
	}

	currentNames := make([]string, 0)
	for name, value := range s {
		switch value := value.(type) {
		case string, int, int8, int16, int32, int64, float32, float64, bool:
			currentNames = append(currentNames, prefixToUse+name)
		case map[string]interface{}:
			currentNames = append(currentNames, prefixToUse+name)
			nestedNames := getFieldsForMapRecursive(name, value)
			currentNames = append(currentNames, nestedNames...)
		}
	}
	return currentNames
}

func getFieldsForMap(s map[string]interface{}) []string {
	return getFieldsForMapRecursive("", s)
}

func validateParameters(userConf []string, p *set.ThreadUnsafeSet) error {
	for _, field := range userConf {
		if p.Has(field) == false {
			return errors.New(field)
		}
	}
	return nil
}

// ValidateConfiguration compares s against p to validate each property
// and section.
func ValidateConfiguration(p interface{}, s map[string]interface{}) error {
	if p == nil {
		return errors.New("configuration cannot be null")
	}

	if len(s) == 0 {
		return errors.New("no configuration provided")
	}

	primaryFieldList := getFieldsForStruct(p)
	primarySet := set.NewSet()
	for _, c := range primaryFieldList {
		primarySet.Add(c)
	}

	secondaryFieldList := getFieldsForMap(s)
	fmt.Println(secondaryFieldList)

	err := validateParameters(secondaryFieldList, primarySet)
	if err != nil {
		var m string
		message := err.Error()
		messageSplit := strings.Split(message, ".")
		if len(messageSplit) == 1 {
			m = "\"" + messageSplit[0] + "\" is not a valid section or property in configuration"
		} else {
			m = "\"" + messageSplit[1] + "\" in section \"" + messageSplit[0] + "\" is not valid configuration"
		}
		return errors.New(m)
	}
	return nil
}
