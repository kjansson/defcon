// defcon is a minimalistic library for parsing tagged config structs, automatically handling default and required values
package defcon

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// CheckConfigStruct accepts any struct (supports nested structs) and will check all exported values and their tags.
// The supported tags are "default", "required" and "requires". Supported types to tag are all ints, floats and string.
// It will modify all struct fields where the tag `default:"<value>"` is present and a valid value is given.
//
// It will return an error if one or more of the following conditions are met.
// A field is tagged with `required:"true"` and A: is a string but the value is empty or B: is a numerical type but the value is zero.
// A field is tagged with `default:"<value>"` but the value is not valid for the type of the field
// A field is tagged with both `default:"<value>"` and `required:"true"`
// A field is tagged with `requires:"<field1, field2, ...>"` is in use, and the value of any the given fields is A: is a string but the value is empty or B: is a numerical type but the value is zero.
func CheckConfigStruct(config interface{}) error {
	c := reflect.ValueOf(config).Elem()
	return checkStruct(&c)
}

func existsIn(subject []string, searchValue string) bool {
	for _, value := range subject {
		if value == searchValue {
			return true
		}
	}
	return false
}

func isStructEmpty(s reflect.Value) bool {
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).Kind() == reflect.Struct { // This is a nested struct
			//c := s.Field(i)

			return isStructEmpty(s.Field(i))
		} else {
			if !s.Field(i).IsZero() {
				return false
			}
		}
	}
	return true
}

func getTypeFamily(v reflect.Value) (string, int) {

	var bits int
	typeRegex, _ := regexp.Compile("^([a-zA-Z]+)([0-9]*)") // Extract type family (int, float, etc) and the number of bits for the type

	family := typeRegex.FindStringSubmatch(v.Kind().String())
	if family[2] == "" {
		bits = 0
	} else {
		bits, _ = strconv.Atoi(family[2])
	}

	return family[1], bits
}

// func isEmpty(v reflect.Value) bool {

// 	family, _ := getTypeFamily(v)

// 	switch family {
// 	case "int":
// 		if v.Int() == 0 { // If zero
// 			return true
// 		}
// 	case "float":
// 		if v.Float() == 0 { // If zero
// 			return true
// 		}
// 	case "string":
// 		if v.Len() == 0 { // If zero
// 			return true
// 		}
// 	}
// 	return false
// }

func setValue(v *reflect.Value, val string) error {

	family, bits := getTypeFamily(*v)

	switch family {
	case "int":
		integer, err := strconv.ParseInt(val, 10, bits)
		if err != nil { // Parse string to int
			return err
		}
		v.SetInt(integer)
	case "float":
		floating, err := strconv.ParseFloat(val, bits)
		if err != nil { // Parse string to int
			return err
		}
		v.SetFloat(floating)
	case "string":
		v.SetString(val)
	}

	return nil
}

func checkStruct(v *reflect.Value) error {

	var requiredTag, defaultTag, requiresTag bool
	var defaultTagValue, requiredTagValue, requiresValue string

	requiresMap := make(map[string][]string)
	setFields := []string{}

	for i := 0; i < v.NumField(); i++ { // Loop through fields in struct

		// Get tags
		requiredTag = false
		if requiredTagValue, _ = v.Type().Field(i).Tag.Lookup("required"); requiredTagValue == "true" || requiredTagValue == "TRUE" {
			requiredTag = true
		}

		defaultTagValue, defaultTag = v.Type().Field(i).Tag.Lookup("default")

		requiresValue, requiresTag = v.Type().Field(i).Tag.Lookup("requires")
		if requiresTag {
			requires := strings.Split(requiresValue, ",")
			for _, r := range requires {
				fieldName := strings.TrimSpace(r)
				match, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_-]+$", fieldName)
				if err != nil {
					return fmt.Errorf("error while evaluating requires value with regex: %s", err)
				}
				if !match {
					return fmt.Errorf("field name %s required by %s does not seem to have a valid name", fieldName, v.Type().Field(i).Name)
				} else {
					requiresMap[v.Type().Field(i).Name] = append(requiresMap[v.Type().Field(i).Name], fieldName)
				}
			}
		}

		if v.Field(i).Kind() == reflect.Struct { // This is a nested struct
			c := v.Field(i)
			if !c.IsZero() { // Store as set if struct nested struct has set fields
				setFields = append(setFields, v.Type().Field(i).Name)
			}

			if c.CanSet() {
				if err := checkStruct(&c); err != nil {
					return err
				}
			}
		} else {

			if defaultTag && requiredTag { // Both default and required is not allowed
				return fmt.Errorf("having both default and required tags present in field %s is not allowed", v.Type().Field(i).Name)
			}

			if v.Type().Field(i).IsExported() {
				if v.Field(i).IsZero() { // If zero
					if requiredTag { // And required, not allowed
						return fmt.Errorf("field %s is marked as required but has zero/empty value", v.Type().Field(i).Name)
					}
					if defaultTag { // If default value exists, set it
						ptr := v.Field(i)
						err := setValue(&ptr, defaultTagValue)
						if err != nil {
							return fmt.Errorf("could not set value in field, %s", err)
						}
						setFields = append(setFields, v.Type().Field(i).Name)
					} else {
						// If field requires other fields but is not itself set, we should ignore the requirements
						delete(requiresMap, v.Type().Field(i).Name)
					}
				} else {
					// This field has a value, save as set
					setFields = append(setFields, v.Type().Field(i).Name)
				}
			}
		}
	}

	for parentField, requiredFields := range requiresMap {
		for _, requiredField := range requiredFields {
			if !existsIn(setFields, requiredField) {
				return fmt.Errorf("field %s requires field %s which is empty/not set", parentField, requiredField)
			}
		}
	}

	return nil
}
