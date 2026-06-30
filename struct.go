package defcon

import (
	"fmt"
	"go/token"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unsafe"
)

type structField struct{}

// getAnnotations retrieves the annotations for a struct field
func (f *structField) getAnnotations(v reflect.StructField) (*annotations, error) {

	var annotations annotations
	var err error

	required, found := v.Tag.Lookup("required")
	// Get and validate boolean value for required
	if found {
		reqBool, err := strconv.ParseBool(required)
		if err != nil {
			return nil, fmt.Errorf("non-boolean value found where expected: %s", err)
		}
		annotations.Required = reqBool
	}

	// Default, defaultfrom are pure string values, no checks required
	annotations.DefaultValue, _ = v.Tag.Lookup("default")
	annotations.DefaultFromField, _ = v.Tag.Lookup("defaultfrom")

	// Get requires fields, clean up whitespace and split by comma
	requires, found := v.Tag.Lookup("requires")
	if found {
		annotations.RequiresField = strings.Split(strings.TrimSpace(requires), ",")
	}

	// Get and clean up environment variable names
	envVar, found := v.Tag.Lookup("env")
	if found {
		annotations.EnvVarName = strings.TrimSpace(envVar)
	}

	// Get and clean up musthave values
	mustHave, found := v.Tag.Lookup("musthave")
	if found {
		annotations.MustHave = strings.Split(strings.TrimSpace(mustHave), ",")
	}

	// Get and validate boolean value for unique
	unique, found := v.Tag.Lookup("unique")
	if found {
		uniqueBool, err := strconv.ParseBool(unique)
		if err != nil {
			return nil, fmt.Errorf("non-boolean value found where expected: %s", err)
		}
		annotations.Unique = uniqueBool
	}

	// Get and cleanup values for alwayshas
	alwaysHas, found := v.Tag.Lookup("alwayshas")
	if found {
		annotations.AlwaysHas = strings.Split(strings.TrimSpace(alwaysHas), ",")
	}

	// Get and compile regex for mustmatch
	mustMatch, found := v.Tag.Lookup("mustmatch")
	if found {
		annotations.MustMatch, err = regexp.Compile(mustMatch)
		if err != nil {
			return nil, fmt.Errorf("could not parse regular expression: %s", err)
		}
	}

	// Get and compile regex for mustnotmatch
	mustNotMatch, found := v.Tag.Lookup("mustnotmatch")
	if found {
		annotations.MustNotMatch, err = regexp.Compile(mustNotMatch)
		if err != nil {
			return nil, fmt.Errorf("could not parse regular expression: %s", err)
		}
	}

	// Get and validate boolean value for unique
	unique, found = v.Tag.Lookup("unique")
	if found {
		uniqueBool, err := strconv.ParseBool(unique)
		if err != nil {
			return nil, fmt.Errorf("non-boolean value found where expected: %s", err)
		}
		annotations.Unique = uniqueBool
	}

	// Get validrange values, these are validated in the numericField handler
	validRange, found := v.Tag.Lookup("validrange")
	if found {
		annotations.ValidRange = validRange
	}
	errMsg, found := v.Tag.Lookup("errormsg")
	if found {
		annotations.ErrorMsg = errMsg
	}

	return &annotations, nil
}

func (f *structField) handle(val *reflect.Value, annotations *annotations) error {

	// Check which fields are set in the struct and store them for validation of "requires" tags
	setFields := []string{}
	for i := 0; i < val.NumField(); i++ {
		v := val.Field(i)
		name := val.Type().Field(i).Name
		if !v.IsZero() {
			setFields = append(setFields, name)
		}
	}

	// Iterate struct fields and handle each field recursively
	for i := 0; i < val.NumField(); i++ {

		var subField reflect.Value
		v := val.Field(i)

		// Create a exported version of the field if it is unexported to allow access to its value
		if !val.Type().Field(i).IsExported() {
			subField = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem() // Get access to unexported field
		} else {
			subField = v
		}

		// Get the type handler for the current field
		fieldType, err := getType(subField)
		if err != nil {
			return fmt.Errorf("failed to get field type: %v", err)
		}

		// Get annotations for the current field
		annotations, err := f.getAnnotations(val.Type().Field(i))
		if err != nil {
			return fmt.Errorf("invalid annotation syntax: %s", err)
		}

		// Check if the field has a "requires" tag and validate that it is set if the current field is set
		// This needs to be handled on the struct level because the "requires" tag can reference other fields in the same struct
		if len(annotations.RequiresField) > 0 && !subField.IsZero() {

			for _, requiredField := range annotations.RequiresField {

				if !token.IsIdentifier(requiredField) {
					return fmt.Errorf("field %s tagged as required by field %s does not seem to have a valid name", requiredField, val.Type().Field(i).Name)
				}
				if !slices.Contains(setFields, requiredField) { // Check if the required field is set
					// Use custom error message if provided in the annotations
					if annotations.ErrorMsg != "" {
						return fmt.Errorf("%s", annotations.ErrorMsg)
					}
					return fmt.Errorf("field %s requires field %s to be set", val.Type().Field(i).Name, requiredField)
				}
			}
		}

		// Handle the field based on its type
		err = fieldType.handle(&subField, annotations)
		if err != nil {
			// Use custom error message if provided in the annotations
			if annotations.ErrorMsg != "" {
				return fmt.Errorf("%s", annotations.ErrorMsg)
			}
			return err
		}
	}

	return nil
}
