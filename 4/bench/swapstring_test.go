// 2014-11-06 Adam Bryt
// Porównanie szybkości zamiany strinów w tablicy z zamianą struktury.
//
// % go test -bench . swap*
// testing: warning: no tests to run
// PASS
// BenchmarkStringSwap     20000000               147 ns/op
// BenchmarkStringSwapTmp  10000000               146 ns/op
// BenchmarkStringSwapCopy 10000000               164 ns/op
// BenchmarkStructSwap     20000000               148 ns/op
// ok      command-line-arguments  9.668s
// 
// czyli zamiana stringów w tablicy trwa tyle samo co zamiana struktury
// o dwóch polach typu int. Są zamieniane tylko "slice headers" stringów.

package main

import (
	"math/rand"
	"testing"
)

func BenchmarkStringSwap(b *testing.B) {
	var ss = []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww",
		"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq",
		"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
	}
	x := len(ss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := rand.Intn(x)
		k := rand.Intn(x)
		ss[j], ss[k] = ss[k], ss[j]
	}
}

func BenchmarkStringSwapTmp(b *testing.B) {
	var ss = []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww",
		"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq",
		"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
	}
	x := len(ss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := rand.Intn(x)
		k := rand.Intn(x)
		t := ss[j]
		ss[j] = ss[k]
		ss[k] = t
	}
}

func BenchmarkStringSwapCopy(b *testing.B) {
	var ss = []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww",
		"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq",
		"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
	}
	x := len(ss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := rand.Intn(x)
		k := rand.Intn(x)
		var t []byte
		copy(t, ss[j])
		ss[j] = ss[k]
		ss[k] = string(t)
	}
}

func BenchmarkStructSwap(b *testing.B) {
	type st struct {
		a int
		b int
	}
	var ss = []st{
		st{123, 123},
		st{22, 11},
		st{3444, 12},
		st{-123, 123},
		st{0, 11},
	}
	x := len(ss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := rand.Intn(x)
		k := rand.Intn(x)
		ss[j], ss[k] = ss[k], ss[j]
	}
}
