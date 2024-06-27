[![codecov](https://codecov.io/github/Insei/fmap/branch/main/graph/badge.svg?token=S8EDJENDSI)](https://codecov.io/github/Insei/fmap) 
[![build](https://github.com/insei/fmap/actions/workflows/go.yml/badge.svg)](https://github.com/Insei/fmap/actions/workflows/go.yml)
[![Goreport](https://goreportcard.com/badge/github.com/insei/fmap)](https://goreportcard.com/report/github.com/insei/fmap)
[![GoDoc](https://godoc.org/github.com/insei/fmap?status.svg)](https://godoc.org/github.com/insei/fmap)
# FMap 
FMap is a simple library for working with structs as map of fields. Switch case and reflect based.

# Installation
Install via go get. Note that Go 1.18 or newer is required.
```sh
go get github.com/insei/fmap/v2@latest
```

# Description
FMap creates new map with filed names as key and with fmap.Field implementations. This is unsafe library, be careful while use.

fmap.Field has default and some advanced methods:
* Default (for interacting with reflect.StructField values)
```
GetName() string
GetPkgPath() string
GetType() reflect.Type
GetTag() reflect.StructTag
GetOffset() uintptr
GetIndex() []int
GetAnonymous() bool
IsExported() bool
```
* Advanced
```
GetStructPath() string
GetTagPath(tag string, ignoreParentTagMissing bool) string
GetParent() IField
Get(obj any) any
Set(obj any, val any)
GetPtr(obj any) any
```

where the `obj` should be not nil pointer to struct.<br>
`GetStructPath() string` - returns full path to the field into struct.<br>
`GetTagPath(tag string, ignoreParentTagMissing bool) string` - returns full path to the field into struct but as tag values.<br>
`GetParent() IField` - returns a parent field, if exist.<br>
`Get(obj any)` -  return expected typed value as interface{}.<br>
`Set(obj any, val any)` -  `val` expected typed value to set, {}interface can be used.<br>
`GetPtr(obj any) any` - return expected typed value pointer to struct field as interface{}.
# Example

```go
package main

import (
	"time"
	"fmt"

	"github.com/insei/fmap/v2"
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
	fields["Name"].Set(p, "Test")
	fields["Age"].Set(p, uint8(5))
	fields["Birthday"].Set(p, time.Now())
	fields["City.Name"].Set(p, "DefaultCity")
	jsonPath := fields["City.Name"].GetTagPath("json", false) // city.name
	cityField := fields["City.Name"].GetParent()
	cityStruct := cityField.Get(p)
	fmt.Print(*p, jsonPath, cityStruct)
}
```

More examples in `fields_test.go`, like slice fields, nested structs, pointers etc.
