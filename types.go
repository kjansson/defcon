package defcon

import (
	"reflect"
	"regexp"
)

// struct field annotations
type annotations struct {
	Required         bool           // Indicates if the field is required
	DefaultValue     string         // Default value for the field if not set
	DefaultFromField string         // Specifies another field from which to derive the default value
	RequiresField    []string       // Specifies another field that must be set if this field is set
	EnvVarName       string         // Name of the environment variable to use for this field
	Unique           bool           // Indicates if the field values must be unique in a slice
	OneOf            string         // Specifies a set of allowed values for the field
	MustMatch        *regexp.Regexp // Specifies a regex pattern that the field value must match
	MustNotMatch     *regexp.Regexp // Specifies a regex pattern that the field value must not match
	MustHave         []string       // Specifies a list of fields that must be present in a slice
	AlwaysHas        []string       // Specifies a list of fields that will always be present in a slice, even if not set
	ValidRange       string         // Specifies a range of allowed values for the field (e.g., "1-10, 44, 100-200")
	ErrorMsg         string         // Custom error message to use when validation fails
}

// common interface for all field types
type field interface {
	handle(*reflect.Value, *annotations) error
}

// getType returns the appropriate field type based on the reflect.Value kind
func getType(v reflect.Value) (field, error) {
	switch v.Kind() {
	case reflect.String:
		return &stringField{}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return &numericField{}, nil
	case reflect.Struct:
		return &structField{}, nil
	case reflect.Bool:
		return &boolField{}, nil
	case reflect.Slice, reflect.Array:
		return &sliceField{}, nil
	default:
		return nil, nil
	}
}
