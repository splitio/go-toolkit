package validator

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/splitio/go-toolkit/v5/datastructures/set"
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
			nestedNames := getFieldsForStructRecursive(prefixToUse+name, varType)
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
			nestedNames := getFieldsForMapRecursive(prefixToUse+name, value)
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
		if !p.Has(field) {
			return errors.New(field)
		}
	}
	return nil
}

// ValidateConfiguration compares s against p to validate each property
// and section.
func ValidateConfiguration(p interface{}, s []byte) error {
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

	m := map[string]interface{}{}
	err := json.Unmarshal(s, &m)
	if err != nil {
		return err
	}

	secondaryFieldList := getFieldsForMap(m)

	err = validateParameters(secondaryFieldList, primarySet)
	if err != nil {
		return errors.New("\"" + err.Error() + "\" is not a valid property in configuration")
	}
	return nil
}
