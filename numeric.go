package defcon

import (
	"fmt"
	"os"
	"reflect"
	"slices"

	intervals "github.com/kjansson/go-intervals"
)

type numericField struct{}

func (f *numericField) handle(val *reflect.Value, annotations *annotations) error {

	// Manage mustmatch

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

	// Manage default values
	if annotations.DefaultValue != "" && val.IsZero() {
		err := setValue(val, annotations.DefaultValue)
		if err != nil {
			return fmt.Errorf("failed to set default value: %v", err)
		}
	}

	// Manage required field
	if annotations.Required {
		if val.IsZero() {
			// Return an error if the field is required but has no value
			return fmt.Errorf("field is marked as required but has no value")
		}
	}

	// Manage valid range
	if annotations.ValidRange != "" && !val.IsZero() {

		if !val.CanInt() {
			return fmt.Errorf("intervals are only supported on integer fields")
		}

		interval, err := intervals.New(annotations.ValidRange)
		if err != nil {
			return fmt.Errorf("failed to create interval: %v", err)
		}

		values := interval.Values()

		if !slices.Contains(values, val.Int()) {
			return fmt.Errorf("integer value is out of the specified range")
		}

	}
	return nil
}
