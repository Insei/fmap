package fmap

import "reflect"

var cache = map[reflect.Type]map[string]Field{}

// Get returns a map of Field objects.
// It takes a parameter `T` of type `any`, representing the type to be used for Fields map creation.
func Get[T any]() map[string]Field {
	obj := new(T)
	return GetFrom(obj)
}

// GetFrom returns a map of Field objects. It takes a parameter `obj` of type `interface{}` representing the object to be analyzed.
// The function first checks if the `obj` type is already in the cache, and if it exists, it returns the cached value.
// Otherwise, it creates a new empty map with fields.
func GetFrom(obj interface{}) map[string]Field {
	typeOf := reflect.TypeOf(obj)
	if tFields, ok := cache[typeOf]; ok {
		return tFields
	}
	tFields := map[string]Field{}
	getFieldsMapRecursive(obj, "", &tFields, 0)
	cache[typeOf] = tFields
	return tFields
}

func getFieldsMapRecursive(conf any, path string, f *map[string]Field, offset uintptr) {
	typeOf := reflect.TypeOf(conf)
	valueOf := reflect.ValueOf(conf)
	if valueOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	if path != "" {
		path += "."
	}
	for i := 0; i < typeOf.NumField(); i++ {
		fieldTypeOf := typeOf.Field(i)
		fieldValueOf := valueOf.Field(i)
		switch fieldTypeOf.Type.Kind() {
		case reflect.Struct:
			(*f)[path+fieldTypeOf.Name] = Field(fieldTypeOf)
			getFieldsMapRecursive(fieldValueOf.Addr().Interface(), path+fieldTypeOf.Name, f, offset+fieldTypeOf.Offset)
		default:
			fld := Field(fieldTypeOf)
			fld.Offset = fld.Offset + offset
			(*f)[path+fieldTypeOf.Name] = fld
		}
	}
}
