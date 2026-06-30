# defcon
Minimalistic library for parsing tagged config structs, handling default values, required values, dependencies and using environment variables.

## Overview
defcon is a minimalistic library that validates and alters struct values using instructions from annotations. It was created to ease the pain and repetative nature of validating config structs.  
Handling default, required, valid and externally fetched values in large configuration structs can cause massive bloat. defcon handles all the validation logic for you, and makes your intents visible and easily readable in the struct itself.

Example - manual validation;
```
type config struct {
	required string
	email string
	foo string
	bar int
	env string
	percent int
}

func myFunc(c config) {

	// Check if required field is set
	if required == "" {
		panic("this value must be set and i need to explain why")
	}

	// Validate an email address
	regex, err := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %v", err)
	}	
	if !regex.MatchString(c.email) {
		panic("not a valid email")
	}

	// Apply a default value
	if c.foo == "" {
		c.foo = "default value"
	}

	// Apply a default value
	if c.bar == 0 {
		c.bar == 42
	}

	// Get a value from environment variable
	e, found := os.LookupEnv("FOO_BAR")
	if found {
		c.env = e
	}

	// Check that a value is in a valid range
	if c.percent < 0 || c.percent > 100 {
		panic("percent must be between 0 and 100")
	}

}
```

Using defcon and annotations;
```
import "github.com/kjansson/defcon"

type config struct {
	required string 	`required:"true" errormsg:"this value must be set and i need to explain why"`
	email string 		`mustmatch:"^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$" errormsg:"not a valid email address"`
	foo string 			`default:"default value"`
	bar int 			`default:"42"`
	env string 			`env`:"FOO_BAR"
	percent int 		`validrange:"0-100"`
}

func myfunc(c config) {

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
| alwayshas | `alwayshas:"foo, bar"`<br>`alwayshas:"1,2,3"` | slices of primitives | correcting | Ensures that a slice always contains a set of given elements. If not present in the slice they will be appended to it. |
| validrange | `validrange:"1, 5, 50-100"` | integers, slices of integers | validating | Ensures that the integer value(s) falls within the given range. |
| errormsg | `errormsg:"custom error"` | any, in combination with validating annotation | informing | When used with a validating annotation, any validation error will use this error message. |

## Behaviour
- Values from environment variables will be applied before defaults.
- Values from defaults and environment variables takes precedence, i.e. a `required` field as with a `default` value will always be filled in and the `required` check will never fail.

# Formatting notes

- Boolean values are evaluated with [strconv.ParseBool](https://pkg.go.dev/strconv#ParseBool).
- Regular expressions are evaluated with [regex.Compile](https://pkg.go.dev/regexp#Compile). Backslashes in a regular expressions should be escaped with another backslash, i.e. "\." -> "\\."
- Ranges are expressed with single values (e.g. `1`, `11`, `1024`) and/or ranges (e.g. `10-20`) separated by commas. Example: `"80, 443, 1024-65535"`.
- Whitespace is ignored in all values representing sets of values and ranges.

## Documentation
https://pkg.go.dev/github.com/kjansson/defcon

