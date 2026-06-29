package defcon

import (
	"fmt"
	"os"
	"reflect"
	"slices"

	intervals "github.com/kjansson/go-intervals"
)

type sliceField struct {
	field reflect.Value
}

func (f *sliceField) new(v *reflect.Value) {
	f.field = *v
}

func (f *sliceField) handle(a *annotations) error {
	// Check if the slice contains structs
	if f.field.Len() > 0 && f.field.Index(0).Kind() == reflect.Struct {

		// Iterate through slice elements and check each struct
		for j := 0; j < f.field.Len(); j++ {

			element := f.field.Index(j)
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
				fieldType.new(&elementPtr)
				err = fieldType.handle(nil)
				if err != nil {
					return fmt.Errorf("error in slice %s at index %d: %s", f.field.Type().Name(), j, err)
				}
			}
		}
	} else {

		// Manage env var, default, required for non-struct slices
		if a.EnvVarName != "" && f.field.IsZero() {
			envValue, found := os.LookupEnv(a.EnvVarName)
			if found {
				err := setValue(&f.field, envValue)
				if err != nil {
					return fmt.Errorf("failed to set value from environment variable: %v", err)
				}
			}
		}

		// Handle default values
		if a.DefaultValue != "" && f.field.IsZero() {
			err := setValue(&f.field, a.DefaultValue)
			if err != nil {
				return fmt.Errorf("failed to set default value: %v", err)
			}
		}

		// Check if slice is required and empty
		if a.Required && f.field.Len() == 0 {
			return fmt.Errorf("field is marked as required but has no value")
		}

	}

	// Check for musthave fields
	if len(a.MustHave) > 0 {
		for _, mustHaveField := range a.MustHave {
			found := false
			for i := 0; i < f.field.Len(); i++ {
				val, err := createTypeFromValue(f.field.Index(i), mustHaveField)
				if err != nil {
					return fmt.Errorf("error comparing values: %s", err)
				}
				if val == f.field.Index(i) {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("field is marked as must have but has no value for field: %s", mustHaveField)
			}
		}
	}

	// Handle alwayshas
	if len(a.AlwaysHas) > 0 {

		//var ex reflect.Value
		for _, alwaysHasField := range a.AlwaysHas {
			found := false
			for i := 0; i < f.field.Len(); i++ {
				// fmt.Println("Looking for field:", alwaysHasField, "in slice:", f.field.Index(i))
				val, err := createTypeFromValue(f.field.Index(i), alwaysHasField)
				if err != nil {
					return fmt.Errorf("error comparing values: %s", err)
				}
				if val == f.field.Index(i) {
					found = true
					fmt.Println("Found")
				}
			}
			if !found {

				// Figure out the type of the slice element
				elemType := f.field.Type().Elem()
				for elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}
				newPtr := reflect.New(elemType).Elem()

				val, err := createTypeFromValue(newPtr, alwaysHasField)
				if err != nil {
					return fmt.Errorf("error comparing values: %s", err)
				}

				f.field.Set(reflect.Append(f.field, val))
			}
		}
	}

	// Handle mustmatch
	if a.MustMatch != nil && f.field.Len() > 0 {
		for i := 0; i < f.field.Len(); i++ {
			if !a.MustMatch.MatchString(f.field.Index(i).String()) {
				return fmt.Errorf("field value '%s' does not match regex '%s'", f.field.Index(i).String(), a.MustMatch)
			}
		}
	}

	// Handle mustnotmatch
	if a.MustNotMatch != nil && f.field.Len() > 0 {
		for i := 0; i < f.field.Len(); i++ {
			if a.MustNotMatch.MatchString(f.field.Index(i).String()) {
				return fmt.Errorf("field value '%s' matches forbidden regex '%s'", f.field.Index(i).String(), a.MustNotMatch)
			}
		}
	}

	// Handle unique values
	if a.Unique && f.field.Len() > 0 {
		seen := make(map[interface{}]bool)
		for i := 0; i < f.field.Len(); i++ {
			val := f.field.Index(i).Interface()
			if seen[val] {
				return fmt.Errorf("field value '%v' is not unique", val)
			}
			seen[val] = true
		}
	}

	// Manage valid range
	if a.ValidRange != "" && !f.field.IsZero() {

		interval, err := intervals.New(a.ValidRange)
		if err != nil {
			return fmt.Errorf("failed to create interval: %v", err)
		}

		values := interval.Values()

		for i := 0; i < f.field.Len(); i++ {
			if !f.field.Index(i).CanInt() {
				return fmt.Errorf("intervals are only supported on integer fields")
			}

			if !slices.Contains(values, f.field.Index(i).Int()) {
				return fmt.Errorf("integer value is out of the specified range")
			}
		}

	}

	return nil
}
