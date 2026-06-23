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

	_ = CheckConfigStruct(&test)

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

	_ = CheckConfigStruct(&test)

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

	CheckConfigStruct(&test)

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
	CheckConfigStruct(&test)
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
	CheckConfigStruct(&test)
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
	CheckConfigStruct(&test)
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
	CheckConfigStruct(&test)
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
	CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Boolean default value failed: %s", err)
	}
	if test.Val != false {
		t.Errorf("Default value did not set correctly. Wanted false, got %t", test.Val)
	}
}

// Test required bool
func TestRequiredBool(t *testing.T) {

	type testStruct struct {
		Val bool `required:"true"`
	}
	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Required boolean field was not required")
	}
}

// Test invalid boolean value
func TestInvalidBoolean(t *testing.T) {

	type testStruct struct {
		Val bool `default:"notabool"`
	}
	test := testStruct{}
	err := CheckConfigStruct(&test)
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
	err = CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Field required by other not detected. %s", err)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Validating required empty struct failed.")
	}
	// Second, test with set values, this should work
	test.Val1 = "set"
	test.Nested.NVal1 = "set"
	err = CheckConfigStruct(&test)
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
	test := testStruct{}
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
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
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Unexported field not handled correctly. %s", err)
	}
	if test.val1 != "test" {
		t.Errorf("Default value in unexported field not handled correctly. Should be 'test', is '%s' Error: %s", test.val1, err)
	}
}

func TestStringArray(t *testing.T) {

	type testStruct struct {
		arr1 []string `default:"{foo, bar}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Default array value not handled correctly. %s", err)
	}
	if test.arr1[0] != "foo" {
		t.Errorf("Default array value not handled correctly. Should be 'foo', is '%s'", test.arr1[0])
	}
	if test.arr1[1] != "bar" {
		t.Errorf("Default array value not handled correctly. Should be 'bar', is '%s'", test.arr1[1])
	}
}

func TestIntArray(t *testing.T) {

	type testStruct struct {
		arr1 []int `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Default array value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default array value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default array value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}
func TestInt8Array(t *testing.T) {

	type testStruct struct {
		arr1 []int8 `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Default array value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default array value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default array value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

func TestInt16Array(t *testing.T) {
	type testStruct struct {
		arr1 []int16 `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Default array value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default array value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default array value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

func TestInt32Array(t *testing.T) {
	type testStruct struct {
		arr1 []int32 `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Default array value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default array value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default array value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

func TestInt64Array(t *testing.T) {
	type testStruct struct {
		arr1 []int64 `default:"{1, 2}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Default array value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1 {
		t.Errorf("Default array value not handled correctly. Should be 'foo', is '%d'", test.arr1[0])
	}
	if test.arr1[1] != 2 {
		t.Errorf("Default array value not handled correctly. Should be 'bar', is '%d'", test.arr1[1])
	}
}

func TestFloat32Array(t *testing.T) {

	type testStruct struct {
		arr1 []float32 `default:"{1.2, 2.3}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Default array value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1.2 {
		t.Errorf("Default array value not handled correctly. Should be 'foo', is '%f'", test.arr1[0])
	}
	if test.arr1[1] != 2.3 {
		t.Errorf("Default array value not handled correctly. Should be 'bar', is '%f'", test.arr1[1])
	}
}

func TestFloat64Array(t *testing.T) {

	type testStruct struct {
		arr1 []float64 `default:"{1.2, 2.3}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Default array value not handled correctly. %s", err)
	}
	if test.arr1[0] != 1.2 {
		t.Errorf("Default array value not handled correctly. Should be 'foo', is '%f'", test.arr1[0])
	}
	if test.arr1[1] != 2.3 {
		t.Errorf("Default array value not handled correctly. Should be 'bar', is '%f'", test.arr1[1])
	}
}

func TestWrongTypeArray(t *testing.T) {

	type testStruct struct {
		arr1 []bool `default:"{foo, bar}"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Wrong type array not handled correctly. %s", err)
	}
}

func TestEnvVarWithoutDefault(t *testing.T) {

	type testStruct struct {
		Val1 string `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "test")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckConfigStruct(&test)
	if test.Val1 != "test" {
		t.Errorf("Env var without default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

func TestEnvVarWithDefault(t *testing.T) {

	type testStruct struct {
		Val1 string `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH" default:"test"`
	}

	test := testStruct{}
	err := CheckConfigStruct(&test)
	if test.Val1 != "test" {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}
}

func TestEnvVarInteger(t *testing.T) {

	type testStruct struct {
		Val1 int `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "123")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckConfigStruct(&test)
	if test.Val1 != 123 {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

func TestEnvVarFloat(t *testing.T) {

	type testStruct struct {
		Val1 float64 `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "123.123")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckConfigStruct(&test)
	if test.Val1 != 123.123 {
		t.Errorf("Env var with default not handled correctly. %s", err)
	}

	err = os.Unsetenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH")
	if err != nil {
		t.Errorf("Could not unset environment variable for test.")
	}
}

func TestEnvVarSlice(t *testing.T) {

	type testStruct struct {
		arr []string `env:"ENV_VAR_ALDKAKDICJAHDNBEBDASH"`
	}

	err := os.Setenv("ENV_VAR_ALDKAKDICJAHDNBEBDASH", "{test1, test2}")
	if err != nil {
		t.Errorf("Could not set environment variable for test.")
	}

	test := testStruct{}
	err = CheckConfigStruct(&test)

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

	err := CheckConfigStruct(&test)
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

	err := CheckConfigStruct(&test)
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

	err = CheckConfigStruct(&test2)
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

	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Required empty slice was not caught")
	}

	// Non-empty slice should pass
	test2 := testStruct{
		Items: []nestedStruct{
			{Name: "item"},
		},
	}

	err = CheckConfigStruct(&test2)
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

	err := CheckConfigStruct(&test)
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

	err := CheckConfigStruct(&test)
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

	err := CheckConfigStruct(&test)
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

	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Empty slice of structs should be valid: %s", err)
	}
}
