package defcon

import (
	"fmt"
	"os"
	"reflect"
	"slices"

	intervals "github.com/kjansson/go-intervals"
)

type sliceField struct{}

func (f *sliceField) handle(val *reflect.Value, annotations *annotations) error {
	// Check if the slice contains structs
	if val.Len() > 0 && val.Index(0).Kind() == reflect.Struct {

		// Iterate through slice elements and check each struct
		for j := 0; j < val.Len(); j++ {

			element := val.Index(j)
			if element.CanSet() || element.CanAddr() {
				elementPtr := element
				if element.CanAddr() {
					elementPtr = element.Addr().Elem()
				}
				// Determine the type of the slice element
				fieldType, err := getType(elementPtr)
				if err != nil {
					return fmt.Errorf("failed to get field type: %v", err)
				}
				// Handle the field based on its type
				//fieldType.new(&elementPtr)
				err = fieldType.handle(&elementPtr, nil)
				if err != nil {
					return fmt.Errorf("error in slice %s at index %d: %s", val.Type().Name(), j, err)
				}
			}
		}
	} else {

		// Manage env var, default, required for non-struct slices
		if annotations.EnvVarName != "" && val.IsZero() {
			envValue, found := os.LookupEnv(annotations.EnvVarName)
			if found {
				err := setValue(val, envValue)
				if err != nil {
					return fmt.Errorf("failed to set value from environment variable: %v", err)
				}
			}
		}

		// Handle default values
		if annotations.DefaultValue != "" && val.IsZero() {
			err := setValue(val, annotations.DefaultValue)
			if err != nil {
				return fmt.Errorf("failed to set default value: %v", err)
			}
		}

		// Check if slice is required and empty
		if annotations.Required && val.Len() == 0 {
			return fmt.Errorf("field is marked as required but has no value")
		}

	}

	// Check for musthave fields
	for _, mustHaveField := range annotations.MustHave {
		found := false
		for i := 0; i < val.Len(); i++ {
			newVal, err := createTypeFromValue(val.Index(i), mustHaveField)
			if err != nil {
				return fmt.Errorf("error comparing values: %s", err)
			}
			if newVal == val.Index(i) {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("field is marked as must have but has no value for field: %s", mustHaveField)
		}
	}

	// Handle alwayshas
	for _, alwaysHasField := range annotations.AlwaysHas {
		found := false
		for i := 0; i < val.Len(); i++ {
			// fmt.Println("Looking for field:", alwaysHasField, "in slice:", f.field.Index(i))
			newVal, err := createTypeFromValue(val.Index(i), alwaysHasField)
			if err != nil {
				return fmt.Errorf("error comparing values: %s", err)
			}
			if newVal == val.Index(i) {
				found = true
				fmt.Println("Found")
			}
		}
		if !found {

			// Figure out the type of the slice element
			elemType := val.Type().Elem()
			for elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			newPtr := reflect.New(elemType).Elem()

			newVal, err := createTypeFromValue(newPtr, alwaysHasField)
			if err != nil {
				return fmt.Errorf("error comparing values: %s", err)
			}

			val.Set(reflect.Append(*val, newVal))
		}
	}

	// Handle mustmatch
	if annotations.MustMatch != nil && val.Len() > 0 {
		for i := 0; i < val.Len(); i++ {
			if !annotations.MustMatch.MatchString(val.Index(i).String()) {
				return fmt.Errorf("field value '%s' does not match regex '%s'", val.Index(i).String(), annotations.MustMatch)
			}
		}
	}

	// Handle mustnotmatch
	if annotations.MustNotMatch != nil && val.Len() > 0 {
		for i := 0; i < val.Len(); i++ {
			if annotations.MustNotMatch.MatchString(val.Index(i).String()) {
				return fmt.Errorf("field value '%s' matches forbidden regex '%s'", val.Index(i).String(), annotations.MustNotMatch)
			}
		}
	}

	// Handle unique values
	if annotations.Unique && val.Len() > 0 {
		seen := make(map[interface{}]bool)
		for i := 0; i < val.Len(); i++ {
			val := val.Index(i).Interface()
			if seen[val] {
				return fmt.Errorf("field value '%v' is not unique", val)
			}
			seen[val] = true
		}
	}

	// Manage valid range
	if annotations.ValidRange != "" && !val.IsZero() {

		interval, err := intervals.New(annotations.ValidRange)
		if err != nil {
			return fmt.Errorf("failed to create interval: %v", err)
		}

		values := interval.Values()

		for i := 0; i < val.Len(); i++ {
			if !val.Index(i).CanInt() {
				return fmt.Errorf("intervals are only supported on integer fields")
			}

			if !slices.Contains(values, val.Index(i).Int()) {
				return fmt.Errorf("integer value is out of the specified range")
			}
		}

	}

	return nil
}
