# defcon
Minimalistic library for parsing tagged config structs, automatically handling default and required values

## Overview
defcon is a minimalistic library that parses structs (and nested structs) and examines certain tags, allowing you to tag certain fields with default values, as required, and mark field as required by others. It was created to ease the pain and repetative nature of sanity checking config structs.
Currently supported types for tagging are all ints, floats and strings. Fields must be exported.  
Allowed tags are `default:"<value>"`, `required:"<true|TRUE>"` and `requires:"field1, field2, ..."`.  
The default tag will modify the struct field with the given value, if the original value is the type default, i.e. zero for numerical values or zero length for strings.  
The required tag will throw an error if the fields value is the type default mentioned earlier.  
The requires tag will throw an error if any of the given fields are missing values when the tagged field is in use.

## Documentation
https://pkg.go.dev/github.com/kjansson/defcon

## Example

```
package main

import (
	"fmt"

	defcon "github.com/kjansson/defcon"
)

type config struct {
	Address  string `default:"localhost"`
	Port     int    `default:"8080"`
	User     string `required:"true"`
	Password string `required:"true"`
}

func main() {

	configuration := config{}
	// Fails if these are empty
	configuration.User = "user"
	configuration.Password = "qwerty"

	err := defcon.CheckConfigStruct(&configuration)
	if err != nil {
		fmt.Println("Parsing error:", err)
	}

	fmt.Println(configuration.Address) // Output: "localhost"
	fmt.Println(configuration.Port)    // Output: "8080"

}
```