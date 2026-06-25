package defcon

import (
	"reflect"
)

type annotations struct {
	Required         bool
	DefaultValue     string
	DefaultFromField string
	RequiresField    string
	EnvVarName       string
	Unique           bool
	OneOf            string
	MustMatch        string
	MustNotMatch     string
	MustHave         []string
	AlwaysHas        []string
	ErrorMsg         string
}

type field interface {
	new(*reflect.Value)
	handle(*annotations) error
}

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
