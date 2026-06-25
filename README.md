# defcon
Minimalistic library for parsing tagged config structs, handling default values, required values, dependencies and using environment variables.

## Overview
defcon is a minimalistic library that validates and alters struct values using instructions from annotations. It was created to ease the pain and repetative nature of validating config structs.  

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
		panic("Required value not set")
	}

	regex, err := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %v", err)
	}
	
	
	if !regex.MatchString(c.email) {
		panic("Not a valid pattern")
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

```
import "github.com/kjansson/defcon"

type config struct {
	required string `required:"true"`
	email string `mustmatch:"^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
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

| Annotation | Example | Target types | Check type | Behaviour |
|:---|:---|:---|:---|:---|
| default | `default:"foo"`<br>`default:"{foo, bar}"` | primitives, slices of primitive | correcting | Replaces value if field has type default. |
| required | `required:"true"` | primitives, slices | validating | Returns an error if field has the type default. |
| requires | `requires:"field1, field2"` | any struct field | validating | Returns error if field is not type default and any of the given fields has type default. |
| env | `env:"ENV_VAR_FOO"` | primitives, slices of primitives | altering | Tries to set the field with the value of the given environment variable if found, overwriting the value. |
| defaultfrom | `defaultfrom:"fieldFoo"` | primitives | correcting | Replaces value with the value of another field if annotated field has type default. |
| mustmatch | `mustmatch:"$foo.*^` | strings, slices of strings | validating | Matches the field(s) against the given regular expression, returns error if not matching. |
| mustnotmatch | `mustnotmatch:"$foo.*^` | strings, slices of strings | validating | Matches the field(s) against the given regular expression, returns error if matching. |
| alwayshas | `alwayshas:"foo, bar"` | slices of primitives | correcting | Ensures that a slice always contains a set of given elements. If not present in the slice they will be appended to it. |

## Behaviour
- Values from environment variables will be applied before defaults.
- Values from defaults and environment variables takes precedence, i.e. a `required` field as with a `default` value will always be filled in and the `required` check will never fail.

## Documentation
https://pkg.go.dev/github.com/kjansson/defcon

## Example

```
package main

import (
	"fmt"

	defcon "github.com/kjansson/defcon"
)

type networkConfig struct {
	Protocol string
}

type config struct {
	Address  string   `default:"localhost"`
	Port     int      `default:"8080" requires:"Network"`
	User     string   `required:"true"`
	Password string   `required:"true"`
	Network  networkConfig	// Implicitly required by field Port
	Options  []string `default:"{foo, bar}"`
	Name     string   `env:"APP_NAME"`
}

func main() {

	_ = os.Setenv("APP_NAME", "myapp") // For illustration purposes, this would normally be set outside of the code 

	configuration := config{}
	// Fails if these are empty
	configuration.User = "user"
	configuration.Password = "qwerty"
	configuration.Network.Protocol = "HTTP"	// Field "Port" also requires the field "Network" which has to have set fields. 

	err := defcon.CheckConfigStruct(&configuration)
	if err != nil {
		fmt.Println("Parsing error:", err)
	}

	fmt.Println(configuration.Address) 	// Output: "localhost"
	fmt.Println(configuration.Port)    	// Output: "8080"
	fmt.Println(configuration.Options)	// Output: "[foo bar]"
	fmt.Println(configuration.Name)     // Output: "myapp"

}
```

Try it out!  
https://go.dev/play/p/6sIpyODLlv7