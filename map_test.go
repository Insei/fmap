package fmap

import (
	"testing"
)

func BenchmarkGetFrom(b *testing.B) {
	tt := TestStruct{}
	for i := 0; i < b.N; i++ {
		_ = GetFrom(&tt)
	}
}

func BenchmarkFieldGet(b *testing.B) {
	tt := TestStruct{}
	fields := GetFrom(&tt)

	for i := 0; i < b.N; i++ {
		_ = fields["String"].Get(&tt)
	}
}

func BenchmarkRawFieldGet(b *testing.B) {
	tt := TestStruct{}
	for i := 0; i < b.N; i++ {
		_ = tt.String
	}
}
