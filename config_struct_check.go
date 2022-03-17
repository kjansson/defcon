package defcon

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// CheckConfigStruct accepts any struct (supports nested structs) and will check all values and their tags.
// It will change all values where the tag `default:"<value>"` is present and a valid value is given.
//
// It will return an error if one or more of the following conditions are met.
// A field is tagged with `required:"true"` and A: is a string but the value is empty or B: is a numerical type but the value is zero.
// A field is tagged with `default:"<value>"` but the value is not valid for the type of the field
// A field is tagged with both `default:"<value>"` and `required:"true"`
func CheckConfigStruct(config interface{}) error {
	c := reflect.ValueOf(config).Elem()
	return checkStruct(&c)
}

func checkStruct(v *reflect.Value) error {

	var bits int
	var err error
	var requiredTag, defaultTag bool
	var defaultTagValue string

	typeRegex, _ := regexp.Compile("^([a-zA-Z]+)([0-9]*)")

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			c := reflect.ValueOf(v.Field(i).Interface())
			if err := checkStruct(&c); err != nil {
				return err
			}
		} else {
			_, requiredTag = v.Type().Field(i).Tag.Lookup("required")
			defaultTagValue, defaultTag = v.Type().Field(i).Tag.Lookup("default")

			if defaultTag || requiredTag {

				typeInfo := typeRegex.FindStringSubmatch(v.Field(i).Kind().String())

				switch typeInfo[1] {
				case "int":
					if typeInfo[2] == "" {
						bits = 0
					} else {
						bits, err = strconv.Atoi(typeInfo[2])
						if err != nil {
							return fmt.Errorf("Error while parsing integer value on field %s. Invalid type?", v.Type().Field(i).Name)
						}
					}
					integer, err := strconv.ParseInt(defaultTagValue, 10, bits)
					if err != nil {
						return err
					}
					if requiredTag && integer == 0 {
						return fmt.Errorf("Integer field %s is marked as required but has zero value", v.Type().Field(i).Name)
					}
					if defaultTag {
						v.Field(i).SetInt(int64(integer))
					}
				case "float":
					if typeInfo[2] == "" {
						bits = 0
					} else {
						bits, err = strconv.Atoi(typeInfo[2])
						if err != nil {
							return fmt.Errorf("Error while parsing integer value on field %s. Invalid type?", v.Type().Field(i).Name)
						}
					}
					floating, err := strconv.ParseFloat(defaultTagValue, bits)
					if err != nil {
						return err
					}
					if requiredTag && floating == 0 {
						return fmt.Errorf("Error while parsing float value on field %s. Invalid type?", v.Type().Field(i).Name)
					}
					if defaultTag {
						v.Field(i).SetFloat(floating)
					}
				case "string":
					if requiredTag && v.Field(i).Len() == 0 {
						return fmt.Errorf("required value missing in string field %s", v.Type().Field(i).Name)
					}
					if defaultTag {
						v.Field(i).SetString(defaultTagValue)
					}
				}
			} else if _, ok := v.Type().Field(i).Tag.Lookup("default"); ok {
				return fmt.Errorf("Having both default and required tags present in field %s is not allowed", v.Type().Field(i).Name)
			}
		}
	}
	return nil
}
