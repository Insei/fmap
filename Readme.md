[![codecov](https://codecov.io/github/Insei/fmap/branch/main/graph/badge.svg?token=S8EDJENDSI)](https://codecov.io/github/Insei/fmap) 
[![build](https://github.com/insei/fmap/actions/workflows/go.yml/badge.svg)](https://github.com/Insei/fmap/actions/workflows/go.yml)
[![Goreport](https://goreportcard.com/badge/github.com/insei/fmap)](https://goreportcard.com/report/github.com/insei/fmap)
[![GoDoc](https://godoc.org/github.com/insei/fmap?status.svg)](https://godoc.org/github.com/insei/fmap)
# FMap 
FMap is a simple library for working with structs as a storage of fields. Switch case and reflect based.
This is unsafe library, be careful while use.

# Installation
Install via go get. Note that Go 1.18 or newer is required.
```sh
go get github.com/insei/fmap/v3@latest
```

# Description
`fmap.GetFrom(obj any)` and `fmap.Get[T any]()` creates new fmap.Storage. This storage manage access to fmap.Field by field path like in struct.

```go
type Storage interface {
    // Find returns the Field object and a boolean value indicating if the field with the given path was found.
    // The path parameter represents the path of the field in the struct.
    // If the field is found, the method returns the Field object and true.
    // If the field is not found, the method returns a nil Field object and false.
    Find(path string) (Field, bool)
    
    // MustFind returns the Field object for the field with the given path in the struct.
    // If the Field is not found, MustFind panics.
    MustFind(path string) Field
    
    // GetAllPaths returns a slice containing all paths of fields in the struct.
    GetAllPaths() []string
}
```

fmap.Field is an advanced abstraction level for reflect.StructField with some advanced methods:
```go
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
    
    // Extended Methods
    
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
```
# Example

```go
package main

import (
	"time"
	"fmt"

	"github.com/insei/fmap/v3"
)

type City struct {
	Name string `json:"name"`
}

type People struct {
	Name     string
	Age      uint8
	Birthday time.Time
	City City `json:"city"`
}

func main() {
	p := &People{}
	fields := fmap.Get[People]() // or fmap.GetFrom(p)
	fields.MustFind("Name").Set(p, "Test")
	fields.MustFind("Age").Set(p, uint8(5))
	fields.MustFind("Birthday").Set(p, time.Now())
	fields.MustFind("City.Name").Set(p, "DefaultCity")
	jsonPath := fields.MustFind("City.Name").GetTagPath("json", false) // city.name
	cityField := fields.MustFind("City.Name").GetParent()
	cityStruct := cityField.Get(p)
	fmt.Print(*p, jsonPath, cityStruct)
}
```

More examples in `field_test.go`, like slice fields, nested structs, pointers etc.

# Benchmarks
`fmap.GetFrom(obj any) map[string]Field`
```
BenchmarkGetFrom
BenchmarkGetFrom-16     93002347                12.62 ns/op            0 B/op          0 allocs/op
```

`Field.Get(obj any) any`
```
BenchmarkFieldGet
BenchmarkFieldGet-16            88818492                14.05 ns/op            0 B/op          0 allocs/op
```

`Raw access to field from struct :)`
```
BenchmarkRawFieldGet
BenchmarkRawFieldGet-16         1000000000               0.2350 ns/op          0 B/op          0 allocs/op
```