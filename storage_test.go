package fmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFmap(t *testing.T) {
	check := func(fields Storage, path string) bool {
		_, ok := fields.Find(path)
		return ok
	}
	test := func(t *testing.T, fields Storage) {
		assert.Equal(t, check(fields, "String"), true)
		assert.Equal(t, check(fields, "Int"), true)
		assert.Equal(t, check(fields, "Int8"), true)
		assert.Equal(t, check(fields, "Int16"), true)
		assert.Equal(t, check(fields, "Int32"), true)
		assert.Equal(t, check(fields, "Int64"), true)
		assert.Equal(t, check(fields, "Uint"), true)
		assert.Equal(t, check(fields, "Uint8"), true)
		assert.Equal(t, check(fields, "Uint16"), true)
		assert.Equal(t, check(fields, "Uint32"), true)
		assert.Equal(t, check(fields, "Uint64"), true)
		assert.Equal(t, check(fields, "Float32"), true)
		assert.Equal(t, check(fields, "Float64"), true)
		assert.Equal(t, check(fields, "Bool"), true)
		assert.Equal(t, check(fields, "PtrString"), true)
		assert.Equal(t, check(fields, "PtrInt"), true)
		assert.Equal(t, check(fields, "PtrInt8"), true)
		assert.Equal(t, check(fields, "PtrInt16"), true)
		assert.Equal(t, check(fields, "PtrInt32"), true)
		assert.Equal(t, check(fields, "PtrInt64"), true)
		assert.Equal(t, check(fields, "PtrUint"), true)
		assert.Equal(t, check(fields, "PtrUint8"), true)
		assert.Equal(t, check(fields, "PtrUint16"), true)
		assert.Equal(t, check(fields, "PtrUint32"), true)
		assert.Equal(t, check(fields, "PtrUint64"), true)
		assert.Equal(t, check(fields, "PtrFloat32"), true)
		assert.Equal(t, check(fields, "PtrFloat64"), true)
		assert.Equal(t, check(fields, "PtrBool"), true)
		assert.Equal(t, check(fields, "NestedStruct"), true)
		assert.Equal(t, check(fields, "NestedStruct.String"), true)
		assert.Equal(t, check(fields, "NestedStruct.PtrString"), true)
		assert.Equal(t, check(fields, "Slice"), true)
		assert.Equal(t, check(fields, "PtrSlice"), true)
		assert.Equal(t, check(fields, "Time"), true)
		assert.Equal(t, check(fields, "PtrTime"), true)
	}

	t.Run("GetFrom", func(t *testing.T) {
		fields, _ := GetFrom(&TestStruct{})
		test(t, fields)
	})
	t.Run("Get_Struct", func(t *testing.T) {
		fields, _ := Get[TestStruct]()
		test(t, fields)
	})
	t.Run("Get_PtrStruct", func(t *testing.T) {
		fields, _ := Get[*TestStruct]()
		test(t, fields)
	})
	t.Run("Get_NotAStruct", func(t *testing.T) {
		fields, err := Get[[]string]()
		assert.Error(t, err)
		assert.Nil(t, fields)
	})
}

func BenchmarkGetFrom(b *testing.B) {
	tt := TestStruct{}
	for i := 0; i < b.N; i++ {
		fields, _ := GetFrom(&tt)
		_ = fields
	}
}

func BenchmarkFieldGet(b *testing.B) {
	tt := TestStruct{}
	fields, _ := GetFrom(&tt)

	for i := 0; i < b.N; i++ {
		val := fields.MustFind("Int").Get(&tt)
		_ = val
	}
}

func BenchmarkFieldGetPtr(b *testing.B) {
	tt := TestStruct{}
	fields, _ := GetFrom(&tt)

	for i := 0; i < b.N; i++ {
		val := fields.MustFind("Int").GetPtr(&tt)
		_ = val
	}
}

func BenchmarkFieldSet(b *testing.B) {
	tt := TestStruct{}
	fields, _ := GetFrom(&tt)

	for i := 0; i < b.N; i++ {
		fields.MustFind("Time").Set(&tt, time.Now())
	}
}

func BenchmarkRawFieldGet(b *testing.B) {
	tt := TestStruct{}
	for i := 0; i < b.N; i++ {
		str := tt.String
		_ = str
	}
}
