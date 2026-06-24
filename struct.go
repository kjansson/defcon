package defcon

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type structField struct {
	field reflect.Value
}

func (f *structField) new(v *reflect.Value) {
	f.field = *v
}

func (f *structField) getAnnotations(v reflect.StructField) *annotations {
	var annotations annotations

	// fmt.Println("yay")
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

	return &annotations
}

func (f *structField) handle(a *annotations) error {

	if a == nil {
		a = f.getAnnotations(f.field.Type().Field(0))
	}

	// Manage required
	if a.Required {
		if f.field.IsZero() {
			// Return an error if the field is required but has no value
			return fmt.Errorf("field is marked as required but has no value")
		}
	}

	// Check which fields are set in the struct
	setFields := []string{}
	for i := 0; i < f.field.NumField(); i++ {
		v := f.field.Field(i)
		name := f.field.Type().Field(i).Name
		if !v.IsZero() {
			setFields = append(setFields, name)
		}
	}

	// Handle struct fields recursively
	for i := 0; i < f.field.NumField(); i++ {

		//setFields := []string{}
		var tmpField reflect.Value
		v := f.field.Field(i)
		//name := f.field.Type().Field(i).Name

		if !f.field.Type().Field(i).IsExported() {
			tmpField = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem() // Get access to unexported field
		} else {
			tmpField = v
		}

		subField := tmpField

		fieldType, err := getType(subField)
		if err != nil {
			return fmt.Errorf("failed to get field type: %v", err)
		}

		annotations := f.getAnnotations(f.field.Type().Field(i))

		// Check if the field has a "requires" tag and validate it
		if annotations.RequiresField != "" {
			requiredFields := strings.Split(annotations.RequiresField, ",")
			for _, requiredField := range requiredFields {
				requiredField = strings.TrimSpace(requiredField)
				if !existsIn(setFields, requiredField) {
					return fmt.Errorf("field %s requires field %s to be set", f.field.Type().Field(i).Name, requiredField)
				}
			}
		}

		fieldType.new(&subField)
		err = fieldType.handle(annotations)
		if err != nil {
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
