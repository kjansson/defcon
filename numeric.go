package defcon

import (
	"fmt"
	"os"
	"reflect"
)

type numericField struct {
	field reflect.Value
}

func (f *numericField) new(v *reflect.Value) {
	f.field = *v
}

func (f *numericField) handle(a *annotations) error {

	// Manage environment variables
	if a.EnvVarName != "" && f.field.IsZero() {
		envValue, found := os.LookupEnv(a.EnvVarName)
		if found {
			err := setValue(&f.field, envValue)
			if err != nil {
				return fmt.Errorf("failed to set value from environment variable: %v", err)
			}
		}
	}

	// Manage default values
	if a.DefaultValue != "" && f.field.IsZero() {
		err := setValue(&f.field, a.DefaultValue)
		if err != nil {
			return fmt.Errorf("failed to set default value: %v", err)
		}
	}

	// Manage required
	if a.Required {
		if f.field.IsZero() {
			// Return an error if the field is required but has no value
			return fmt.Errorf("field is marked as required but has no value")
		}
	}
	return nil
}
