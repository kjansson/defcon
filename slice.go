package defcon

import (
	"fmt"
	"os"
	"reflect"
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
				fieldType, err := getType(elementPtr)
				if err != nil {
					return fmt.Errorf("failed to get field type: %v", err)
				}
				fieldType.new(&elementPtr)
				err = fieldType.handle(nil)
				if err != nil {
					return fmt.Errorf("error in slice %s at index %d: %s", f.field.Type().Name(), j, err)
				}
			}
		}
		// Deny
		if f.field.Len() == 0 {
			return fmt.Errorf("field is marked as required but has no value")
		}
	} else {

		// Manage required
		if a.EnvVarName != "" && f.field.IsZero() {
			envValue, found := os.LookupEnv(a.EnvVarName)
			if found {
				err := setValue(&f.field, envValue)
				if err != nil {
					return fmt.Errorf("failed to set value from environment variable: %v", err)
				}
			}
		}
		if a.DefaultValue != "" && f.field.IsZero() {
			err := setValue(&f.field, a.DefaultValue)
			if err != nil {
				return fmt.Errorf("failed to set default value: %v", err)
			}
		} else if a.Required {
			if f.field.IsZero() {
				// Return an error if the field is required but has no value
				return fmt.Errorf("field is marked as required but has no value")
			}
		}

		// Check if slice is required and empty
		if a.Required && f.field.Len() == 0 {
			return fmt.Errorf("field is marked as required but has no value")
		}

	}
	return nil

}
