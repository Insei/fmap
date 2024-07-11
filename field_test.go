package fmap

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type NestedStruct struct {
	String    string
	PtrString *string
}

type TestStruct struct {
	String       string
	Int          int
	Int8         int8
	Int16        int16
	Int32        int32
	Int64        int64
	Uint         uint
	Uint8        uint8
	Uint16       uint16
	Uint32       uint32
	Uint64       uint64
	Float32      float32
	Float64      float64
	Bool         bool
	Time         time.Time
	PtrTime      *time.Time
	Slice        []string
	PtrString    *string
	PtrInt       *int
	PtrInt8      *int8
	PtrInt16     *int16
	PtrInt32     *int32
	PtrInt64     *int64
	PtrUint      *uint
	PtrUint8     *uint8
	PtrUint16    *uint16
	PtrUint32    *uint32
	PtrUint64    *uint64
	PtrFloat32   *float32
	PtrFloat64   *float64
	PtrBool      *bool
	NestedStruct NestedStruct
	PtrSlice     *[]string
}

func TestField_Get(t *testing.T) {
	t.Run("Get storage from struct", func(t *testing.T) {
		strVal := "Test2"
		intVal := 1
		int8Val := int8(2)
		int16Val := int16(3)
		int32Val := int32(4)
		int64Val := int64(5)
		uintVal := uint(6)
		uint8Val := uint8(7)
		uint16Val := uint16(8)
		uint32Val := uint32(9)
		uint64Val := uint64(10)
		f32Val := float32(6.6)
		f64Val := float64(7.7)
		boolVal := true
		slice := []string{"Test1", "Test2"}
		timeVal := time.Now()
		source := &TestStruct{
			String:     strVal,
			Int:        intVal,
			Int8:       int8Val,
			Int16:      int16Val,
			Int32:      int32Val,
			Int64:      int64Val,
			Uint:       uintVal,
			Uint8:      uint8Val,
			Uint16:     uint16Val,
			Uint32:     uint32Val,
			Uint64:     uint64Val,
			Float32:    f32Val,
			Float64:    f64Val,
			Bool:       boolVal,
			Time:       timeVal,
			Slice:      slice,
			PtrString:  &strVal,
			PtrInt:     &intVal,
			PtrInt8:    &int8Val,
			PtrInt16:   &int16Val,
			PtrInt32:   &int32Val,
			PtrInt64:   &int64Val,
			PtrUint:    &uintVal,
			PtrUint8:   &uint8Val,
			PtrUint16:  &uint16Val,
			PtrUint32:  &uint32Val,
			PtrUint64:  &uint64Val,
			PtrFloat32: &f32Val,
			PtrFloat64: &f64Val,
			PtrBool:    &boolVal,
			PtrTime:    &timeVal,
			PtrSlice:   &slice,
			NestedStruct: NestedStruct{
				String:    strVal,
				PtrString: &strVal,
			},
		}
		fields, _ := GetFrom(source)
		assert.Equal(t, source.String, fields.MustFind("String").Get(source))
		assert.Equal(t, source.Int, fields.MustFind("Int").Get(source))
		assert.Equal(t, source.Int8, fields.MustFind("Int8").Get(source))
		assert.Equal(t, source.Int16, fields.MustFind("Int16").Get(source))
		assert.Equal(t, source.Int32, fields.MustFind("Int32").Get(source))
		assert.Equal(t, source.Int64, fields.MustFind("Int64").Get(source))
		assert.Equal(t, source.Uint, fields.MustFind("Uint").Get(source))
		assert.Equal(t, source.Uint8, fields.MustFind("Uint8").Get(source))
		assert.Equal(t, source.Uint16, fields.MustFind("Uint16").Get(source))
		assert.Equal(t, source.Uint32, fields.MustFind("Uint32").Get(source))
		assert.Equal(t, source.Uint64, fields.MustFind("Uint64").Get(source))
		assert.Equal(t, source.Float32, fields.MustFind("Float32").Get(source))
		assert.Equal(t, source.Float64, fields.MustFind("Float64").Get(source))
		assert.Equal(t, source.Bool, fields.MustFind("Bool").Get(source))
		assert.Equal(t, source.Time, fields.MustFind("Time").Get(source))
		assert.Equal(t, source.PtrString, fields.MustFind("PtrString").Get(source))
		assert.Equal(t, source.PtrInt, fields.MustFind("PtrInt").Get(source))
		assert.Equal(t, source.PtrInt8, fields.MustFind("PtrInt8").Get(source))
		assert.Equal(t, source.PtrInt16, fields.MustFind("PtrInt16").Get(source))
		assert.Equal(t, source.PtrInt32, fields.MustFind("PtrInt32").Get(source))
		assert.Equal(t, source.PtrInt64, fields.MustFind("PtrInt64").Get(source))
		assert.Equal(t, source.PtrUint, fields.MustFind("PtrUint").Get(source))
		assert.Equal(t, source.PtrUint8, fields.MustFind("PtrUint8").Get(source))
		assert.Equal(t, source.PtrUint16, fields.MustFind("PtrUint16").Get(source))
		assert.Equal(t, source.PtrUint32, fields.MustFind("PtrUint32").Get(source))
		assert.Equal(t, source.PtrUint64, fields.MustFind("PtrUint64").Get(source))
		assert.Equal(t, source.PtrFloat32, fields.MustFind("PtrFloat32").Get(source))
		assert.Equal(t, source.PtrFloat64, fields.MustFind("PtrFloat64").Get(source))
		assert.Equal(t, source.PtrBool, fields.MustFind("PtrBool").Get(source))
		assert.Equal(t, source.PtrTime, fields.MustFind("PtrTime").Get(source))
		assert.Equal(t, source.Slice, fields.MustFind("Slice").Get(source))
		assert.Equal(t, source.NestedStruct.String, fields.MustFind("NestedStruct.String").Get(source))
		assert.Equal(t, source.NestedStruct.PtrString, fields.MustFind("NestedStruct.PtrString").Get(source))
	})
}

func TestField_Set(t *testing.T) {
	t.Run("Set storage in struct", func(t *testing.T) {
		dest := &TestStruct{}

		strVal := "String"
		intVal := 6
		int8Val := int8(7)
		int16Val := int16(8)
		int32Val := int32(9)
		int64Val := int64(10)
		uintVal := uint(6)
		uint8Val := uint8(7)
		uint16Val := uint16(8)
		uint32Val := uint32(9)
		uint64Val := uint64(10)
		f32Val := float32(11.11)
		f64Val := 12.12
		boolVal := true
		timeVal := time.Now()
		nestedStructString := "NewNested"
		slice := []string{"Test1", "Test2"}

		fields, _ := GetFrom(dest)
		fields.MustFind("String").Set(dest, strVal)
		fields.MustFind("Int").Set(dest, intVal)
		fields.MustFind("Int8").Set(dest, int8Val)
		fields.MustFind("Int16").Set(dest, int16Val)
		fields.MustFind("Int32").Set(dest, int32Val)
		fields.MustFind("Int64").Set(dest, int64Val)
		fields.MustFind("Uint").Set(dest, uintVal)
		fields.MustFind("Uint8").Set(dest, uint8Val)
		fields.MustFind("Uint16").Set(dest, uint16Val)
		fields.MustFind("Uint32").Set(dest, uint32Val)
		fields.MustFind("Uint64").Set(dest, uint64Val)
		fields.MustFind("Float32").Set(dest, f32Val)
		fields.MustFind("Float64").Set(dest, f64Val)
		fields.MustFind("Bool").Set(dest, boolVal)
		fields.MustFind("Time").Set(dest, timeVal)
		fields.MustFind("Slice").Set(dest, slice)
		fields.MustFind("PtrString").Set(dest, &strVal)
		fields.MustFind("PtrInt").Set(dest, &intVal)
		fields.MustFind("PtrInt8").Set(dest, &int8Val)
		fields.MustFind("PtrInt16").Set(dest, &int16Val)
		fields.MustFind("PtrInt32").Set(dest, &int32Val)
		fields.MustFind("PtrInt64").Set(dest, &int64Val)
		fields.MustFind("PtrUint").Set(dest, &uintVal)
		fields.MustFind("PtrUint8").Set(dest, &uint8Val)
		fields.MustFind("PtrUint16").Set(dest, &uint16Val)
		fields.MustFind("PtrUint32").Set(dest, &uint32Val)
		fields.MustFind("PtrUint64").Set(dest, &uint64Val)
		fields.MustFind("PtrFloat32").Set(dest, &f32Val)
		fields.MustFind("PtrFloat64").Set(dest, &f64Val)
		fields.MustFind("PtrBool").Set(dest, &boolVal)
		fields.MustFind("PtrTime").Set(dest, &timeVal)
		fields.MustFind("PtrSlice").Set(dest, &slice)
		fields.MustFind("NestedStruct.String").Set(dest, nestedStructString)
		fields.MustFind("NestedStruct.PtrString").Set(dest, &nestedStructString)

		assert.Equal(t, strVal, dest.String)
		assert.Equal(t, intVal, dest.Int)
		assert.Equal(t, int8Val, dest.Int8)
		assert.Equal(t, int16Val, dest.Int16)
		assert.Equal(t, int32Val, dest.Int32)
		assert.Equal(t, int64Val, dest.Int64)
		assert.Equal(t, uintVal, dest.Uint)
		assert.Equal(t, uint8Val, dest.Uint8)
		assert.Equal(t, uint16Val, dest.Uint16)
		assert.Equal(t, uint32Val, dest.Uint32)
		assert.Equal(t, uint64Val, dest.Uint64)
		assert.Equal(t, f32Val, dest.Float32)
		assert.Equal(t, f64Val, dest.Float64)
		assert.Equal(t, boolVal, dest.Bool)
		assert.Equal(t, timeVal, dest.Time)
		assert.Equal(t, slice, dest.Slice)
		assert.Equal(t, &strVal, dest.PtrString)
		assert.Equal(t, &intVal, dest.PtrInt)
		assert.Equal(t, &int8Val, dest.PtrInt8)
		assert.Equal(t, &int16Val, dest.PtrInt16)
		assert.Equal(t, &int32Val, dest.PtrInt32)
		assert.Equal(t, &int64Val, dest.PtrInt64)
		assert.Equal(t, &uintVal, dest.PtrUint)
		assert.Equal(t, &uint8Val, dest.PtrUint8)
		assert.Equal(t, &uint16Val, dest.PtrUint16)
		assert.Equal(t, &uint32Val, dest.PtrUint32)
		assert.Equal(t, &uint64Val, dest.PtrUint64)
		assert.Equal(t, &f32Val, dest.PtrFloat32)
		assert.Equal(t, &f64Val, dest.PtrFloat64)
		assert.Equal(t, &boolVal, dest.PtrBool)
		assert.Equal(t, &slice, dest.PtrSlice)
		assert.Equal(t, nestedStructString, dest.NestedStruct.String)
		assert.Equal(t, &nestedStructString, dest.NestedStruct.PtrString)
	})
}

func TestField_GetPtr(t *testing.T) {
	t.Run("Get pointer to storage from struct", func(t *testing.T) {
		strVal := "Test2"
		intVal := 1
		int8Val := int8(2)
		int16Val := int16(3)
		int32Val := int32(4)
		int64Val := int64(5)
		f32Val := float32(6.6)
		f64Val := float64(7.7)
		boolVal := true
		slice := []string{"Test1", "Test2"}
		source := &TestStruct{
			String:     strVal,
			Int:        intVal,
			Int8:       int8Val,
			Int16:      int16Val,
			Int32:      int32Val,
			Int64:      int64Val,
			Float32:    f32Val,
			Float64:    f64Val,
			Bool:       boolVal,
			Slice:      slice,
			PtrString:  &strVal,
			PtrInt:     &intVal,
			PtrInt8:    &int8Val,
			PtrInt16:   &int16Val,
			PtrInt32:   &int32Val,
			PtrInt64:   &int64Val,
			PtrFloat32: &f32Val,
			PtrFloat64: &f64Val,
			PtrBool:    &boolVal,
			PtrSlice:   &slice,
			NestedStruct: NestedStruct{
				String:    strVal,
				PtrString: &strVal,
			},
		}
		fields, _ := GetFrom(source)
		overrideString := "override"

		stringField := fields.MustFind("String").GetPtr(source).(*string)
		*stringField = overrideString
		assert.Equal(t, source.String, *stringField)
		assert.NotEqual(t, source.String, strVal)

		stringPointerField := fields.MustFind("PtrString").GetPtr(source).(**string)
		*stringPointerField = &overrideString
		assert.Equal(t, source.PtrString, *stringPointerField)
		assert.NotEqual(t, source.PtrString, &strVal)

		nestedStructStringFiled := fields.MustFind("NestedStruct.String").GetPtr(source).(*string)
		*nestedStructStringFiled = overrideString
		assert.Equal(t, source.NestedStruct.String, *nestedStructStringFiled)
		assert.NotEqual(t, source.NestedStruct.String, strVal)

		nestedStructPtrStringFiled := fields.MustFind("NestedStruct.PtrString").GetPtr(source).(**string)
		*nestedStructPtrStringFiled = &overrideString
		assert.Equal(t, source.NestedStruct.PtrString, *nestedStructPtrStringFiled)
		assert.NotEqual(t, source.NestedStruct.PtrString, &strVal)
	})
}

func getMockField(tag reflect.StructTag, parent *field) *field {
	return &field{
		StructField: reflect.StructField{
			Tag: tag,
		},
		parent: parent,
	}
}

func TestField_GetTagPath(t *testing.T) {
	tests := []struct {
		name            string
		field           *field
		tag             string
		ignoreParentTag bool
		want            string
	}{
		{
			"Single Field Tag",
			getMockField(`json:"tag1"`, nil),
			"json",
			false,
			"tag1",
		},
		{
			"Nested Field Tag",
			getMockField(`json:"tag2"`, getMockField(`json:"tag1"`, nil)),
			"json",
			false,
			"tag1.tag2",
		},
		{
			"Missing Parent Tag",
			getMockField(`json:"tag2"`, getMockField(``, nil)),
			"json",
			false,
			"",
		},
		{
			"Missing Current Tag",
			getMockField(``, getMockField(`json:"tag1"`, nil)),
			"json",
			false,
			"",
		},
		{
			"Ignore Missing Parent Tag",
			getMockField(`json:"tag2"`, getMockField(``, nil)),
			"json",
			true,
			"tag2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.GetTagPath(tt.tag, tt.ignoreParentTag); got != tt.want {
				t.Errorf("field.GetTagPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_GetDereferenced(t *testing.T) {
	s := "field1 value"
	type testStruct struct {
		field1 *string
		field2 string
		field3 *int
	}
	testObj := testStruct{
		field1: &s,
		field2: "field2 value",
	}

	rtype := reflect.TypeOf(testObj)
	tests := []struct {
		name     string
		field    *field
		expected bool
	}{
		{
			name: "IsDereferencedPtr",
			field: &field{
				StructField: rtype.Field(0),
			},
			expected: true,
		},
		{
			name: "AlreadyDereferencedType",
			field: &field{
				StructField: rtype.Field(1),
			},
			expected: true,
		},
		{
			name: "NilDereferenceNotOk",
			field: &field{
				StructField: rtype.Field(2),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := tt.field.GetDereferenced(testObj)
			if ok != tt.expected {
				t.Errorf("Got %t, want %t", ok, tt.expected)
			}
		})
	}
}

func TestField_GetDereferencedType(t *testing.T) {
	type someStruct struct {
		MyField int
	}
	dereferencedTypeIntPtrField := &field{StructField: reflect.StructField{
		Type: reflect.TypeOf(new(int)),
	}}

	testCases := []struct {
		name     string
		field    Field
		expected reflect.Type
	}{
		{
			name: "DereferencedTypeInt",
			field: &field{StructField: reflect.StructField{
				Type: reflect.TypeOf(0),
			}},
			expected: reflect.TypeOf(0),
		},
		{
			name:     "DereferencedTypeIntPtr",
			field:    dereferencedTypeIntPtrField,
			expected: reflect.TypeOf(0),
		},
		{
			name:     "DereferencedTypeIntPtrCached",
			field:    dereferencedTypeIntPtrField,
			expected: reflect.TypeOf(0),
		},
		{
			name: "DereferencedTypeStruct",
			field: &field{StructField: reflect.StructField{
				Type: reflect.TypeOf(someStruct{}),
			}},
			expected: reflect.TypeOf(someStruct{}),
		},
		{
			name: "DereferencedTypeStructPtr",
			field: &field{StructField: reflect.StructField{
				Type: reflect.TypeOf(new(someStruct)),
			}},
			expected: reflect.TypeOf(someStruct{}),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.field.GetDereferencedType()
			if actual != testCase.expected {
				t.Errorf("Expected '%v', got '%v'", testCase.expected, actual)
			}
		})
	}
}

func TestField_ReflectStructField(t *testing.T) {
	fields, _ := GetFrom(&struct {
		TestField string
	}{})
	fld := fields.MustFind("TestField")
	fld.GetType()
	fld.GetTag()
	fld.GetParent()
	fld.GetAnonymous()
	fld.GetIndex()
	fld.GetName()
	fld.GetOffset()
	fld.GetPkgPath()
}
