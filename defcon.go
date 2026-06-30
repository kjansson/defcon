// defcon is a minimalistic library for parsing tagged config structs, automatically handling required/default values and field dependencies
package defcon

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// CheckStruct accepts a struct and will validate and alter structs field values according to instructions in their annotations. It will recursively process any containing nested structs an slices of structs.
// Supported annotations and applicable types are;
// "default" - all primitive types and slices of primitives - modifies the struct field with the given value if field is not set
// "required" - all types - returns an error if field is not set
// "env" - all primitive types - modifies struct field with value of environment variable if found
// "requires" - all fields - declares a dependency to another field(s) in the same struct, returns an error if dependent field(s) is not set
// "musthave" - slices of primitives - defines a list of values that must be present in a slice, returns an error if any of the values are not present
// "unique" - slices of primitives - returns an error of the slice contains duplicate values
// "alwayshas" - slices of primitives - modifies a slice to always contain the given values, if not present they will be appended at validation time
// "mustmatch" - strings and slices of strings - returns an error if string(s) do not match the given regular expression
// "mustnotmatch - strings and slices of strings - returns an error if string(s) does match the given regular expression
// "validrange" - integers and slices of integers - returns an error if value(s) are not within the given range, e.g. "1-10, 44, 100-200"
// "errormsg" - all types - allows for a custom error message to be returned if validation fails for the field

func CheckStruct(config interface{}) error {

	s := reflect.ValueOf(config).Elem()

	field := structField{}
	field.new(&s)
	err := field.handle(nil)
	if err != nil {
		return err
	}

	return nil

}

// Get reflection type and returns its type family and number of bits
func getTypeDetails(v reflect.Type) (string, int, error) {

	var bits int
	var err error

	// Extract type family and bits (if any) from reflect type string, e.g. "int32" => ["int", "32"]
	family := regexp.MustCompile("^([a-zA-Z]+)([0-9]*)").FindStringSubmatch(v.Kind().String())

	if family[2] == "" {
		bits = 0
	} else {
		bits, err = strconv.Atoi(family[2])
		if err != nil {
			return "", 0, fmt.Errorf("could not convert type: %s", err)
		}
	}
	return family[1], bits, nil
}

// Sets a value back into reflect.Value by determing it's value and parsing the value from a string
func setValue(v *reflect.Value, val string) error {

	family, bits, err := getTypeDetails(v.Type()) // Get type family and number of bits if applicable
	if err != nil {
		return fmt.Errorf("could not determine type: %s", err)
	}

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
	case "bool":
		boolean, err := strconv.ParseBool(val) // Parse string to bool (accepts "true", "false", "TRUE", "FALSE", "True", "False", "1", "0")
		if err != nil {
			return err
		}
		v.SetBool(boolean)
	case "string":
		v.SetString(val)
	case "slice":

		t := v.Type().Elem() // Get the type of the slice elements
		eType, bits, err := getTypeDetails(t)
		if err != nil {
			return fmt.Errorf("could not determine element type: %s", err)
		}
		var uv = reflect.Value{}
		switch eType {
		case "int":
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
		case "float":
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
		case "string":
			values := strings.Split(regexp.MustCompile("^{(.*)}").FindStringSubmatch(val)[1], ",")
			us := make([]string, 0)    // Create a slice of strings
			u := &us                   // Create a pointer to the slice
			uv := reflect.ValueOf(u)   // Get the reflect value of the pointer
			for _, x := range values { // Loop through all values
				y := reflect.ValueOf(strings.TrimSpace(x))  // Get the reflect value of the value
				uv.Elem().Set(reflect.Append(uv.Elem(), y)) // Append the value to the slice
			}
			v.Set(uv.Elem()) // Set the value of the reflect value to the slice
		default:
			return fmt.Errorf("slice type %s is not supported", eType)
		}
	}
	return nil
}

// createTypeFromValue returns a reflect.Value of same type as typedValue, containing the value from untypedValueString if value is parseable for the type
func createTypeFromValue(typedValue reflect.Value, untypedValueString string) (reflect.Value, error) {

	family, bits, err := getTypeDetails(typedValue.Type()) // Get type family and number of bits if applicable
	if err != nil {
		return reflect.Value{}, fmt.Errorf("could not determine type: %s", err)
	}
	untypedValueString = strings.TrimSpace(untypedValueString)

	// Parse numerical values if needed and set values
	switch family {
	case "int":
		integer, err := strconv.ParseInt(untypedValueString, 10, bits) // Parse string to int
		if err != nil {
			return reflect.Value{}, err
		}

		switch bits {
		case 8:
			return reflect.ValueOf(int8(integer)), nil
		case 16:
			return reflect.ValueOf(int16(integer)), nil
		case 32:
			return reflect.ValueOf(int32(integer)), nil
		case 64:
			return reflect.ValueOf(int64(integer)), nil
		default:
			return reflect.ValueOf(int(integer)), nil
		}
	case "float":
		floating, err := strconv.ParseFloat(untypedValueString, bits) // Parse string to float
		if err != nil {
			return reflect.Value{}, err
		}

		switch bits {
		case 32:
			return reflect.ValueOf(float32(floating)), nil
		case 64:
			return reflect.ValueOf(float64(floating)), nil
		default:
			return reflect.ValueOf(float64(floating)), nil
		}

	case "bool":
		boolean, err := strconv.ParseBool(untypedValueString) // Parse string to bool (accepts "true", "false", "TRUE", "FALSE", "True", "False", "1", "0")
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(boolean), nil

	case "string":
		return reflect.ValueOf(untypedValueString), nil
	default:
		return reflect.Value{}, fmt.Errorf("could not determine type")
	}
}
