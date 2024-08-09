package fmap

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var cache = map[reflect.Type]Storage{}

// Get returns a map of field objects.
// It takes a parameter `T` of type `any`, representing the type to be used for Fields map creation.
func Get[T any]() (Storage, error) {
	var tt T
	return getFrom(reflect.TypeOf(tt))
}

type storage struct {
	asMap map[string]Field
	paths []string
}

func (s *storage) Find(path string) (Field, bool) {
	field, ok := s.asMap[path]
	return field, ok
}

func (s *storage) MustFind(path string) Field {
	return s.asMap[path]
}

func (s *storage) GetAllPaths() []string {
	return s.paths
}

// GetFrom returns a map of field objects. It takes a parameter `obj` of type `interface{}` representing the object to be analyzed.
// The function first checks if the `obj` type is already in the cache, and if it exists, it returns the cached value.
// Otherwise, it creates a new empty map with storage.
func GetFrom(obj interface{}) (Storage, error) {
	typeOf := reflect.TypeOf(obj)
	return getFrom(typeOf)
}

func getFrom(typeOf reflect.Type) (Storage, error) {
	if typeOf.Kind() == reflect.Struct {
		typeOf = reflect.PointerTo(typeOf)
	}
	if typeOf.Kind() != reflect.Pointer ||
		(typeOf.Kind() == reflect.Pointer && typeOf.Elem().Kind() != reflect.Struct) {
		return nil, fmt.Errorf("not supported type: %v, only struct and ptr to struct is supported", typeOf)
	}
	if tFields, ok := cache[typeOf]; ok {
		return tFields, nil
	}
	tFields := map[string]Field{}
	count := new(int)
	calculateFields(typeOf, count)
	slice := make([]string, 0, *count)
	getFieldsMapRecursive(typeOf, "", &tFields, &slice, 0)
	cache[typeOf] = &storage{
		asMap: tFields,
		paths: slice,
	}
	return cache[typeOf], nil
}

func calculateFields(confTypeOf reflect.Type, count *int) {
	if confTypeOf.Kind() == reflect.Ptr {
		confTypeOf = confTypeOf.Elem()
	}
	*count += confTypeOf.NumField()
	for i := 0; i < confTypeOf.NumField(); i++ {
		if confTypeOf.Field(i).Type.Kind() == reflect.Struct {
			calculateFields(confTypeOf.Field(i).Type, count)
		}
	}
}

func getFieldsMapRecursive(confTypeOf reflect.Type, path string, f *map[string]Field, s *[]string, offset uintptr) {
	if confTypeOf.Kind() == reflect.Ptr {
		confTypeOf = confTypeOf.Elem()
	}
	if path != "" {
		path += "."
	}
	for i := 0; i < confTypeOf.NumField(); i++ {
		fieldTypeOf := confTypeOf.Field(i)
		var parent *field = nil
		if path != "" {
			parentPath := path[:len(path)-1]
			parent, _ = (*f)[parentPath].(*field)
		}
		switch fieldTypeOf.Type.Kind() {
		case reflect.Struct:
			fld := &field{StructField: fieldTypeOf, structPath: path + fieldTypeOf.Name, parent: parent}
			(*f)[path+fieldTypeOf.Name] = fld
			*s = append(*s, fld.structPath)
			getFieldsMapRecursive(fieldTypeOf.Type, path+fieldTypeOf.Name, f, s, offset+fieldTypeOf.Offset)
		default:
			fld := &field{StructField: fieldTypeOf, structPath: path + fieldTypeOf.Name, parent: parent}
			fld.Offset = fld.Offset + offset
			*s = append(*s, fld.structPath)
			(*f)[path+fieldTypeOf.Name] = fld
		}
	}
}

func (s *storage) GetFieldByPtr(structPtr, fieldPtr any) (Field, error) {
	fldType := reflect.TypeOf(fieldPtr)
	if fldType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not supported type: %v, only ptr to types is supported", fldType)
	}

	if fldType.Kind() == reflect.Ptr {
		fldType = fldType.Elem()
	}

	structType := reflect.TypeOf(structPtr)
	if structType.Kind() != reflect.Ptr ||
		(structType.Kind() == reflect.Ptr && structType.Elem().Kind() != reflect.Struct) {
		return nil, fmt.Errorf("not supported type: %v, only struct and ptr to struct is supported", structType)
	}

	fPtr := ((*[2]unsafe.Pointer)(unsafe.Pointer(&fieldPtr)))[1]
	sPtr := ((*[2]unsafe.Pointer)(unsafe.Pointer(&structPtr)))[1]
	offset := uintptr(fPtr) - uintptr(sPtr)

	for _, path := range s.GetAllPaths() {
		fld := s.MustFind(path)

		if fld.GetOffset() != offset {
			continue
		}

		matchType := fld.GetType().Kind() == fldType.Kind()
		if !matchType {
			continue
		}

		return fld, nil
	}

	return nil, errors.New("field not found")
}
