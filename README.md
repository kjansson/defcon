# defcon
Minimalistic library for parsing tagged config structs, automatically handling default values, required values and field dependencies

## Overview
defcon is a minimalistic library that parses structs (and nested structs) and examines certain tags, allowing you to tag certain fields with default values or values from environment variables, mark field as required, or required by other fields. It was created to ease the pain and repetative nature of validating config structs.  

Currently supported types for tagging are all ints, floats and strings and slices. Structs are supported with the required and requires tags.
Allowed tags for primitive types are `default:"<value>"`, `required:"<true|TRUE>"`, `requires:"field1, field2, ..."` and `env:"<envvar_name>"`.  
Allowed tags for slices are `default:"{foo, bar, ...}"` `required:"<true|TRUE>"`, `requires:"field1, field2, ..."` and `env:"<envvar_name>"`. 

## Behaviour
The `default` tag will modify the tagged struct field with the given value, if its original value is the primitive type default, i.e. zero for numerical values, or zero length string.
The `default` tag will be applied before checking if a field is required, so if a field is tagged with both default and required, the required tag will have no effect.
The `env` tag will try to set the value of the tagged field with the value from the given environment variable, wether it has one or not.
The `required` tag will return an error if the tagged fields value is the primitive type default. If applied to a struct, an error will be returned if the tagged struct is considered empty/unset, i.e. if all of its primitive type fields have their default values.  
The `requires` tag will return an error if any of the given fields in the tag value have the primitive type default or is considered an empty struct according the the definition above.  

Tags with invalid values such as references to non-existing fields, values that will overflow the numerical types, invalid numerical values, etc. will result in an error.  


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