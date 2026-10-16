package fmap

import (
	"reflect"
	"strings"
	"unsafe"
)

type field struct {
	reflect.StructField
	structPath      string
	parent          *field
	dereferenceType reflect.Type
	// owner - the nearest field with the pointer to the struct, where the field is stored.
	// Offset of the field is calculated from the beginning of that struct, nil if there are no pointers in the path.
	owner *field
}

func newField(structField reflect.StructField, structPath string, parent *field) *field {
	return &field{
		StructField:     structField,
		structPath:      structPath,
		parent:          parent,
		dereferenceType: dereferenceType(structField.Type),
	}
}

// isStructPtr checks that the type is a pointer to the struct, fields of such struct are available by the path.
func isStructPtr(typeOf reflect.Type) bool {
	return typeOf.Kind() == reflect.Ptr && typeOf.Elem().Kind() == reflect.Struct
}

func dereferenceType(typeOf reflect.Type) reflect.Type {
	for typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	return typeOf
}

func (f *field) GetName() string {
	return f.Name
}

func (f *field) GetPkgPath() string {
	return f.PkgPath
}

func (f *field) GetType() reflect.Type {
	return f.Type
}

func (f *field) GetTag() reflect.StructTag {
	return f.Tag
}

func (f *field) GetOffset() uintptr {
	return f.Offset
}

func (f *field) GetIndex() []int {
	return f.Index
}

func (f *field) GetAnonymous() bool {
	return f.Anonymous
}

func (f *field) IsExported() bool {
	return f.PkgPath == ""
}

func (f *field) GetStructPath() string {
	return f.structPath
}

func (f *field) GetParent() Field {
	if f.parent == nil {
		// nil pointer in the interface is not equal to nil
		return nil
	}
	return f.parent
}

func (f *field) GetPtr(obj interface{}) interface{} {
	ptrToField := f.getPtr(obj)
	if ptrToField == nil {
		// the struct with the field is not created, typed nil pointer is returned
		return reflect.Zero(reflect.PtrTo(f.Type)).Interface()
	}
	return reflect.NewAt(f.Type, ptrToField).Interface()
}

// Get returns the value of the storage in the provided object.
// It takes a parameter `obj` of type `interface{}`, representing the object.
// It returns the value of the storage as an `interface{}`.
func (f *field) Get(obj interface{}) interface{} {
	ptrToField := f.getPtr(obj)
	if ptrToField == nil {
		// the struct with the field is not created
		return nil
	}
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
		default:

			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
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
		default:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		}
	}
}

// tagName returns the name from the tag value, the part before the comma: json:"name,omitempty" -> name.
func (f *field) tagName(tag string) string {
	return strings.Split(f.Tag.Get(tag), ",")[0]
}

func (f *field) GetTagPath(tag string, ignoreParentTagMissing bool) string {
	tagPath := f.tagName(tag)
	if tagPath == "" {
		return tagPath
	}
	// embedded structs without the name in the tag are transparent, their fields are on the level of the embedding struct,
	// like in encoding/json and the most of the encoding libraries; yaml:",inline" has no name too
	parent := f.parent
	for parent != nil && parent.Anonymous && parent.tagName(tag) == "" {
		parent = parent.parent
	}
	if parent == nil {
		return tagPath
	}
	parentTag := parent.GetTagPath(tag, ignoreParentTagMissing)
	if parentTag == "" && !ignoreParentTagMissing {
		return ""
	}
	if parentTag == "" {
		return tagPath
	}
	return parentTag + "." + tagPath
}

// getPtr returns a pointer to the field's value in the provided configuration object.
// It takes a parameter `conf` of type `any`, representing the configuration object.
// It returns an `unsafe.Pointer` to the `field's` value in the configuration object,
// nil if the field is stored in the struct by the pointer and that pointer is nil.
func (f *field) getPtr(obj interface{}) unsafe.Pointer {
	confPointer := ((*[2]unsafe.Pointer)(unsafe.Pointer(&obj)))[1]
	return f.addr(confPointer, false)
}

// addr returns a pointer to the field's value in the object by the root pointer.
// Nil pointers to the structs in the path are replaced with new structs if create is true, otherwise nil is returned.
func (f *field) addr(root unsafe.Pointer, create bool) unsafe.Pointer {
	if f.owner == nil {
		return unsafe.Add(root, f.Offset)
	}
	ownerPtr := f.owner.addr(root, create)
	if ownerPtr == nil {
		return nil
	}
	structPtr := *(*unsafe.Pointer)(ownerPtr)
	if structPtr == nil {
		if !create {
			return nil
		}
		structPtr = reflect.New(f.owner.Type.Elem()).UnsafePointer()
		*(*unsafe.Pointer)(ownerPtr) = structPtr
	}
	return unsafe.Add(structPtr, f.Offset)
}

func setPtrValue[T any](ptr unsafe.Pointer, val any) {
	valSet := (*T)(ptr)
	if typed, ok := val.(T); ok {
		*valSet = typed
		return
	}
	// value of the named type: time.Duration for int64 field kind, *time.Duration for *int64 and so on
	*valSet = reflect.ValueOf(val).Convert(reflect.TypeOf(*valSet)).Interface().(T)
}

func getPtrValue[T any](ptr unsafe.Pointer) T {
	return *(*T)(ptr)
}

// Set updates the value of the storage in the provided object with the provided value.
// It takes two parameters:
//   - obj: interface{}, representing the object containing the field.
//   - val: interface{}, representing the new value for the field.
//
// The Set method uses the getPtr method to get a pointer to the storage in the object.
// It then performs a type switch on the kind of the storage to determine its type, and sets the value accordingly.
// The supported storage types are string, int, and bool.
// If the storage type is not one of the supported types, it panics with the message "unhandled default case".
func (f *field) Set(obj interface{}, val interface{}) {
	// nil pointers to the structs in the path of the field are replaced with new structs
	ptrToField := f.addr(((*[2]unsafe.Pointer)(unsafe.Pointer(&obj)))[1], true)
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

func (f *field) GetDereferencedType() reflect.Type {
	// calculated once in newField, the field is read only after that and can be used from many goroutines
	return f.dereferenceType
}

// GetDereferenced - uses reflect package for casting field value from obj to direct field value, i.e. dereferenced value.
func (f *field) GetDereferenced(obj any) (any, bool) {
	derefType := f.GetDereferencedType()
	value := f.Get(obj)
	valOf := reflect.ValueOf(value)
	for valOf.IsValid() && valOf.Kind() == reflect.Ptr {
		valOf = valOf.Elem()
	}
	if valOf.IsValid() && derefType == valOf.Type() {
		return valOf.Interface(), true
	}
	return nil, false
}
