package fmap

import (
	"reflect"
	"unsafe"
)

type Field reflect.StructField

func (f Field) GetPtr(obj interface{}) interface{} {
	return reflect.NewAt(f.Type, f.getPtr(obj)).Interface()
}

// Get returns the value of the fields in the provided object.
// It takes a parameter `obj` of type `interface{}`, representing the object.
// It returns the value of the fields as an `interface{}`.
func (f Field) Get(obj interface{}) interface{} {
	ptrToField := f.getPtr(obj)
	kind := f.Type.Kind()
	isPtr := false
	if kind == reflect.Ptr {
		isPtr = true
		kind = f.Type.Elem().Kind()
	}
	if isPtr {
		switch kind {
		case reflect.String:
			return getPtrValue[*string](ptrToField)
		case reflect.Int:
			return getPtrValue[*int](ptrToField)
		case reflect.Int8:
			return getPtrValue[*int8](ptrToField)
		case reflect.Int16:
			return getPtrValue[*int16](ptrToField)
		case reflect.Int32:
			return getPtrValue[*int32](ptrToField)
		case reflect.Int64:
			return getPtrValue[*int64](ptrToField)
		case reflect.Uint:
			return getPtrValue[*uint](ptrToField)
		case reflect.Uint8:
			return getPtrValue[*uint8](ptrToField)
		case reflect.Uint16:
			return getPtrValue[*uint16](ptrToField)
		case reflect.Uint32:
			return getPtrValue[*uint32](ptrToField)
		case reflect.Uint64:
			return getPtrValue[*uint64](ptrToField)
		case reflect.Float32:
			return getPtrValue[*float32](ptrToField)
		case reflect.Float64:
			return getPtrValue[*float64](ptrToField)
		case reflect.Bool:
			return getPtrValue[*bool](ptrToField)
		case reflect.Struct:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		case reflect.Slice:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		case reflect.Array:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		default:
			panic("unhandled default case")
		}
	} else {
		switch kind {
		case reflect.String:
			return getPtrValue[string](ptrToField)
		case reflect.Int:
			return getPtrValue[int](ptrToField)
		case reflect.Int8:
			return getPtrValue[int8](ptrToField)
		case reflect.Int16:
			return getPtrValue[int16](ptrToField)
		case reflect.Int32:
			return getPtrValue[int32](ptrToField)
		case reflect.Int64:
			return getPtrValue[int64](ptrToField)
		case reflect.Uint:
			return getPtrValue[uint](ptrToField)
		case reflect.Uint8:
			return getPtrValue[uint8](ptrToField)
		case reflect.Uint16:
			return getPtrValue[uint16](ptrToField)
		case reflect.Uint32:
			return getPtrValue[uint32](ptrToField)
		case reflect.Uint64:
			return getPtrValue[uint64](ptrToField)
		case reflect.Float32:
			return getPtrValue[float32](ptrToField)
		case reflect.Float64:
			return getPtrValue[float64](ptrToField)
		case reflect.Bool:
			return getPtrValue[bool](ptrToField)
		case reflect.Struct:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		case reflect.Slice:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		case reflect.Array:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		default:
			panic("unhandled default case")
		}
	}
}

// getPtr returns a pointer to the field's value in the provided configuration object.
// It takes a parameter `conf` of type `any`, representing the configuration object.
// It returns an `unsafe.Pointer` to the `field's` value in the configuration object.
func (f Field) getPtr(obj interface{}) unsafe.Pointer {
	confPointer := ((*[2]unsafe.Pointer)(unsafe.Pointer(&obj)))[1]
	ptToField := unsafe.Add(confPointer, f.Offset)
	return ptToField
}

func setPtrValue[T any](ptr unsafe.Pointer, val any) {
	valSet := (*T)(ptr)
	*valSet = val.(T)
}

func getPtrValue[T any](ptr unsafe.Pointer) T {
	return *(*T)(ptr)
}

// Set updates the value of the fields in the provided object with the provided value.
// It takes two parameters:
//   - obj: interface{}, representing the object containing the field.
//   - val: interface{}, representing the new value for the field.
//
// The Set method uses the getPtr method to get a pointer to the fields in the object.
// It then performs a type switch on the kind of the fields to determine its type, and sets the value accordingly.
// The supported fields types are string, int, and bool.
// If the fields type is not one of the supported types, it panics with the message "unhandled default case".
func (f Field) Set(obj interface{}, val interface{}) {
	ptrToField := f.getPtr(obj)
	kind := f.Type.Kind()
	isPtr := false
	if kind == reflect.Ptr {
		isPtr = true
		kind = f.Type.Elem().Kind()
	}
	if isPtr {
		switch kind {
		case reflect.String:
			setPtrValue[*string](ptrToField, val)
		case reflect.Int:
			setPtrValue[*int](ptrToField, val)
		case reflect.Int8:
			setPtrValue[*int8](ptrToField, val)
		case reflect.Int16:
			setPtrValue[*int16](ptrToField, val)
		case reflect.Int32:
			setPtrValue[*int32](ptrToField, val)
		case reflect.Int64:
			setPtrValue[*int64](ptrToField, val)
		case reflect.Uint:
			setPtrValue[*uint](ptrToField, val)
		case reflect.Uint8:
			setPtrValue[*uint8](ptrToField, val)
		case reflect.Uint16:
			setPtrValue[*uint16](ptrToField, val)
		case reflect.Uint32:
			setPtrValue[*uint32](ptrToField, val)
		case reflect.Uint64:
			setPtrValue[*uint64](ptrToField, val)
		case reflect.Float32:
			setPtrValue[*float32](ptrToField, val)
		case reflect.Float64:
			setPtrValue[*float64](ptrToField, val)
		case reflect.Bool:
			setPtrValue[*bool](ptrToField, val)
		default:
			dest := reflect.NewAt(f.Type, ptrToField)
			dest = dest.Elem()
			source := reflect.ValueOf(val)
			dest.Set(source)
		}
	} else {
		switch kind {
		case reflect.String:
			setPtrValue[string](ptrToField, val)
		case reflect.Int:
			setPtrValue[int](ptrToField, val)
		case reflect.Int8:
			setPtrValue[int8](ptrToField, val)
		case reflect.Int16:
			setPtrValue[int16](ptrToField, val)
		case reflect.Int32:
			setPtrValue[int32](ptrToField, val)
		case reflect.Int64:
			setPtrValue[int64](ptrToField, val)
		case reflect.Uint:
			setPtrValue[uint](ptrToField, val)
		case reflect.Uint8:
			setPtrValue[uint8](ptrToField, val)
		case reflect.Uint16:
			setPtrValue[uint16](ptrToField, val)
		case reflect.Uint32:
			setPtrValue[uint32](ptrToField, val)
		case reflect.Uint64:
			setPtrValue[uint64](ptrToField, val)
		case reflect.Float32:
			setPtrValue[float32](ptrToField, val)
		case reflect.Float64:
			setPtrValue[float64](ptrToField, val)
		case reflect.Bool:
			setPtrValue[bool](ptrToField, val)
		default:
			dest := reflect.NewAt(f.Type, ptrToField)
			dest = dest.Elem()
			source := reflect.ValueOf(val)
			dest.Set(source)
		}
	}
}
