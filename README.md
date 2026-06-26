# defcon
Minimalistic library for parsing tagged config structs, handling default values, required values, dependencies and using environment variables.

## Overview
defcon is a minimalistic library that validates and alters struct values using instructions from annotations. It was created to ease the pain and repetative nature of validating config structs.  
Handling default, required, valid and externally fetched values in large structs can cause massive bloat. With defcon you can handle all of this with a simple function call and visible and readable annotations on the struct.

Example - manual validation;
```
type config struct {
	required string
	email string
	foo string
	bar int
	env string
}

func main() {
	c := config{}

	if required == "" {
		panic("this value must be set and i need to explain why")
	}

	regex, err := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %v", err)
	}
	
	if !regex.MatchString(c.email) {
		panic("not a valid email")
	}

	if c.foo == "" {
		c.foo = "default value"
	}

	if c.bar == 0 {
		c.bar == 42
	}

	e, found := os.LookupEnv("FOO_BAR")
	if found {
		c.env = e
	}
}
```

Can be handled using annotations;
```
import "github.com/kjansson/defcon"

type config struct {
	required string `required:"true" errormsg:"this value must be set and i need to explain why"`
	email string `mustmatch:"^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$" errormsg:"not a valid email address"`
	foo string `default:"default value"`
	bar int `default:"42"`
	env string `env`:"FOO_BAR"
}

func main() {
	c := config{}

	err := defcon.CheckStruct(&c)
	if err != nil {
		panic(err)
	}
}
```

## Supported annotations

| Annotation | Example | Target types | Action | Behaviour |
|:---|:---|:---|:---|:---|
| default | `default:"foo"`<br>`default:"{foo, bar}"` | primitives, slices of primitive | correcting | Replaces value if field is unset. |
| required | `required:"true"` | primitives, slices | validating | Returns an error if field is unset. |
| requires | `requires:"field1, field2"` | any struct field | validating | Returns error if field is not unset and any of the given required fields are unset. |
| env | `env:"ENV_VAR_FOO"` | primitives, slices of primitives | altering | Tries to set the field with the value of the given environment variable if found, overwriting the value. |
| defaultfrom | `defaultfrom:"fieldFoo"` | primitives | correcting | Replaces value with the value of another field if annotated field is unset. |
| mustmatch | `mustmatch:"$foo.*^` | strings, slices of strings | validating | Matches the field(s) against the given regular expression, returns error if not matching. |
| mustnotmatch | `mustnotmatch:"$foo.*^` | strings, slices of strings | validating | Matches the field(s) against the given regular expression, returns error if matching. |
| alwayshas | `alwayshas:"foo, bar"` | slices of primitives | correcting | Ensures that a slice always contains a set of given elements. If not present in the slice they will be appended to it. |
| errormsg | `errormsg:"custom error"` | any, in combination with validating annotation | informing | When used with a validating annotation, any validation error will use this error message. |
## Behaviour
- Values from environment variables will be applied before defaults.
- Values from defaults and environment variables takes precedence, i.e. a `required` field as with a `default` value will always be filled in and the `required` check will never fail.

## Documentation
https://pkg.go.dev/github.com/kjansson/defcon

