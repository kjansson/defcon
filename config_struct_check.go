// defcon is a minimalistic library for parsing tagged config structs, automatically handling required/default values and field dependencies
package defcon

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

// CheckConfigStruct accepts any struct (supports nested structs) and will inspect all fields and their tags.
// The package supports the tags "default", "required" and "requires". Supported types to tag are all ints, floats and string and slices of these types (structs support the "required" tag).
// Behaviour;
// The "default" tag will modify the struct field with the tag value, if the original value is the primitive type default, i.e. zero for numerical values, or zero length string.
// The "required" tag will return an error if the fields value is the primitive type default. If applied to a struct, the struct will be considered empty if all of its fields have primitive type default values.
// The "default" tag will be applied first, so if a field is tagged with both "default" and "required", the "required" tag will have no effect.
// The "requires" tag will return an error if any of the given fields values (within the same struct) have their primitive type default or is an empty struct.
// Tags with invalid values such as references to non-existing fields, values that will overflow the numerical types, invalid numerical values, etc. will result in an error.

func CheckConfigStruct(config interface{}) error {

	c := reflect.ValueOf(config).Elem()
	return checkStruct(&c)
}

// Finds a string value in an array of strings
func existsIn(subject []string, searchValue string) bool {

	for _, value := range subject {
		if value == searchValue {
			return true
		}
	}
	return false
}

// Get reflection type and returns its type family and number of bits
func getTypeDetails(v reflect.Value) (string, int) {

	var bits int
	// Extract type family and bits (if any) from reflect type, e.g. "int32" => ["int", "32"]
	family := regexp.MustCompile("^([a-zA-Z]+)([0-9]*)").FindStringSubmatch(v.Kind().String())

	if family[2] == "" {
		bits = 0
	} else {
		bits, _ = strconv.Atoi(family[2]) // This should be safe w/o error checking since the vaule come from the reflect kind
	}
	return family[1], bits
}

// Get reflection type and returns its type family and number of bits
func getElementTypeDetails(v reflect.Value) (string, int) {

	var bits int
	// Extract type family and bits (if any) from reflect type, e.g. "int32" => ["int", "32"]
	e := v.Type().Elem()
	family := regexp.MustCompile("^([a-zA-Z]+)([0-9]*)").FindStringSubmatch(e.Kind().String())

	if family[2] == "" {
		bits = 0
	} else {
		bits, _ = strconv.Atoi(family[2]) // This should be safe w/o error checking since the vaule come from the reflect kind
	}
	return family[1], bits
}

// Sets a value
func setValue(v *reflect.Value, val string) error {

	family, bits := getTypeDetails(*v) // Get type family and number of bits if applicable

	// Parse numerical values if needed and set values
	switch family {
	case "int":
		integer, err := strconv.ParseInt(val, 10, bits) // Parse string to int
		if err != nil {
			return err
		}
		v.SetInt(integer)
	case "float":
		floating, err := strconv.ParseFloat(val, bits) // Parse string to float
		if err != nil {
			return err
		}
		v.SetFloat(floating)
	case "string":
		v.SetString(val)
	case "slice":
		eType, bits := getElementTypeDetails(*v)
		var uv = reflect.Value{}
		if eType == "int" {
			values := strings.Split(regexp.MustCompile("^{(.*)}").FindStringSubmatch(val)[1], ",")
			switch bits {
			case 8:
				us := make([]int8, 0)   // Create a slice of strings
				u := &us                // Create a pointer to the slice
				uv = reflect.ValueOf(u) // Get the reflect value of the pointer
			case 16:
				us := make([]int16, 0)
				u := &us
				uv = reflect.ValueOf(u)
			case 32:
				us := make([]int32, 0)
				u := &us
				uv = reflect.ValueOf(u)
			case 64:
				us := make([]int64, 0)
				u := &us
				uv = reflect.ValueOf(u)
			default:
				us := make([]int, 0)
				u := &us
				uv = reflect.ValueOf(u)
			}
			for _, x := range values { // Loop through all values
				numVal, err := strconv.ParseInt(strings.TrimSpace(x), 10, bits) // Parse string to int
				if err != nil {
					return err
				}
				switch bits {
				case 8:
					y := reflect.ValueOf(int8(numVal)) // Get the reflect value of the value
					uv.Elem().Set(reflect.Append(uv.Elem(), y))
				case 16:
					y := reflect.ValueOf(int16(numVal))
					uv.Elem().Set(reflect.Append(uv.Elem(), y))
				case 32:
					y := reflect.ValueOf(int32(numVal))
					uv.Elem().Set(reflect.Append(uv.Elem(), y))
				case 64:
					y := reflect.ValueOf(int64(numVal))
					uv.Elem().Set(reflect.Append(uv.Elem(), y))
				default:
					y := reflect.ValueOf(int(numVal))
					uv.Elem().Set(reflect.Append(uv.Elem(), y))
				}
			}
			v.Set(uv.Elem()) // Set the value of the reflect value to the slice
		} else if eType == "float" {
			values := strings.Split(regexp.MustCompile("^{(.*)}").FindStringSubmatch(val)[1], ",")
			switch bits {
			case 32:
				us := make([]float32, 0) // Create a slice of strings
				u := &us                 // Create a pointer to the slice
				uv = reflect.ValueOf(u)  // Get the reflect value of the pointer
			case 64:
				us := make([]float64, 0)
				u := &us
				uv = reflect.ValueOf(u)
			default:
				us := make([]float64, 0)
				u := &us
				uv = reflect.ValueOf(u)
			}
			for _, x := range values { // Loop through all values
				numVal, err := strconv.ParseFloat(strings.TrimSpace(x), bits) // Parse string to int
				if err != nil {
					return err
				}
				switch bits {
				case 32:
					y := reflect.ValueOf(float32(numVal)) // Get the reflect value of the value
					uv.Elem().Set(reflect.Append(uv.Elem(), y))
				case 64:
					y := reflect.ValueOf(float64(numVal))
					uv.Elem().Set(reflect.Append(uv.Elem(), y))
				default:
					y := reflect.ValueOf(float64(numVal))
					uv.Elem().Set(reflect.Append(uv.Elem(), y))
				}
			}
			v.Set(uv.Elem()) // Set the value of the reflect value to the slice
		} else if eType == "string" {
			values := strings.Split(regexp.MustCompile("^{(.*)}").FindStringSubmatch(val)[1], ",")
			us := make([]string, 0)    // Create a slice of strings
			u := &us                   // Create a pointer to the slice
			uv := reflect.ValueOf(u)   // Get the reflect value of the pointer
			for _, x := range values { // Loop through all values
				y := reflect.ValueOf(strings.TrimSpace(x))  // Get the reflect value of the value
				uv.Elem().Set(reflect.Append(uv.Elem(), y)) // Append the value to the slice
			}
			v.Set(uv.Elem()) // Set the value of the reflect value to the slice
		} else {
			return fmt.Errorf("slice type %s is not supported", eType)
		}
	}
	return nil
}

func checkStruct(v *reflect.Value) error {

	requiresMap := make(map[string][]string) // Map containing all fields tagged as "requires" by other
	setFields := []string{}                  // Map containing all fields with a non-empty value
	var field reflect.Value

	for i := 0; i < v.NumField(); i++ { // Loop through fields in struct
		// Get tags
		requiredValue, isRequired := v.Type().Field(i).Tag.Lookup("required")
		if requiredValue != "true" && requiredValue != "TRUE" {
			isRequired = false
		}
		defaultValue, hasDefault := v.Type().Field(i).Tag.Lookup("default")
		requiresValue, requiresField := v.Type().Field(i).Tag.Lookup("requires")

		if !v.Type().Field(i).IsExported() {
			field = reflect.NewAt(v.Field(i).Type(), unsafe.Pointer(v.Field(i).UnsafeAddr())).Elem() // Get access to unexported field
		} else {
			field = v.Field(i)
		}
		fieldName := v.Type().Field(i).Name

		if requiresField {
			for _, reqVal := range strings.Split(requiresValue, ",") { // Split all field names in requires tag value
				reqsFieldName := strings.TrimSpace(reqVal)
				match := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]+$").MatchString(reqsFieldName) // Check if name of field is valid
				if !match {
					return fmt.Errorf("field %s tagged as required by field %s does not seem to have a valid name", reqsFieldName, fieldName)
				} else {
					requiresMap[fieldName] = append(requiresMap[fieldName], reqsFieldName) // Add to requires map
				}
			}
		}

		if field.Kind() == reflect.Struct { // This is a nested struct
			if field.CanSet() {
				if err := checkStruct(&field); err != nil { // Drill down in nested struct
					return err
				}
			}
			if field.IsZero() && isRequired { // Empty but required should return error
				return fmt.Errorf("field %s (struct) is marked as required but has no fields with non-empty values", fieldName)
			}
			if !field.IsZero() { // Check if field is set and add it to set fields map
				setFields = append(setFields, fieldName)
			}
		} else {
			if field.IsZero() { // If zero
				if hasDefault { // If default value exists, set it
					err := setValue(&field, defaultValue)
					if err != nil {
						return fmt.Errorf("could not set value in field, %s", err)
					}
					setFields = append(setFields, fieldName)
				}
				if isRequired { // And required, not allowed
					return fmt.Errorf("field %s (%s) is marked as required but has zero/empty value", fieldName, field.Type().String())
				}
			} else {
				setFields = append(setFields, fieldName) // This field has a value, save as set
			}
		}
	}

	for parentField, requiredFields := range requiresMap { // Range trough all requires
		for _, requiredField := range requiredFields {
			if !existsIn(setFields, parentField) {
				return fmt.Errorf("field %s requires field %s but is itself empty/not set", requiredField, parentField)
			}
			if !existsIn(setFields, requiredField) { // Check if the requires field was registered as set
				return fmt.Errorf("field %s requires field %s which is empty/not set", parentField, requiredField)
			}
		}
	}
	return nil
}
