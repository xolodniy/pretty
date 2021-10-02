package pretty

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmptyStruct(t *testing.T) {
	type MyStruct struct{}
	require.Equal(t, "MyStruct{}", Print(MyStruct{}))
}

func TestPrimitiveTypes(t *testing.T) {
	var vTrue bool = true
	require.Equal(t, "bool{true}", Print(vTrue))

	var vFalse bool = false
	require.Equal(t, "bool{false}", Print(vFalse))

	var vNil interface{} = nil
	require.Equal(t, "nil", Print(vNil))

	var vString string = "string"
	require.Equal(t, "string{string}", Print(vString))

	var vInt int = 2
	require.Equal(t, "int{2}", Print(vInt))

	var vInt64 int64 = 2
	require.Equal(t, "int64{2}", Print(vInt64))

	var vFloat32 float32 = 3.2
	require.Equal(t, "float32{3.2}", Print(vFloat32))

	var vFloat64 float64 = 6.4
	require.Equal(t, "float64{6.4}", Print(vFloat64))

	var vUint uint = 3
	require.Equal(t, "uint{3}", Print(vUint))
}

func TestStructWithPublicFields(t *testing.T) {
	type MyStruct struct {
		A int
		B string
		C int
		D string
	}
	var ms = MyStruct{A: 1, B: "hello", C: 0, D: ""}
	require.Equal(t, "MyStruct{A: 1, B: hello}", Print(ms))
}

func TestStructWithPrivateFields(t *testing.T) {
	type MyStruct struct {
		a int
		b string
		c int
		d string
	}
	var ms = MyStruct{a: 1, b: "hello", c: 0, d: ""}
	require.Equal(t, "MyStruct{a: 1, b: hello}", Print(ms))
}

func TestStructWithSubStructs(t *testing.T) {
	type a struct {
		myIntFilled int
		myIntEmpty  int
	}
	type b struct {
		a a
	}
	type c struct {
		b b
	}
	var ms = c{
		b: b{
			a: a{
				myIntFilled: 1,
			},
		},
	}
	require.Equal(t, "c{b: pretty.b{a: pretty.a{myIntFilled: 1}}}", Print(ms))
}
