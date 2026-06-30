package defcon

import (
	"fmt"
	"os"
	"reflect"
)

type stringField struct {
	field reflect.Value
}

func (f *stringField) new(v *reflect.Value) {
	f.field = *v
}

func (f *stringField) handle(a *annotations) error {

	// Lookup environment variable if specified and field is empty
	if a.EnvVarName != "" && f.field.IsZero() {
		envValue, found := os.LookupEnv(a.EnvVarName)
		if found {
			err := setValue(&f.field, envValue)
			if err != nil {
				return fmt.Errorf("failed to set value from environment variable: %v", err)
			}
		}
	}

	// Mangage default value
	if a.DefaultValue != "" && f.field.IsZero() {
		err := setValue(&f.field, a.DefaultValue)
		if err != nil {
			return fmt.Errorf("failed to set default value: %v", err)
		}
	}

	// Manage required field
	if a.Required {
		if f.field.IsZero() {
			// Return an error if the field is required but has no value
			return fmt.Errorf("field is marked as required but has no value")
		}
	}

	// Manage mustmatch
	if a.MustMatch != nil && !f.field.IsZero() {
		if !a.MustMatch.MatchString(f.field.String()) {
			return fmt.Errorf("field value '%s' does not match regex '%s'", f.field.String(), a.MustMatch)
		}
	}

	// Manage mustnotmatch
	if a.MustNotMatch != nil && !f.field.IsZero() {
		if a.MustNotMatch.MatchString(f.field.String()) {
			return fmt.Errorf("field value '%s' matches forbidden regex '%s'", f.field.String(), a.MustNotMatch)
		}
	}

	return nil
}
