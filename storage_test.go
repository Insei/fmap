package fmap

import (
	"github.com/stretchr/testify/assert"
	"reflect"
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

func Test_storage_GetFieldByPtr(t *testing.T) {
	type TestStruct struct {
		Name         string
		Age          uint
		NestedStruct struct {
			Field    string
			Slice    []string
			SlicePtr []*string
			MapSlice map[string][]string
		}
		Slice    []string
		SlicePtr []*string
		MapSlice map[string][]string
	}

	a, b, c := "a", "b", "c"
	test := TestStruct{
		Name: "meow",
		Age:  17,
		NestedStruct: struct {
			Field    string
			Slice    []string
			SlicePtr []*string
			MapSlice map[string][]string
		}{
			Field:    "f",
			Slice:    []string{"a", "b", "c"},
			SlicePtr: []*string{&a, &b, &c},
			MapSlice: map[string][]string{
				"a": []string{"apricot", "avocado", "apple"},
				"b": []string{"banana", "blackberry", "blueberry"},
				"c": []string{"coconut", "cherry", "cashew"},
			},
		},
		Slice:    []string{"a", "b", "c"},
		SlicePtr: []*string{&a, &b, &c},
		MapSlice: map[string][]string{
			"a": []string{"apricot", "avocado", "apple"},
			"b": []string{"banana", "blackberry", "blueberry"},
			"c": []string{"coconut", "cherry", "cashew"},
		},
	}

	s, _ := GetFrom(test)

	sideValue := 555

	type args struct {
		structPtr any
		fieldPtr  any
	}
	tests := []struct {
		name       string
		getStorage func() Storage
		args       args
		want       Field
		wantErr    bool
	}{
		{
			name: "exist uint field",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.Age,
			},
			want:    s.MustFind("Age"),
			wantErr: false,
		},
		{
			name: "non-exist field",
			args: args{
				structPtr: &test,
				fieldPtr:  &sideValue,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-exist field",
			args: args{
				structPtr: &test,
				fieldPtr:  &sideValue,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "exist slice field",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.Slice,
			},
			want:    s.MustFind("Slice"),
			wantErr: false,
		},
		{
			name: "exist slice ptr field",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.SlicePtr,
			},
			want:    s.MustFind("SlicePtr"),
			wantErr: false,
		},
		{
			name: "exist map slice field",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.MapSlice,
			},
			want:    s.MustFind("MapSlice"),
			wantErr: false,
		},
		{
			name: "exist struct",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.NestedStruct,
			},
			want:    s.MustFind("NestedStruct"),
			wantErr: false,
		},
		{
			name: "exist struct field",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.NestedStruct.Field,
			},
			want:    s.MustFind("NestedStruct.Field"),
			wantErr: false,
		},
		{
			name: "exist struct slice",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.NestedStruct.Slice,
			},
			want:    s.MustFind("NestedStruct.Slice"),
			wantErr: false,
		},
		{
			name: "exist struct slice ptr",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.NestedStruct.SlicePtr,
			},
			want:    s.MustFind("NestedStruct.SlicePtr"),
			wantErr: false,
		},
		{
			name: "exist struct map slice",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.NestedStruct.MapSlice,
			},
			want:    s.MustFind("NestedStruct.MapSlice"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFieldByPtr(tt.args.structPtr, tt.args.fieldPtr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldByPtr() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldByPtr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
