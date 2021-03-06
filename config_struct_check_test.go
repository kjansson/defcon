package defcon

import (
	"fmt"
	"testing"
)

func TestInt(t *testing.T) {

	type testStruct struct {
		Val int `default:"127"`
	}

	test := testStruct{}

	err := CheckConfigStruct(&test)

	if test.Val != 127 {
		fmt.Println(err)
		t.Errorf("Default value did not set correctly. Wanted 127, got %d", test.Val)
	}

}

func TestInt8(t *testing.T) {

	type testStruct struct {
		Val int8 `default:"127"`
	}

	test := testStruct{}

	err := CheckConfigStruct(&test)

	if test.Val != 127 {
		fmt.Println(err)
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
