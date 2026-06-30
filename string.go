package defcon

import (
	"fmt"
	"os"
	"reflect"
)

type stringField struct{}

func (f *stringField) handle(val *reflect.Value, annotations *annotations) error {

	// Lookup environment variable if specified and field is empty
	if annotations.EnvVarName != "" && val.IsZero() {
		envValue, found := os.LookupEnv(annotations.EnvVarName)
		if found {
			err := setValue(val, envValue)
			if err != nil {
				return fmt.Errorf("failed to set value from environment variable: %v", err)
			}
		}
	}

	// Mangage default value
	if annotations.DefaultValue != "" && val.IsZero() {
		err := setValue(val, annotations.DefaultValue)
		if err != nil {
			return fmt.Errorf("failed to set default value: %v", err)
		}
	}

	// Manage required field
	if annotations.Required {
		if val.IsZero() {
			return fmt.Errorf("field is marked as required but has no value")
		}
	}

	// Manage mustmatch
	if annotations.MustMatch != nil && !val.IsZero() {
		if !annotations.MustMatch.MatchString(val.String()) {
			return fmt.Errorf("field value '%s' does not match regex '%s'", val.String(), annotations.MustMatch)
		}
	}

	// Manage mustnotmatch
	if annotations.MustNotMatch != nil && !val.IsZero() {
		if annotations.MustNotMatch.MatchString(val.String()) {
			return fmt.Errorf("field value '%s' matches forbidden regex '%s'", val.String(), annotations.MustNotMatch)
		}
	}

	return nil
}
