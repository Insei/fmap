package fmap

import (
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

func TestGetFrom(t *testing.T) {
	t.Run("Get fields map from struct", func(t *testing.T) {
		fieldMap := GetFrom(&TestStruct{})
		assert.Contains(t, fieldMap, "String")
		assert.Contains(t, fieldMap, "Int")
		assert.Contains(t, fieldMap, "Int8")
		assert.Contains(t, fieldMap, "Int16")
		assert.Contains(t, fieldMap, "Int32")
		assert.Contains(t, fieldMap, "Int64")
		assert.Contains(t, fieldMap, "Uint")
		assert.Contains(t, fieldMap, "Uint8")
		assert.Contains(t, fieldMap, "Uint16")
		assert.Contains(t, fieldMap, "Uint32")
		assert.Contains(t, fieldMap, "Uint64")
		assert.Contains(t, fieldMap, "Float32")
		assert.Contains(t, fieldMap, "Float64")
		assert.Contains(t, fieldMap, "Bool")
		assert.Contains(t, fieldMap, "PtrString")
		assert.Contains(t, fieldMap, "PtrInt")
		assert.Contains(t, fieldMap, "PtrInt8")
		assert.Contains(t, fieldMap, "PtrInt16")
		assert.Contains(t, fieldMap, "PtrInt32")
		assert.Contains(t, fieldMap, "PtrInt64")
		assert.Contains(t, fieldMap, "PtrUint")
		assert.Contains(t, fieldMap, "PtrUint8")
		assert.Contains(t, fieldMap, "PtrUint16")
		assert.Contains(t, fieldMap, "PtrUint32")
		assert.Contains(t, fieldMap, "PtrUint64")
		assert.Contains(t, fieldMap, "PtrFloat32")
		assert.Contains(t, fieldMap, "PtrFloat64")
		assert.Contains(t, fieldMap, "PtrBool")
		assert.Contains(t, fieldMap, "NestedStruct")
		assert.Contains(t, fieldMap, "NestedStruct.String")
		assert.Contains(t, fieldMap, "NestedStruct.PtrString")
		assert.Contains(t, fieldMap, "Slice")
		assert.Contains(t, fieldMap, "PtrSlice")
		assert.Contains(t, fieldMap, "Time")
		assert.Contains(t, fieldMap, "PtrTime")
	})
}

func TestGet(t *testing.T) {
	t.Run("Get fields from struct", func(t *testing.T) {
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
		fieldsMap := GetFrom(source)
		assert.Equal(t, source.String, fieldsMap["String"].Get(source))
		assert.Equal(t, source.Int, fieldsMap["Int"].Get(source))
		assert.Equal(t, source.Int8, fieldsMap["Int8"].Get(source))
		assert.Equal(t, source.Int16, fieldsMap["Int16"].Get(source))
		assert.Equal(t, source.Int32, fieldsMap["Int32"].Get(source))
		assert.Equal(t, source.Int64, fieldsMap["Int64"].Get(source))
		assert.Equal(t, source.Uint, fieldsMap["Uint"].Get(source))
		assert.Equal(t, source.Uint8, fieldsMap["Uint8"].Get(source))
		assert.Equal(t, source.Uint16, fieldsMap["Uint16"].Get(source))
		assert.Equal(t, source.Uint32, fieldsMap["Uint32"].Get(source))
		assert.Equal(t, source.Uint64, fieldsMap["Uint64"].Get(source))
		assert.Equal(t, source.Float32, fieldsMap["Float32"].Get(source))
		assert.Equal(t, source.Float64, fieldsMap["Float64"].Get(source))
		assert.Equal(t, source.Bool, fieldsMap["Bool"].Get(source))
		assert.Equal(t, source.Time, fieldsMap["Time"].Get(source))
		assert.Equal(t, source.PtrString, fieldsMap["PtrString"].Get(source))
		assert.Equal(t, source.PtrInt, fieldsMap["PtrInt"].Get(source))
		assert.Equal(t, source.PtrInt8, fieldsMap["PtrInt8"].Get(source))
		assert.Equal(t, source.PtrInt16, fieldsMap["PtrInt16"].Get(source))
		assert.Equal(t, source.PtrInt32, fieldsMap["PtrInt32"].Get(source))
		assert.Equal(t, source.PtrInt64, fieldsMap["PtrInt64"].Get(source))
		assert.Equal(t, source.PtrUint, fieldsMap["PtrUint"].Get(source))
		assert.Equal(t, source.PtrUint8, fieldsMap["PtrUint8"].Get(source))
		assert.Equal(t, source.PtrUint16, fieldsMap["PtrUint16"].Get(source))
		assert.Equal(t, source.PtrUint32, fieldsMap["PtrUint32"].Get(source))
		assert.Equal(t, source.PtrUint64, fieldsMap["PtrUint64"].Get(source))
		assert.Equal(t, source.PtrFloat32, fieldsMap["PtrFloat32"].Get(source))
		assert.Equal(t, source.PtrFloat64, fieldsMap["PtrFloat64"].Get(source))
		assert.Equal(t, source.PtrBool, fieldsMap["PtrBool"].Get(source))
		assert.Equal(t, source.PtrTime, fieldsMap["PtrTime"].Get(source))
		assert.Equal(t, source.Slice, fieldsMap["Slice"].Get(source))
		assert.Equal(t, source.NestedStruct.String, fieldsMap["NestedStruct.String"].Get(source))
		assert.Equal(t, source.NestedStruct.PtrString, fieldsMap["NestedStruct.PtrString"].Get(source))
	})
}

func TestSet(t *testing.T) {
	t.Run("Set fields in struct", func(t *testing.T) {
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

		fieldMap := GetFrom(dest)
		fieldMap["String"].Set(dest, strVal)
		fieldMap["Int"].Set(dest, intVal)
		fieldMap["Int8"].Set(dest, int8Val)
		fieldMap["Int16"].Set(dest, int16Val)
		fieldMap["Int32"].Set(dest, int32Val)
		fieldMap["Int64"].Set(dest, int64Val)
		fieldMap["Uint"].Set(dest, uintVal)
		fieldMap["Uint8"].Set(dest, uint8Val)
		fieldMap["Uint16"].Set(dest, uint16Val)
		fieldMap["Uint32"].Set(dest, uint32Val)
		fieldMap["Uint64"].Set(dest, uint64Val)
		fieldMap["Float32"].Set(dest, f32Val)
		fieldMap["Float64"].Set(dest, f64Val)
		fieldMap["Bool"].Set(dest, boolVal)
		fieldMap["Time"].Set(dest, timeVal)
		fieldMap["Slice"].Set(dest, slice)
		fieldMap["PtrString"].Set(dest, &strVal)
		fieldMap["PtrInt"].Set(dest, &intVal)
		fieldMap["PtrInt8"].Set(dest, &int8Val)
		fieldMap["PtrInt16"].Set(dest, &int16Val)
		fieldMap["PtrInt32"].Set(dest, &int32Val)
		fieldMap["PtrInt64"].Set(dest, &int64Val)
		fieldMap["PtrUint"].Set(dest, &uintVal)
		fieldMap["PtrUint8"].Set(dest, &uint8Val)
		fieldMap["PtrUint16"].Set(dest, &uint16Val)
		fieldMap["PtrUint32"].Set(dest, &uint32Val)
		fieldMap["PtrUint64"].Set(dest, &uint64Val)
		fieldMap["PtrFloat32"].Set(dest, &f32Val)
		fieldMap["PtrFloat64"].Set(dest, &f64Val)
		fieldMap["PtrBool"].Set(dest, &boolVal)
		fieldMap["PtrTime"].Set(dest, &timeVal)
		fieldMap["PtrSlice"].Set(dest, &slice)
		fieldMap["NestedStruct.String"].Set(dest, nestedStructString)
		fieldMap["NestedStruct.PtrString"].Set(dest, &nestedStructString)

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

func TestGetPtr(t *testing.T) {
	t.Run("Get pointer to fields from struct", func(t *testing.T) {
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
		fieldsMap := GetFrom(source)
		overrideString := "override"

		stringField := fieldsMap["String"].GetPtr(source).(*string)
		*stringField = overrideString
		assert.Equal(t, source.String, *stringField)
		assert.NotEqual(t, source.String, strVal)

		stringPointerField := fieldsMap["PtrString"].GetPtr(source).(**string)
		*stringPointerField = &overrideString
		assert.Equal(t, source.PtrString, *stringPointerField)
		assert.NotEqual(t, source.PtrString, &strVal)

		nestedStructStringFiled := fieldsMap["NestedStruct.String"].GetPtr(source).(*string)
		*nestedStructStringFiled = overrideString
		assert.Equal(t, source.NestedStruct.String, *nestedStructStringFiled)
		assert.NotEqual(t, source.NestedStruct.String, strVal)

		nestedStructPtrStringFiled := fieldsMap["NestedStruct.PtrString"].GetPtr(source).(**string)
		*nestedStructPtrStringFiled = &overrideString
		assert.Equal(t, source.NestedStruct.PtrString, *nestedStructPtrStringFiled)
		assert.NotEqual(t, source.NestedStruct.PtrString, &strVal)
	})
}
