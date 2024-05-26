[![codecov](https://codecov.io/github/Insei/fmap/branch/main/graph/badge.svg?token=S8EDJENDSI)](https://codecov.io/github/Insei/fmap) 
[![build](https://github.com/insei/fmap/actions/workflows/go.yml/badge.svg)](https://github.com/Insei/fmap/actions/workflows/go.yml)
[![Goreport](https://goreportcard.com/badge/github.com/insei/fmap)](https://goreportcard.com/report/github.com/insei/fmap)
[![GoDoc](https://godoc.org/github.com/insei/fmap?status.svg)](https://godoc.org/github.com/insei/fmap)
# FMap 
FMap is a simple library for working with structs as map of fields. Switch case and reflect based.

# Description
FMap creates new map with filed names as key and with fmap.Field(reflect.StructField) values. This is unsafe library, be careful while use.

fmap.Field has 3 advanced methods:

```
Get[T any]() any
Set(obj any, val any)
GetPtr(obj any) any
```

where the `obj` should be not nil pointer to struct.<br>
`Get[T any]() any` -  return expected typed value as interface{}.<br>
`Set(obj any, val any)` -  `val` expected typed value to set, {}interface can be used.<br>
`GetPtr(obj any) any` - return expected typed value pointer to struct field as interface{}.
# Example

```go
package main

import (
	"time"
	"fmt"

	"github.com/insei/fmap"
)

type City struct {
	Name string
}

type People struct {
	Name     string
	Age      uint8
	Birthday time.Time
	City City
}

func main() {
	p := &People{}
	fields := fmap.Get[People]() // or fmap.GetFrom(p)
	fields["Name"].Set(p, "Test")
	fields["Age"].Set(p, uint8(5))
	fields["Birthday"].Set(p, time.Now())
	fields["City.Name"].Set(p, "DefaultCity")
	fmt.Print(*p)
}
```

More examples in `fields_test.go`, like slice fields, nested structs, pointers etc.
