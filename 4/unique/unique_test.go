// 2014-11-10 Adam Bryt

package main

import (
	"testing"
)

func BenchmarkStringAssign(b *testing.B) {
	var a string = "abcdefghijkląąąąąęęęśśśśzzzzzzzzzzzzzzzzzzzz"
	var s string
	for i := 0; i < b.N; i++ {
		s = a
	}
	_ = s
}

func BenchmarkStringAssign_ConvertToBytes(b *testing.B) {
	var a string = "abcdefghijkląąąąąęęęśśśśzzzzzzzzzzzzzzzzzzzz"
	var s []byte
	for i := 0; i < b.N; i++ {
		s = []byte(a)
	}
	_ = s
}
