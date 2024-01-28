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
