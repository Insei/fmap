package fmap

import "reflect"

type Storage interface {
	// Find returns the Field object and a boolean value indicating if the field with the given path was found.
	// The path parameter represents the path of the field in the struct.
	// If the field is found, the method returns the Field object and true.
	// If the field is not found, the method returns a nil Field object and false.
	Find(path string) (Field, bool)

	// MustFind returns the Field object for the field with the given path in the struct.
	// If the Field is not found, MustFind panics.
	MustFind(path string) Field

	// GetAllPaths returns a slice containing all paths of fields in the struct ordered like field struct definition.
	GetAllPaths() []string

	GetFieldByPtr(structPtr, fieldPtr any) (Field, error)
}

type Field interface {
	// GetName returns the name of the field.
	GetName() string

	// GetPkgPath returns the package import path of the Field struct type.
	GetPkgPath() string

	// GetType returns the reflect.Type of the field.
	GetType() reflect.Type

	// GetTag returns the reflect.StructTag of the field. The reflect.StructTag is a string.
	GetTag() reflect.StructTag

	// GetOffset returns the offset of the field in memory relative to the start of the struct.
	GetOffset() uintptr

	// GetIndex returns the index of the field within its containing struct as a slice of integers.
	GetIndex() []int

	// GetAnonymous returns a boolean value indicating whether the field is anonymous.
	GetAnonymous() bool

	// IsExported checks if a field is exported by checking its PkgPath property.
	IsExported() bool

	// Get returns the value of the storage in the provided object.
	// It takes a parameter `obj` of type `interface{}`, representing the object.
	// It returns the value of the storage as an `interface{}`.
	Get(obj any) any

	// GetPtr returns the pointer to the field's value in the provided object.
	// It takes a parameter `obj` of type `any`, representing the pointer to object.
	// It returns the pointer to the field's value as an `any`.
	GetPtr(obj any) any

	// Set updates the value of the storage in the provided object with the provided value.
	// It takes two parameters:
	//   - obj: interface{}, representing the object pointer containing the field.
	//   - val: interface{}, representing the new value for the field.
	Set(obj any, val any)

	// GetStructPath returns the struct path of the field.
	// It returns the struct path as a string.
	GetStructPath() string

	// GetTagPath returns the path of the field's tag value with the given tag name.
	// It takes two parameters:
	//   - tag: string, representing the tag name.
	//   - ignoreParentTagMissing: bool, representing whether to ignore the missing parent tags or not.
	// It returns the tag value path as a string.
	GetTagPath(tag string, ignoreParentTagMissing bool) string

	// GetParent returns the parent field of the current field, if not exist return nil.
	GetParent() Field

	// GetDereferencedType returns the dereferenced type of the field.
	// It returns the dereferenced type as a reflect.Type.
	GetDereferencedType() reflect.Type

	// GetDereferenced - uses reflect package for casting field value from obj to direct field value, i.e. dereferenced value.
	GetDereferenced(obj any) (any, bool)
}
