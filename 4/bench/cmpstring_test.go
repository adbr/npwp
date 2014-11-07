// 2014-11-07 Adam Bryt
// Sprawdzenie jak szybko działa prównywanie stringów funkcją
// collate.CompareString() w porównaniu z operatorem <.
// 
// go test -bench . cmp*
// testing: warning: no tests to run
// PASS
// BenchmarkCompareOperatorL       200000000                9.33 ns/op
// BenchmarkCompareCollate  1000000              2340 ns/op
// ok      command-line-arguments  5.174s
//
// czyli collate.CompareString() jest 200 razy wolniejsze.

package main

import (
	"testing"

	"code.google.com/p/go.text/collate"
	"code.google.com/p/go.text/language"
)

func BenchmarkCompareOperatorL(b *testing.B) {
	s1 := "abc ąęś 123"
	s2 := "123 abc zzz"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s1 < s2
	}
}

func BenchmarkCompareCollate(b *testing.B) {
	s1 := "abc ąęś 123"
	s2 := "123 abc zzz"
	c := collate.New(language.Polish)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.CompareString(s1, s2)
	}
}
