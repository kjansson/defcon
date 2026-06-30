package defcon

import (
	"fmt"
	"os"
	"reflect"
)

type boolField struct{}

func (f *boolField) handle(val *reflect.Value, annotations *annotations) error {

	// Manage environment variables
	if annotations.EnvVarName != "" && val.IsZero() {
		envValue, found := os.LookupEnv(annotations.EnvVarName)
		if found {
			err := setValue(val, envValue)
			if err != nil {
				return fmt.Errorf("failed to set value from environment variable: %v", err)
			}
		}
	}

	// Manage default value
	if annotations.DefaultValue != "" && val.IsZero() {
		err := setValue(val, annotations.DefaultValue)
		if err != nil {
			return fmt.Errorf("failed to set default value: %v", err)
		}
	}

	// Manage required
	if annotations.Required && val.IsZero() {
		// Return an error if the field is required but has no value
		return fmt.Errorf("field is marked as required but has no value")
	}

	return nil
}
