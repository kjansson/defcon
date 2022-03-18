# defcon
Minimalistic library for parsing tagged config structs, automatically handling default and required values

## Overview
defcon is a minimalistic library that parses structs (and nested structs) and examines certain tags, allowing you to tag certain fields with default values, and as required. It was created to ease the pain and repetative nature of sanity checking config structs.
Currently supported types for tagging are all ints, floats and strings.
Allowed tags are `default:"<value>"` and `required:"<true|TRUE>"`. The default tag will modify the struct field with the given value, if the original value is the type default, i.e. zero for numerical values or zero length for strings. The required tag will throw an error if the fields value is the type default mentioned earlier.

## Example

To be done