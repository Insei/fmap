package fmap

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// cache - storages by the struct type, reflect.Type -> Storage
var cache sync.Map

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
	if cached, ok := cache.Load(typeOf); ok {
		return cached.(Storage), nil
	}
	tFields := map[string]Field{}
	count := new(int)
	calculateFields(typeOf, count)
	slice := make([]string, 0, *count)
	getFieldsMapRecursive(typeOf, "", &tFields, &slice, 0, nil, []reflect.Type{typeOf.Elem()})

	// storage can be already created by another goroutine, in this case the stored one is used
	cached, _ := cache.LoadOrStore(typeOf, &storage{
		asMap: tFields,
		paths: slice,
	})
	return cached.(Storage), nil
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

// getFieldsMapRecursive fills the map and the slice with fields of the struct and all nested structs.
//   - offset: offset of the struct from the root struct or from the struct of the owner.
//   - owner: the nearest field with the pointer to the struct in the path, nil if there are no pointers.
//   - visited: types of the structs in the path, recursive types are not traversed again.
func getFieldsMapRecursive(confTypeOf reflect.Type, path string, f *map[string]Field, s *[]string, offset uintptr,
	owner *field, visited []reflect.Type) {
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
		fld := newField(fieldTypeOf, path+fieldTypeOf.Name, parent)
		// offset from the root struct or from the struct of the owner, not from the parent
		fld.Offset = fld.Offset + offset
		fld.owner = owner
		(*f)[fld.structPath] = fld
		*s = append(*s, fld.structPath)

		structType := fld.dereferenceType
		if isVisited(visited, structType) {
			continue
		}
		switch {
		case fieldTypeOf.Type.Kind() == reflect.Struct:
			getFieldsMapRecursive(structType, fld.structPath, f, s, fld.Offset, owner, append(visited, structType))
		case isStructPtr(fieldTypeOf.Type):
			// struct is stored in another place of the memory, offsets of its fields are calculated from its beginning
			getFieldsMapRecursive(structType, fld.structPath, f, s, 0, fld, append(visited, structType))
		}
	}
}

func isVisited(visited []reflect.Type, typeOf reflect.Type) bool {
	for _, v := range visited {
		if v == typeOf {
			return true
		}
	}
	return false
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

	// struct and its first field has the same address, so the field with the same type is preferred
	var matchKind Field
	for _, path := range s.GetAllPaths() {
		fld := s.asMap[path].(*field)

		// fields of the structs stored by the pointer have no constant offset from the struct, addresses are compared
		// nil is the field of the struct that is not created
		if addr := fld.addr(sPtr, false); addr == nil || addr != fPtr {
			continue
		}

		if fld.GetType() == fldType {
			return fld, nil
		}

		if matchKind == nil && fld.GetType().Kind() == fldType.Kind() {
			matchKind = fld
		}
	}
	if matchKind != nil {
		return matchKind, nil
	}

	return nil, errors.New("field not found")
}
