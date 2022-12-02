package defcon

import (
	"testing"
)

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

func TestRequiredString(t *testing.T) {

	type testStruct struct {
		Val string `required:"true"`
	}
	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Required field was not required.")
	}
}

func TestRequiredInteger(t *testing.T) {

	type testStruct struct {
		Val int `required:"true"`
	}
	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Required field was not required.")
	}
}

func TestRequiredAndDefault(t *testing.T) {

	type testStruct struct {
		Val int `required:"true" default:"123"`
	}
	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Default and required tagged field was not detected.")
	}
}

func TestInvalidNumerical(t *testing.T) {

	type testStruct struct {
		Val int `default:"somevalue"`
	}
	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Invalid numerical tag value was not detected.")
	}
}

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

func TestRequiresNestedStruct(t *testing.T) {

	type nestedTestStruct struct {
		Val1 string
	}
	type testStruct struct {
		Val1   string `requires:"Nested"`
		Nested nestedTestStruct
	}

	test := testStruct{Val1: "set"}
	test.Nested.Val1 = "set"
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Validating required nested struct failed: %s", err)
	}
}

func TestRequiresVarSyntax(t *testing.T) {

	type testStruct struct {
		Val1 string `requires:"Val**2"`
		Val2 string
		Val3 string
	}
	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err == nil {
		t.Errorf("Illegal characters in field name not detected")
	}
}

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

func TestUnexported(t *testing.T) {

	type testStruct struct {
		val1 string `default:"test"`
	}
	test := testStruct{}
	err := CheckConfigStruct(&test)
	if err != nil {
		t.Errorf("Unexported field not handled correctly. %s", err)
	}
}
