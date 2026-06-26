package defcon

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unsafe"
)

type structField struct {
	field reflect.Value
}

func (f *structField) new(v *reflect.Value) {
	f.field = *v
}

// getAnnotations retrieves the annotations for a struct field
func (f *structField) getAnnotations(v reflect.StructField) *annotations {
	var annotations annotations

	required, found := v.Tag.Lookup("required")
	if found && isTrue(required) {
		annotations.Required = true
	}
	annotations.DefaultValue, found = v.Tag.Lookup("default")
	annotations.DefaultFromField, found = v.Tag.Lookup("defaultfrom")
	annotations.RequiresField, found = v.Tag.Lookup("requires")
	envVar, found := v.Tag.Lookup("env")
	if found {
		annotations.EnvVarName = strings.TrimSpace(envVar)
	}
	mustHave, found := v.Tag.Lookup("musthave")
	if found {
		annotations.MustHave = strings.Split(mustHave, ",")
	}
	unique, found := v.Tag.Lookup("unique")
	if found && !isTrue(unique) {
		annotations.Unique = false
	}
	alwaysHas, found := v.Tag.Lookup("alwayshas")
	if found {
		annotations.AlwaysHas = strings.Split(alwaysHas, ",")
	}
	mustMatch, found := v.Tag.Lookup("mustmatch")
	if found {
		annotations.MustMatch = mustMatch
	}
	mustNotMatch, found := v.Tag.Lookup("mustnotmatch")
	if found {
		annotations.MustNotMatch = mustNotMatch
	}
	unique, found = v.Tag.Lookup("unique")
	if found && isTrue(unique) {
		annotations.Unique = true
	}
	validRange, found := v.Tag.Lookup("validrange")
	if found {
		annotations.ValidRange = validRange
	}
	errMsg, found := v.Tag.Lookup("errormsg")
	if found {
		annotations.ErrorMsg = errMsg
	}

	return &annotations
}

func (f *structField) handle(a *annotations) error {

	if a == nil {
		a = f.getAnnotations(f.field.Type().Field(0))
	}

	// Check which fields are set in the struct and store them for validation of "requires" tags
	setFields := []string{}
	for i := 0; i < f.field.NumField(); i++ {
		v := f.field.Field(i)
		name := f.field.Type().Field(i).Name
		if !v.IsZero() {
			setFields = append(setFields, name)
		}
	}

	// Iterate struct fields and handle each field recursively
	for i := 0; i < f.field.NumField(); i++ {

		var subField reflect.Value
		v := f.field.Field(i)

		// Create a exported version of the field if it is unexported to allow access to its value
		if !f.field.Type().Field(i).IsExported() {
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
		annotations := f.getAnnotations(f.field.Type().Field(i))

		// Check if the field has a "requires" tag and validate that it is set if the current field is set
		// This needs to be handled on the struct level because the "requires" tag can reference other fields in the same struct
		if annotations.RequiresField != "" && !subField.IsZero() {

			requiredFields := strings.Split(annotations.RequiresField, ",")
			for _, requiredField := range requiredFields {
				requiredField = strings.TrimSpace(requiredField)
				match := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]+$").MatchString(requiredField) // Check if name of required field is valid (no special characters, starts with a letter)
				if !match {
					return fmt.Errorf("field %s tagged as required by field %s does not seem to have a valid name", requiredField, f.field.Type().Field(i).Name)
				}
				if !existsIn(setFields, requiredField) { // Check if the required field is set
					// Use custom error message if provided in the annotations
					if annotations.ErrorMsg != "" {
						return fmt.Errorf("%s", annotations.ErrorMsg)
					}
					return fmt.Errorf("field %s requires field %s to be set", f.field.Type().Field(i).Name, requiredField)
				}
			}
		}

		// Handle the field based on its type
		fieldType.new(&subField)
		err = fieldType.handle(annotations)
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

// Finds a string value in an array of strings
func existsIn(subject []string, searchValue string) bool {

	for _, value := range subject {
		if value == searchValue {
			return true
		}
	}
	return false
}
