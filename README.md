# defcon
Minimalistic library for parsing tagged config structs, automatically handling default and required values

## Overview
defcon is a minimalistic library that parses structs (and nested structs) and examines certain tags, allowing you to tag certain fields with default values, and as required. It was created to ease the pain and repetative nature of sanity checking config structs.
Currently supported types for tagging are all ints, floats and strings.
Allowed tags are `default:"<value>"`, which will set the field to the given value if the value is of the types default value, i.e. 0 for numerical values and empty string for strings, and `required:"true"`, which will cause an error if the the value is of the types default value mentioned earlier.

## Example

To be done