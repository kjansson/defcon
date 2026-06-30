package defcon

import (
	"os"
	"testing"
)

// Test default values for default int
func TestInt(t *testing.T) {

	type testStruct struct {
		Val int `default:"127"`
	}

	test := testStruct{}

	_ = CheckStruct(&test)

	if test.Val != 127 {
		t.Errorf("Default value did not set correctly. Wanted 127, got %d", test.Val)
	}

}

// Test default values for int8
func TestInt8(t *testing.T) {

	type testStruct struct {
		Val int8 `default:"127"`
	}

	test := testStruct{}

	_ = CheckStruct(&test)

	if test.Val != 127 {
		t.Errorf("Default value did not set correctly. Wanted 127, got %d", test.Val)
	}

}

// Test default values for int16
func TestInt16(t *testing.T) {

	type testStruct struct {
		Val int16 `default:"127"`
	}

	test := testStruct{}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Error checking struct: %s", err)
	}

	if test.Val != 127 {
		t.Errorf("Default value did not set correctly. Wanted 127, got %d", test.Val)
	}

}

// Test default values for int32
func TestInt32(t *testing.T) {

	type testStruct struct {
		Val int32 `default:"127"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Error checking struct: %s", err)
	}
	if test.Val != 127 {
		t.Errorf("Default value did not set correctly. Wanted 127, got %d", test.Val)
	}
}

// Test default values for int64
func TestInt64(t *testing.T) {

	type testStruct struct {
		Val int64 `default:"127"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Error checking struct: %s", err)
	}
	if test.Val != 127 {
		t.Errorf("Default value did not set correctly. Wanted 127, got %d", test.Val)
	}
}

// Test default values for float32
func TestFloat32(t *testing.T) {

	type testStruct struct {
		Val float32 `default:"127.127"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Error checking struct: %s", err)
	}
	if test.Val != 127.127 {
		t.Errorf("Default value did not set correctly. Wanted 127.127, got %f", test.Val)
	}
}

// Test default values for float64
func TestFloat64(t *testing.T) {

	type testStruct struct {
		Val float64 `default:"127.127"`
	}
	test := testStruct{}
	CheckStruct(&test)
	if test.Val != 127.127 {
		t.Errorf("Default value did not set correctly. Wanted 127.127, got %f", test.Val)
	}
}

// Test default values for string
func TestString(t *testing.T) {

	type testStruct struct {
		Val string `default:"test"`
	}
	test := testStruct{}
	CheckStruct(&test)
	if test.Val != "test" {
		t.Errorf("Default value did not set correctly. Wanted 'test', got '%s'", test.Val)
	}
}

// Test default values for bool with lowercase "true"
func TestBoolTrue(t *testing.T) {

	type testStruct struct {
		Val bool `default:"true"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Boolean default value failed: %s", err)
	}
	if test.Val != true {
		t.Errorf("Default value did not set correctly. Wanted true, got %t", test.Val)
	}
}

// Test default values for bool with lowercase "false"
func TestBoolFalse(t *testing.T) {

	type testStruct struct {
		Val bool `default:"false"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Boolean default value failed: %s", err)
	}
	if test.Val != false {
		t.Errorf("Default value did not set correctly. Wanted false, got %t", test.Val)
	}
}

// Test default values for bool with uppercase "TRUE"
func TestBoolTRUE(t *testing.T) {

	type testStruct struct {
		Val bool `default:"TRUE"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Boolean default value failed: %s", err)
	}
	if test.Val != true {
		t.Errorf("Default value did not set correctly. Wanted true, got %t", test.Val)
	}
}

// Test default values for bool with uppercase "FALSE"
func TestBoolFALSE(t *testing.T) {

	type testStruct struct {
		Val bool `default:"FALSE"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Boolean default value failed: %s", err)
	}
	if test.Val != false {
		t.Errorf("Default value did not set correctly. Wanted false, got %t", test.Val)
	}
}

// Test default values for bool with mixed case "True"
func TestBoolMixedCaseTrue(t *testing.T) {

	type testStruct struct {
		Val bool `default:"True"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Boolean default value failed: %s", err)
	}
	if test.Val != true {
		t.Errorf("Default value did not set correctly. Wanted true, got %t", test.Val)
	}
}

// Test default values for bool with mixed case "False"
func TestBoolMixedCaseFalse(t *testing.T) {

	type testStruct struct {
		Val bool `default:"False"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Boolean default value failed: %s", err)
	}
	if test.Val != false {
		t.Errorf("Default value did not set correctly. Wanted false, got %t", test.Val)
	}
}

// Test invalid boolean value
func TestInvalidBoolean(t *testing.T) {

	type testStruct struct {
		Val bool `default:"notabool"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Invalid boolean tag value was not detected")
	}
}

// Test bool with env var
func TestEnvVarBool(t *testing.T) {

	type testStruct struct {
		Val bool `env:"ENV_VAR_BOOL_TEST"`
	}

	err := os.Setenv("ENV_VAR_BOOL_TEST", "true")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckStruct(&test)
	if err != nil {
		t.Errorf("Env var bool not handled correctly: %s", err)
	}
	if test.Val != true {
		t.Errorf("Env var bool not handled correctly. Wanted true, got %t", test.Val)
	}

	err = os.Unsetenv("ENV_VAR_BOOL_TEST")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

// Test bool doesn't get overridden when already set to true
func TestBoolAlreadySetTrue(t *testing.T) {

	type testStruct struct {
		Val bool `default:"false"`
	}
	test := testStruct{Val: true}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Boolean handling failed: %s", err)
	}
	if test.Val != true {
		t.Errorf("Boolean value was incorrectly overridden. Wanted true, got %t", test.Val)
	}
}

// Test required string
func TestRequiredString(t *testing.T) {

	type testStruct struct {
		Val string `required:"true"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Required field was not required: %s", err)
	}
}

// Test required int
func TestRequiredInteger(t *testing.T) {

	type testStruct struct {
		Val int `required:"true"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Required field was not required: %s", err)
	}
}

// Test detection of setting an invalid numerical value in int field
func TestInvalidNumerical(t *testing.T) {

	type testStruct struct {
		Val int `default:"somevalue"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Invalid numerical tag value was not detected: %s", err)
	}
}

// Test detection of setting a numerical value too big for its type
func TestOverflowingNumerical(t *testing.T) {

	type testStruct struct {
		Val int8 `default:"555"`
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Overflowing numerical tag value was not detected. %s", err)
	}
}

// Test detection of field requiring other field that is empty
func TestRequires(t *testing.T) {

	type testStruct struct {
		Val1 string `requires:"Val2"`
		Val2 string
		Val3 string
	}
	test := testStruct{Val1: "set"}
	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Field required by other not detected. %s", err)
	}
}

// Test detection of field requiring other field but main field is empty
func TestRequiresPrimaryEmpty(t *testing.T) {

	type testStruct struct {
		Val1 string `requires:"Val2"`
		Val2 string
		Val3 string
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Field requires other field but is not set. %s", err)
	}
}

// Test detection of field requiring multiple other fields that are empty
func TestRequiresMultiple(t *testing.T) {

	type testStruct struct {
		Val1 string `requires:"Val2, Val3"`
		Val2 string
		Val3 string
	}
	test := testStruct{Val1: "set"}
	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Fields required by other not detected. %s", err)
	}
}

// test field requiring other field that is struct
func TestRequiresNestedStruct(t *testing.T) {

	type nestedTestStruct struct {
		NVal1 string
	}
	type testStruct struct {
		Val1   string `requires:"Nested"`
		Nested nestedTestStruct
	}

	// First, test empty value in field that requires other field. This should result in an error.
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Validating required empty struct failed.")
	}
	// Second, test with set values, this should work
	test.Val1 = "set"
	test.Nested.NVal1 = "set"
	err = CheckStruct(&test)
	if err != nil {
		t.Errorf("Validating required nested struct failed.")
	}
}

// Test requires with a invalid field name, should return error
func TestRequiresVarSyntax(t *testing.T) {

	type testStruct struct {
		Val1 string `requires:"Val**2"`
		Val2 string
		Val3 string
	}
	test := testStruct{
		Val1: "set",
	}
	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Illegal characters in field name not detected: %s", err)
	}
}

// Test nested structs
func TestNestedStruct(t *testing.T) {
	type nestedTestStruct struct {
		Val1 string `default:"testnest"`
	}
	type testStruct struct {
		Val1   string `default:"test"`
		Nested nestedTestStruct
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Nested struct not checked correctly. %s", err)
	}
	if test.Nested.Val1 != "testnest" {
		t.Errorf("Nested struct default value not working.")
	}
	if test.Val1 != "test" {
		t.Errorf("Nested struct default value not working.")
	}
}

// Test unexported field, with default value
func TestUnexported(t *testing.T) {

	type testStruct struct {
		val1 string `default:"test"` // Unexported field
	}
	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Unexported field not handled correctly. %s", err)
	}
	if test.val1 != "test" {
		t.Errorf("Default value in unexported field not handled correctly. Should be 'test', is '%s' Error: %s", test.val1, err)
	}
}

// Test default values for string slice
func TestStringSlice(t *testing.T) {

	type testStruct struct {
		arr1 []string `default:"{foo, bar}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Default Slice value not handled correctly. %s", err)
	}
	if test.arr1[0] != "foo" {
		t.Errorf("Default Slice value not handled correctly. Should be 'foo', is '%s'", test.arr1[0])
	}
	if test.arr1[1] != "bar" {
		t.Errorf("Default Slice value not handled correctly. Should be 'bar', is '%s'", test.arr1[1])
	}
}

// Test default values for int slice
func TestIntSlice(t *testing.T) {

	type testStruct struct {
		arr1 []int `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Default Slice value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default Slice value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default Slice value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

// Test default values for int8 slice
func TestInt8Slice(t *testing.T) {

	type testStruct struct {
		arr1 []int8 `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Default Slice value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default Slice value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default Slice value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

// Test default values for int16 slice
func TestInt16Slice(t *testing.T) {
	type testStruct struct {
		arr1 []int16 `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Default Slice value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default Slice value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default Slice value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

// Test default values for int32 slice
func TestInt32Slice(t *testing.T) {
	type testStruct struct {
		arr1 []int32 `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Default Slice value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default Slice value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default Slice value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

// Test default values for int64 slice
func TestInt64Slice(t *testing.T) {
	type testStruct struct {
		arr1 []int64 `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Default Slice value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default Slice value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default Slice value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

// Test default values for float32 slice
func TestFloat32Slice(t *testing.T) {

	type testStruct struct {
		arr1 []float32 `default:"{1.2, 2.3}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Default Slice value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1.2 {
		t.Errorf("Default Slice value not handled correctly. Should be 'foo', is '%f'", test.arr1[0])
	}
	if test.arr1[1] != 2.3 {
		t.Errorf("Default Slice value not handled correctly. Should be 'bar', is '%f'", test.arr1[1])
	}
}

// Test default values for float64 slice
func TestFloat64Slice(t *testing.T) {

	type testStruct struct {
		arr1 []float64 `default:"{1.2, 2.3}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Default Slice value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1.2 {
		t.Errorf("Default Slice value not handled correctly. Should be 'foo', is '%f'", test.arr1[0])
	}
	if test.arr1[1] != 2.3 {
		t.Errorf("Default Slice value not handled correctly. Should be 'bar', is '%f'", test.arr1[1])
	}
}

// Test wrong type slice, should return error
func TestWrongTypeSlice(t *testing.T) {

	type testStruct struct {
		arr1 []bool `default:"{foo, bar}"`
	}

	test := testStruct{}
	err := CheckStruct(&test)

	if test.arr1 != nil {
		t.Errorf("Wrong type Slice not handled correctly. %s", err)
	}

	if err == nil {
		t.Errorf("Wrong type Slice not handled correctly. %s", err)
	}
}

// Test environment variable without default value
func TestEnvVarWithoutDefault(t *testing.T) {

	type testStruct struct {
		Val1 string `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "test")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckStruct(&test)
	if test.Val1 != "test" {
		t.Errorf("Env var without default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

// Test environment variable with default value
func TestEnvVarWithDefault(t *testing.T) {

	type testStruct struct {
		Val1 string `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH" default:"test"`
	}

	test := testStruct{}
	err := CheckStruct(&test)
	if test.Val1 != "test" {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}
}

// Test environment variable with integer value
func TestEnvVarInteger(t *testing.T) {

	type testStruct struct {
		Val1 int `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "123")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckStruct(&test)
	if test.Val1 != 123 {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

// Test environment variable with float value
func TestEnvVarFloat(t *testing.T) {

	type testStruct struct {
		Val1 float64 `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "123.123")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckStruct(&test)
	if test.Val1 != 123.123 {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

// Test environment variable with string slice value
func TestEnvVarStringSlice(t *testing.T) {

	type testStruct struct {
		arr []string `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "{test1, test2}")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckStruct(&test)

	if test.arr[0] != "test1" {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	if test.arr[1] != "test2" {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

// Test environment variable with int slice value
func TestEnvVarIntSlice(t *testing.T) {

	type testStruct struct {
		arr []int `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "{1, 2}")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckStruct(&test)

	if test.arr[0] != 1 {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	if test.arr[1] != 2 {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

// Test environment variable with float slice value
func TestEnvVarFloatSlice(t *testing.T) {

	type testStruct struct {
		arr []float32 `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "{1.1, 2.2}")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckStruct(&test)

	if test.arr[0] != 1.1 {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	if test.arr[1] != 2.2 {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

// Test slice of structs with default values
func TestSliceOfStructs(t *testing.T) {
	type nestedStruct struct {
		Name  string `default:"default_name"`
		Value int    `default:"42"`
	}
	type testStruct struct {
		Items []nestedStruct
	}

	test := testStruct{
		Items: []nestedStruct{
			{},                         // Empty struct, should get defaults
			{Name: "custom"},           // Partially filled, only Value should get default
			{Name: "full", Value: 100}, // Fully filled, no defaults applied
		},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice of structs not handled correctly: %s", err)
	}

	// First element should have all defaults
	if test.Items[0].Name != "default_name" {
		t.Errorf("First element Name should be 'default_name', got '%s'", test.Items[0].Name)
	}
	if test.Items[0].Value != 42 {
		t.Errorf("First element Value should be 42, got %d", test.Items[0].Value)
	}

	// Second element should have custom Name and default Value
	if test.Items[1].Name != "custom" {
		t.Errorf("Second element Name should be 'custom', got '%s'", test.Items[1].Name)
	}
	if test.Items[1].Value != 42 {
		t.Errorf("Second element Value should be 42, got %d", test.Items[1].Value)
	}

	// Third element should keep its values
	if test.Items[2].Name != "full" {
		t.Errorf("Third element Name should be 'full', got '%s'", test.Items[2].Name)
	}
	if test.Items[2].Value != 100 {
		t.Errorf("Third element Value should be 100, got %d", test.Items[2].Value)
	}
}

// Test slice of structs with required fields
func TestSliceOfStructsWithRequired(t *testing.T) {
	type nestedStruct struct {
		Name string `required:"true"`
	}
	type testStruct struct {
		Items []nestedStruct
	}

	// Test with empty field in slice element - should fail
	test := testStruct{
		Items: []nestedStruct{
			{Name: "valid"},
			{}, // Empty Name, should trigger error
		},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Required field in slice element was not caught")
	}

	// Test with all valid fields - should pass
	test2 := testStruct{
		Items: []nestedStruct{
			{Name: "first"},
			{Name: "second"},
		},
	}

	err = CheckStruct(&test2)
	if err != nil {
		t.Errorf("Valid slice of structs with required fields failed: %s", err)
	}
}

// Test required slice of structs (slice itself is required)
func TestRequiredSliceOfStructs(t *testing.T) {
	type nestedStruct struct {
		Name string
	}
	type testStruct struct {
		Items []nestedStruct `required:"true"`
	}

	// Empty slice should fail
	test := testStruct{
		Items: []nestedStruct{},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Required empty slice was not caught")
	}

	// Non-empty slice should pass
	test2 := testStruct{
		Items: []nestedStruct{
			{Name: "item"},
		},
	}

	err = CheckStruct(&test2)
	if err != nil {
		t.Errorf("Valid required slice failed: %s", err)
	}
}

// Test deeply nested structs in slices
func TestSliceOfNestedStructs(t *testing.T) {
	type deepStruct struct {
		DeepValue string `default:"deep_default"`
	}
	type nestedStruct struct {
		Name string `default:"nested_default"`
		Deep deepStruct
	}
	type testStruct struct {
		Items []nestedStruct
	}

	test := testStruct{
		Items: []nestedStruct{
			{}, // Both levels should get defaults
		},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Deeply nested structs in slice not handled correctly: %s", err)
	}

	if test.Items[0].Name != "nested_default" {
		t.Errorf("Nested struct default not applied, got '%s'", test.Items[0].Name)
	}
	if test.Items[0].Deep.DeepValue != "deep_default" {
		t.Errorf("Deep struct default not applied, got '%s'", test.Items[0].Deep.DeepValue)
	}
}

// Test slice of structs with bool fields
func TestSliceOfStructsWithBool(t *testing.T) {
	type nestedStruct struct {
		Enabled bool   `default:"true"`
		Name    string `default:"item"`
	}
	type testStruct struct {
		Items []nestedStruct
	}

	test := testStruct{
		Items: []nestedStruct{
			{},                              // Should get defaults
			{Enabled: true, Name: "custom"}, // Explicitly set to true (non-zero value)
		},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice of structs with bool not handled correctly: %s", err)
	}

	// First element should get defaults
	if test.Items[0].Enabled != true {
		t.Errorf("First element Enabled should be true, got %t", test.Items[0].Enabled)
	}
	if test.Items[0].Name != "item" {
		t.Errorf("First element Name should be 'item', got '%s'", test.Items[0].Name)
	}

	// Second element should keep its values
	if test.Items[1].Enabled != true {
		t.Errorf("Second element Enabled should be true, got %t", test.Items[1].Enabled)
	}
	if test.Items[1].Name != "custom" {
		t.Errorf("Second element Name should be 'custom', got '%s'", test.Items[1].Name)
	}
}

// Test that bool fields with default "true" cannot be set to false
// This is a known limitation: false is the zero value, so it's indistinguishable from "not set"
func TestSliceOfStructsWithBoolLimitation(t *testing.T) {
	type nestedStruct struct {
		Enabled bool `default:"true"`
	}
	type testStruct struct {
		Items []nestedStruct
	}

	// Note: Setting Enabled to false in the struct initializer won't work
	// because false is the zero value - the default will be applied
	test := testStruct{
		Items: []nestedStruct{
			{Enabled: false}, // This will get the default "true" applied
		},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// The field will have the default value, not false
	if test.Items[0].Enabled != true {
		t.Errorf("Field got default value as expected: %t", test.Items[0].Enabled)
	}
}

// Test empty slice (not required)
func TestEmptySliceOfStructs(t *testing.T) {
	type nestedStruct struct {
		Name string `default:"default"`
	}
	type testStruct struct {
		Items []nestedStruct
	}

	test := testStruct{
		Items: []nestedStruct{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Empty slice of structs should be valid: %s", err)
	}
}

// Test slice with musthave
func TestSliceMustHaveRequiredFieldString(t *testing.T) {

	type testStruct struct {
		Items []string `musthave:"item1, item2"`
	}

	test := testStruct{
		Items: []string{"item1", "item3"},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice element missing required string field was not detected")
	}
}

// Test slice with musthave for int
func TestSliceMustHaveRequiredFieldInt(t *testing.T) {

	type testStruct struct {
		Items []int `musthave:"1, 3"`
	}

	test := testStruct{
		Items: []int{1, 2},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice element missing required int field was not detected")
	}
}

// Test slice with musthave for int8
func TestSliceMustHaveRequiredFieldInt8(t *testing.T) {

	type testStruct struct {
		Items []int8 `musthave:"1, 3"`
	}

	test := testStruct{
		Items: []int8{1, 2},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice element missing required int8 field was not detected")
	}
}

// Test slice with musthave for int16
func TestSliceMustHaveRequiredFieldFloat64(t *testing.T) {

	type testStruct struct {
		Items []float64 `musthave:"1.1, 3.3"`
	}

	test := testStruct{
		Items: []float64{1.1, 2.2},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice element missing required float64 field was not detected")
	}
}

// Test slice with musthave for int32
func TestSliceMustHaveRequiredFieldFloat32(t *testing.T) {

	type testStruct struct {
		Items []float32 `musthave:"1.1, 3.3"`
	}

	test := testStruct{
		Items: []float32{1.1, 2.2},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice element missing required float32 field was not detected")
	}
}

// Test string slice with alwayshas
func TestSliceAlwaysHasFieldString(t *testing.T) {

	type testStruct struct {
		Items []string `alwayshas:"item1, item2"`
	}

	test := testStruct{
		Items: []string{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice element with empty slice should be valid: %s", err)
	}

	if len(test.Items) < 2 {
		t.Errorf("Slice element with alwayshas should have required items")
	} else {
		if test.Items[0] != "item1" {
			t.Errorf("Slice element with alwayshas should have first required item")
		}
		if test.Items[1] != "item2" {
			t.Errorf("Slice element with alwayshas should have second required item")
		}
	}
}

// Test int slice with alwayshas
func TestSliceAlwaysHasFieldInt(t *testing.T) {

	type testStruct struct {
		Items []int `alwayshas:"1, 2"`
	}

	test := testStruct{
		Items: []int{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice element with empty slice should be valid: %s", err)
	}

	if len(test.Items) < 2 {
		t.Errorf("Slice element with alwayshas should have required items")
	} else {
		if test.Items[0] != 1 {
			t.Errorf("Slice element with alwayshas should have first required item")
		}
		if test.Items[1] != 2 {
			t.Errorf("Slice element with alwayshas should have second required item")
		}
	}
}

// Test int8 slice with alwayshas
func TestSliceAlwaysHasFieldInt8(t *testing.T) {

	type testStruct struct {
		Items []int8 `alwayshas:"1, 2"`
	}

	test := testStruct{
		Items: []int8{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice element with empty slice should be valid: %s", err)
	}

	if len(test.Items) < 2 {
		t.Errorf("Slice element with alwayshas should have required items")
	} else {
		if test.Items[0] != 1 {
			t.Errorf("Slice element with alwayshas should have first required item")
		}
		if test.Items[1] != 2 {
			t.Errorf("Slice element with alwayshas should have second required item")
		}
	}
}

// Test int16 slice with alwayshas
func TestSliceAlwaysHasFieldInt16(t *testing.T) {

	type testStruct struct {
		Items []int16 `alwayshas:"1, 2"`
	}

	test := testStruct{
		Items: []int16{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice element with empty slice should be valid: %s", err)
	}

	if len(test.Items) < 2 {
		t.Errorf("Slice element with alwayshas should have required items")
	} else {
		if test.Items[0] != 1 {
			t.Errorf("Slice element with alwayshas should have first required item")
		}
		if test.Items[1] != 2 {
			t.Errorf("Slice element with alwayshas should have second required item")
		}
	}
}

// Test int32 slice with alwayshas
func TestSliceAlwaysHasFieldInt32(t *testing.T) {

	type testStruct struct {
		Items []int32 `alwayshas:"1, 2"`
	}

	test := testStruct{
		Items: []int32{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice element with empty slice should be valid: %s", err)
	}

	if len(test.Items) < 2 {
		t.Errorf("Slice element with alwayshas should have required items")
	} else {
		if test.Items[0] != 1 {
			t.Errorf("Slice element with alwayshas should have first required item")
		}
		if test.Items[1] != 2 {
			t.Errorf("Slice element with alwayshas should have second required item")
		}
	}
}

// Test int64 slice with alwayshas
func TestSliceAlwaysHasFieldInt64(t *testing.T) {

	type testStruct struct {
		Items []int64 `alwayshas:"1, 2"`
	}

	test := testStruct{
		Items: []int64{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice element with empty slice should be valid: %s", err)
	}

	if len(test.Items) < 2 {
		t.Errorf("Slice element with alwayshas should have required items")
	} else {
		if test.Items[0] != 1 {
			t.Errorf("Slice element with alwayshas should have first required item")
		}
		if test.Items[1] != 2 {
			t.Errorf("Slice element with alwayshas should have second required item")
		}
	}
}

// Test float32 slice with alwayshas
func TestSliceAlwaysHasFieldFloat32(t *testing.T) {

	type testStruct struct {
		Items []float32 `alwayshas:"1.1, 2.2"`
	}

	test := testStruct{
		Items: []float32{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice element with empty slice should be valid: %s", err)
	}

	if len(test.Items) < 2 {
		t.Errorf("Slice element with alwayshas should have required items")
	} else {
		if test.Items[0] != 1.1 {
			t.Errorf("Slice element with alwayshas should have first required item")
		}
		if test.Items[1] != 2.2 {
			t.Errorf("Slice element with alwayshas should have second required item")
		}
	}
}

// Test float64 slice with alwayshas
func TestSliceAlwaysHasFieldFloat64(t *testing.T) {

	type testStruct struct {
		Items []float64 `alwayshas:"1.1, 2.2"`
	}

	test := testStruct{
		Items: []float64{},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("Slice element with empty slice should be valid: %s", err)
	}

	if len(test.Items) < 2 {
		t.Errorf("Slice element with alwayshas should have required items")
	} else {
		if test.Items[0] != 1.1 {
			t.Errorf("Slice element with alwayshas should have first required item")
		}
		if test.Items[1] != 2.2 {
			t.Errorf("Slice element with alwayshas should have second required item")
		}
	}
}

// Test string field with mustmatch
func TestStringMustMatch(t *testing.T) {

	type testStruct struct {
		Val string `mustmatch:"^test.*"`
	}

	test := testStruct{
		Val: "testValue",
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("String field with mustmatch should be valid: %s", err)
	}
}

// Test string slice with mustmatch
func TestSliceMustMatch(t *testing.T) {

	type testStruct struct {
		Val []string `mustmatch:"^test.*"`
	}

	test := testStruct{
		Val: []string{"testValue", "testAnother"},
	}

	err := CheckStruct(&test)
	if err != nil {
		t.Errorf("String field with mustmatch should be valid: %s", err)
	}
}

// Test string field with mustnotmatch
func TestStringMustNotMatch(t *testing.T) {

	type testStruct struct {
		Val string `mustnotmatch:"^foo.*"`
	}

	test := testStruct{
		Val: "foo",
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("String field with mustnotmatch should not be valid")
	}
}

// Test string slice with mustnotmatch
func TestSliceMustNotMatch(t *testing.T) {

	type testStruct struct {
		Val []string `mustnotmatch:"^foo.*"`
	}

	test := testStruct{
		Val: []string{"foo", "bar"},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice field with mustnotmatch should not be valid")
	}
}

// Test string slice with unique values
func TestSliceUniqueString(t *testing.T) {

	type testStruct struct {
		Val []string `unique:"true"`
	}

	test := testStruct{
		Val: []string{"foo", "foo"},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice field with duplicate values should not be valid")
	}
}

// Test int slice with unique values
func TestSliceUniqueInt(t *testing.T) {

	type testStruct struct {
		Val []int `unique:"true"`
	}

	test := testStruct{
		Val: []int{1, 1},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice field with duplicate values should not be valid")
	}
}

// Test int8 slice with unique values
func TestSliceUniqueFloat32(t *testing.T) {

	type testStruct struct {
		Val []float32 `unique:"true"`
	}

	test := testStruct{
		Val: []float32{1.1, 1.1},
	}

	err := CheckStruct(&test)
	if err == nil {
		t.Errorf("Slice field with duplicate values should not be valid")
	}
}

// Test default and mustmatch together
func TestDefaultAndMustMatch(t *testing.T) {

	type invalid struct {
		Val string `default:"foo" mustmatch:"^test.*"`
	}

	iv := invalid{}

	err := CheckStruct(&iv)
	if err == nil {
		t.Errorf("String field defaulting to value not allowed by mustmatch should not be valid")
	}
}

// Test custom error message for required field
func TestCustomErrorMessageRequired(t *testing.T) {

	type invalid struct {
		Val string `required:"true" errormsg:"Val is required and cannot be empty"`
	}

	iv := invalid{}

	err := CheckStruct(&iv)
	if err == nil {
		t.Errorf("Custom error message should be returned for invalid struct")
	}
	if err.Error() != "Val is required and cannot be empty" {
		t.Errorf("Custom error message not returned correctly. Got: %s", err.Error())
	}
}

// Test custom error message for mustmatch slice
func TestCustomErrorMessageMustMatchSlice(t *testing.T) {

	type invalid struct {
		Val []string `mustmatch:"^foo.*$" errormsg:"Val must match the pattern ^foo.*$"`
	}

	iv := invalid{
		Val: []string{"bar"},
	}

	err := CheckStruct(&iv)
	if err == nil {
		t.Errorf("Custom error message should be returned for invalid struct")
	}
	if err.Error() != "Val must match the pattern ^foo.*$" {
		t.Errorf("Custom error message not returned correctly. Got: %s", err.Error())
	}
}

// Test custom error message for mustmatch string
func TestCustomErrorMessageMustMatch(t *testing.T) {

	type invalid struct {
		Val string `mustmatch:"^test.*" errormsg:"Val must match the pattern ^test.*"`
	}

	iv := invalid{
		Val: "foo",
	}

	err := CheckStruct(&iv)
	if err == nil {
		t.Errorf("Custom error message should be returned for invalid struct")
	}
	if err.Error() != "Val must match the pattern ^test.*" {
		t.Errorf("Custom error message not returned correctly. Got: %s", err.Error())
	}
}

// Test custom error for validrange
func TestIntValidRange(t *testing.T) {

	type invalid struct {
		Val int `validrange:"2-4,99"`
	}

	iv := invalid{
		Val: 1,
	}

	err := CheckStruct(&iv)
	if err == nil {
		t.Errorf("Value %d should be out of the specified valid range", iv.Val)
	}

	iv.Val = 3
	err = CheckStruct(&iv)
	if err != nil {
		t.Errorf("Value %d should be in the specified valid range", iv.Val)
	}

	iv.Val = 99
	err = CheckStruct(&iv)
	if err != nil {
		t.Errorf("Value %d should be in the specified valid range", iv.Val)
	}
}

// Test validrange with invalid type
func TestInvalidTypeValidRange(t *testing.T) {

	type invalid struct {
		Val float32 `validrange:"2-4,99"`
	}

	iv := invalid{
		Val: 2,
	}

	err := CheckStruct(&iv)
	if err == nil {
		t.Errorf("Type float should not be used with valid range")
	}
}

// Test validrange with int slice
func TestIntSliceValidRange(t *testing.T) {

	type invalid struct {
		Val []int `validrange:"2-4,99"`
	}

	iv := invalid{
		Val: []int{1, 2},
	}

	err := CheckStruct(&iv)
	if err == nil {
		t.Errorf("Values should be out of the specified valid range")
	}

	iv.Val = []int{2, 3, 99}
	err = CheckStruct(&iv)
	if err != nil {
		t.Errorf("Value %d should be in the specified valid range", iv.Val)
	}
}

// Test validrange with invalid type slice
func TestInvalidTypeSliceValidRange(t *testing.T) {

	type invalid struct {
		Val []float32 `validrange:"2-4,99"`
	}

	iv := invalid{
		Val: []float32{2},
	}

	err := CheckStruct(&iv)
	if err == nil {
		t.Errorf("Type float should not be used with valid range")
	}
}
