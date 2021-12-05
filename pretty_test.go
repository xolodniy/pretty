package pretty

import (
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"time"
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

func TestEmptyMap(t *testing.T) {
	var m = map[string]interface{}{}
	require.Equal(t, "map{map[]}", Print(m))
}

func TestMap(t *testing.T) {
	var m = map[string]interface{}{
		"int":    1,
		"string": "myString",
	}
	require.Equal(t, "map{map[int:1 string:myString]}", Print(m))
}

func TestInterface(t *testing.T) {
	var i interface{}

	i = 1
	require.Equal(t, "int{1}", Print(i))

	i = "hello world"
	require.Equal(t, "string{hello world}", Print(i))

	i = nil
	require.Equal(t, "nil", Print(i))
}

func TestStructWithPublicFields(t *testing.T) {
	type MyStruct struct {
		A int
		B string
		C int
		D string
	}
	var ms = MyStruct{A: 1, B: "hello", C: 0, D: ""}
	require.Equal(t, "MyStruct{A: 1, B: 'hello'}", Print(ms))
}

func TestStructWithPrivateFields(t *testing.T) {
	type MyStruct struct {
		a int
		b string
		c int
		d string
	}
	var ms = MyStruct{a: 1, b: "hello", c: 0, d: ""}
	require.Equal(t, "MyStruct{a: 1, b: 'hello'}", Print(ms))
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

func TestPointerStruct(t *testing.T) {
	type MyStruct struct {
		myIntFilled int
		myIntEmpty  int
	}
	var (
		ms   = MyStruct{myIntFilled: 1}
		pms  = &ms
		ppms = &pms
	)
	require.Equal(t, "**MyStruct{myIntFilled: 1}", Print(ppms))
}

func TestPointerEmptyString(t *testing.T) {
	var s = ""
	require.Equal(t, "*string{}", Print(&s))
}

func TestEmptyStructFieldsOmitted(t *testing.T) {
	type MyStruct struct {
		a string
		b *string
		c int
		d *int
		e bool
		f *bool
		g int64
		h *int64
		i uint
		j *uint
		k float64
		l *float64
	}
	require.Equal(t, "MyStruct{}", Print(MyStruct{}))
}

func TestFilledStructFieldsPrinted(t *testing.T) {
	type MyStruct struct {
		a string
		b *string
		c int
		d *int
		e bool
		f *bool
		g int64
		h *int64
		i uint
		j *uint
		k float64
		l *float64
	}
	var (
		b         = ""
		d         = 0
		f         = false
		h int64   = 0
		j uint    = 0
		l float64 = 0
	)
	var ms = MyStruct{
		a: "a",
		b: &b,
		c: 1,
		d: &d,
		e: true,
		f: &f,
		g: 1,
		h: &h,
		i: 1,
		j: &j,
		k: 0.1,
		l: &l,
	}
	require.Equal(t, "MyStruct{"+
		"a: 'a', "+
		"b: '', "+
		"c: 1, "+
		"d: 0, "+
		"e: true, "+
		"f: false, "+
		"g: 1, "+
		"h: 0, "+
		"i: 1, "+
		"j: 0, "+
		"k: 0.1, "+
		"l: 0"+
		"}", Print(ms))
}

func TestSliceOfStrings(t *testing.T) {
	var sArr = []string{"1", "2"}
	require.Equal(t, "[]string: [string{1}, string{2}]", Print(sArr))
}

func TestSliceOfInterfaces(t *testing.T) {
	var q, e = 0, 0
	var ww *int = nil
	var arr = []interface{}{q, nil, &e, ww}
	require.Equal(t, "[]interface {}: [int{0}, nil, *int{0}, *int{nil}]", Print(arr))
}

func TestSliceOfStructs(t *testing.T) {
	type Temp struct {
		a int
		b *string
		c bool
	}
	var s = "hello"
	var arr = []Temp{{a: 1, b: &s}, {a: -1, c: true}}
	require.Equal(t, "[]pretty.Temp: [Temp{a: 1, b: 'hello'}, Temp{a: -1, c: true}]", Print(arr))
}

func TestTimeStamp(t *testing.T) {
	asia, err := time.LoadLocation("Asia/Kolkata")
	require.Nil(t, err)
	var date = time.Date(2010, 7, 20, 10, 13, 14, 100000, asia)

	require.Equal(t, "2010-07-20 10:13:14.0001 +0530 IST", Print(date))
	require.Equal(t, "*2010-07-20 10:13:14.0001 +0530 IST", Print(&date))

	var st = struct {
		date  time.Time
		date2 *time.Time
		date3 *time.Time
	}{
		date:  date,
		date2: &date,
	}
	require.Equal(t, "{"+
		"date: unexported date, "+
		"date2: unexported date"+
		"}",
		Print(st))
}

func TestStructWithSqlNullTime(t *testing.T) {
	type testStruct struct {
		gorm.Model
		t     time.Time
		tNil  time.Time
		pT    *time.Time
		pNilT *time.Time
	}
	date := time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)

	require.Equal(t, "testStruct{"+
		"Model: gorm.Model{"+
		"CreatedAt: 2020-01-01T01:01:01Z, "+
		"UpdatedAt: 2020-01-01T01:01:01Z, "+
		"DeletedAt: gorm.DeletedAt{"+
		"Time: 2020-01-01T01:01:01Z"+
		"}}, "+
		"t: unexported date, "+
		"pT: unexported date"+
		"}", Print(testStruct{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: date,
			UpdatedAt: date,
			DeletedAt: gorm.DeletedAt{
				Time:  date,
				Valid: false,
			},
		},
		t:  date,
		pT: &date,
	}))
}

func TestGormModelWithEmptyFields(t *testing.T) {
	type testStruct struct {
		gorm.Model
	}
	require.Equal(t, `testStruct{}`, Print(testStruct{}))
}
