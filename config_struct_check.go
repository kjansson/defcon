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
// The supported tags are "default" and "required". Supported types to tag are all ints, floats and string.
// It will modify all struct fields where the tag `default:"<value>"` is present and a valid value is given.
//
// It will return an error if one or more of the following conditions are met.
// A field is tagged with `required:"true"` and A: is a string but the value is empty or B: is a numerical type but the value is zero.
// A field is tagged with `default:"<value>"` but the value is not valid for the type of the field
// A field is tagged with both `default:"<value>"` and `required:"true"`
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

func checkStruct(v *reflect.Value) error {

	var bits int
	var err error
	var requiredTag, defaultTag, requiresTag bool
	var defaultTagValue, requiredTagValue, requiresValue string
	var integer int64
	var floating float64

	requiresMap := make(map[string][]string)
	setFields := []string{}

	typeRegex, _ := regexp.Compile("^([a-zA-Z]+)([0-9]*)") // Extract type family (int, float, etc) and the number of bits for the type

	for i := 0; i < v.NumField(); i++ { // Loop through fields in struct
		if v.Field(i).Kind() == reflect.Struct { // This is a nested struct
			c := reflect.ValueOf(v.Field(i).Interface())
			if err := checkStruct(&c); err != nil {
				return err
			}
		} else {
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
						return fmt.Errorf("Error while evaluating requires value with regex: %s", err)
					}
					if !match {
						return fmt.Errorf("Field name %s required by %s does not seem to have a valid name", fieldName, v.Type().Field(i).Name)
					} else {
						requiresMap[v.Type().Field(i).Name] = append(requiresMap[v.Type().Field(i).Name], fieldName)
					}
				}
			}

			if defaultTag && requiredTag { // Both default and required is not allowed
				return fmt.Errorf("Having both default and required tags present in field %s is not allowed.", v.Type().Field(i).Name)
			}

			if v.Type().Field(i).IsExported() {
				// Get the type family and number of bits
				typeInfo := typeRegex.FindStringSubmatch(v.Field(i).Kind().String())
				if typeInfo[2] == "" {
					bits = 0
				} else {
					bits, _ = strconv.Atoi(typeInfo[2])
				}
				switch typeInfo[1] {
				case "int": // Integer type
					if integer, err = strconv.ParseInt(defaultTagValue, 10, bits); err != nil { // Parse string to int
						return err
					}
					if v.Field(i).Int() == 0 { // If zero
						if requiredTag { // And required, not allowed
							return fmt.Errorf("Integer field %s is marked as required but has zero value.", v.Type().Field(i).Name)
						}
						if defaultTag { // If default value exists, set it
							v.Field(i).SetInt(int64(integer))
							setFields = append(setFields, v.Type().Field(i).Name)
						}
						if _, ok := requiresMap[v.Type().Field(i).Name]; ok {
							// If field requires other fields but is not itself set, we should ignore the requirements
							delete(requiresMap, v.Type().Field(i).Name)
						}
					} else {
						// This field has a value, save as set
						setFields = append(setFields, v.Type().Field(i).Name)
					}
				case "float": // Float type
					if floating, err = strconv.ParseFloat(defaultTagValue, bits); err != nil { // Parse string to float
						return err
					}
					if v.Field(i).Float() == 0 { // If zero
						if requiredTag { // And required, not allowd
							return fmt.Errorf("Float field %s is marked as required but has zero value.", v.Type().Field(i).Name)
						}
						if defaultTag { // If default value exists, set it
							v.Field(i).SetFloat(floating)
							setFields = append(setFields, v.Type().Field(i).Name)
						}
						if _, ok := requiresMap[v.Type().Field(i).Name]; ok {
							// If field requires other fields but is not itself set, we should ignore the requirements
							delete(requiresMap, v.Type().Field(i).Name)
						}
					} else {
						// This field has a value, save as set
						setFields = append(setFields, v.Type().Field(i).Name)
					}
				case "string": // String type
					if v.Field(i).Len() == 0 { // If zero length
						if requiredTag { // And requred, not allowed
							return fmt.Errorf("required value missing in string field %s", v.Type().Field(i).Name)
						}
						if defaultTag { // If default value exists, set it
							v.Field(i).SetString(defaultTagValue)
							setFields = append(setFields, v.Type().Field(i).Name)
						}
						if _, ok := requiresMap[v.Type().Field(i).Name]; ok {
							// If field requires other fields but is not itself set, we should ignore the requirements
							delete(requiresMap, v.Type().Field(i).Name)
						}
					} else {
						// This field has a value, save as set
						setFields = append(setFields, v.Type().Field(i).Name)
					}
				}
			} else {
				fmt.Printf("Warning: default or required tag detected on unexported field %s", v.Type().Field(i).Name)
			}
		}
	}

	for parentField, requiredFields := range requiresMap {
		for _, requiredField := range requiredFields {
			if !existsIn(setFields, requiredField) {
				return fmt.Errorf("Field %s is required by field %s but is not set", requiredField, parentField)
			}
		}
	}

	return nil
}
