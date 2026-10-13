package fmap

import (
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
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
		String       string
		StringPtr    *string
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
		String:    "meow",
		StringPtr: &a,
		Age:       17,
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
			name: "exist string ptr field",
			args: args{
				structPtr: &test,
				fieldPtr:  &test.StringPtr,
			},
			want:    s.MustFind("StringPtr"),
			wantErr: false,
		},
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

func TestNestedStructFieldOffset(t *testing.T) {
	type Inner struct {
		At   time.Time
		Name string
	}
	type Middle struct {
		Pad   int64
		Inner Inner
		Deep  struct {
			Pad   bool
			Inner Inner
		}
	}
	type Root struct {
		Pad    string
		First  Inner
		Middle Middle
	}
	storage, err := Get[Root]()
	assert.NoError(t, err)

	root := &Root{}
	offsets := map[string]uintptr{
		"First":                unsafe.Offsetof(root.First),
		"First.At":             unsafe.Offsetof(root.First) + unsafe.Offsetof(root.First.At),
		"Middle":               unsafe.Offsetof(root.Middle),
		"Middle.Inner":         unsafe.Offsetof(root.Middle) + unsafe.Offsetof(root.Middle.Inner),
		"Middle.Inner.At":      unsafe.Offsetof(root.Middle) + unsafe.Offsetof(root.Middle.Inner) + unsafe.Offsetof(root.Middle.Inner.At),
		"Middle.Inner.Name":    unsafe.Offsetof(root.Middle) + unsafe.Offsetof(root.Middle.Inner) + unsafe.Offsetof(root.Middle.Inner.Name),
		"Middle.Deep":          unsafe.Offsetof(root.Middle) + unsafe.Offsetof(root.Middle.Deep),
		"Middle.Deep.Inner":    unsafe.Offsetof(root.Middle) + unsafe.Offsetof(root.Middle.Deep) + unsafe.Offsetof(root.Middle.Deep.Inner),
		"Middle.Deep.Inner.At": unsafe.Offsetof(root.Middle) + unsafe.Offsetof(root.Middle.Deep) + unsafe.Offsetof(root.Middle.Deep.Inner) + unsafe.Offsetof(root.Middle.Deep.Inner.At),
	}
	for path, offset := range offsets {
		assert.Equal(t, offset, storage.MustFind(path).GetOffset(), path)
	}

	date := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	// struct fields are set as a whole value
	storage.MustFind("Middle.Inner.At").Set(root, date)
	assert.Equal(t, date, root.Middle.Inner.At)
	assert.Equal(t, date, storage.MustFind("Middle.Inner.At").Get(root))
	assert.Same(t, &root.Middle.Inner.At, storage.MustFind("Middle.Inner.At").GetPtr(root))

	storage.MustFind("Middle.Deep.Inner").Set(root, Inner{At: date, Name: "deep"})
	assert.Equal(t, Inner{At: date, Name: "deep"}, root.Middle.Deep.Inner)
	assert.Equal(t, Inner{At: date, Name: "deep"}, storage.MustFind("Middle.Deep.Inner").Get(root))

	// other fields are not changed
	assert.Equal(t, Root{Middle: Middle{
		Inner: Inner{At: date},
		Deep: struct {
			Pad   bool
			Inner Inner
		}{Inner: Inner{At: date, Name: "deep"}},
	}}, *root)

	field, err := storage.GetFieldByPtr(root, &root.Middle.Deep.Inner.At)
	assert.NoError(t, err)
	assert.Equal(t, "Middle.Deep.Inner.At", field.GetStructPath())

	// struct has the same offset as its first field
	field, err = storage.GetFieldByPtr(root, &root.Middle.Deep.Inner)
	assert.NoError(t, err)
	assert.Equal(t, "Middle.Deep.Inner", field.GetStructPath())
}

func TestGetFieldByPtr_SameKind(t *testing.T) {
	type Level int
	type Config struct {
		Pad   bool
		Level Level
	}
	conf := &Config{}
	storage, err := GetFrom(conf)
	assert.NoError(t, err)

	// pointer to the underlying type of the field, only kind is matched
	field, err := storage.GetFieldByPtr(conf, (*int)(unsafe.Pointer(&conf.Level)))
	assert.NoError(t, err)
	assert.Equal(t, "Level", field.GetStructPath())

	// field with another kind at the same offset is not matched
	_, err = storage.GetFieldByPtr(conf, (*uint)(unsafe.Pointer(&conf.Level)))
	assert.Error(t, err)
}

func TestGetFrom_Concurrent(t *testing.T) {
	type First struct {
		Name *string
		Sub  struct{ Value **int }
	}
	type Second struct{ First First }
	type Third struct{ Second *Second }

	var wg sync.WaitGroup
	storages := make([]Storage, 16)
	for i := range storages {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for _, obj := range []any{&First{}, Second{}, &Third{}} {
				storage, err := GetFrom(obj)
				assert.NoError(t, err)
				for _, path := range storage.GetAllPaths() {
					field := storage.MustFind(path)
					assert.NotEqual(t, reflect.Ptr, field.GetDereferencedType().Kind())
				}
			}
			storages[i], _ = Get[First]()
		}(i)
	}
	wg.Wait()
	// all goroutines use the same storage
	for _, storage := range storages {
		assert.Same(t, storages[0], storage)
	}
}
