// 2014-08-29 Adam Bryt

package util

import (
	"testing"
)

func TestItoc(t *testing.T) {
	type test struct {
		n int
		s string
	}

	tests := []test{
		{
			3,
			"3",
		},
		{
			-5,
			"-5",
		},
		{
			+123,
			"123",
		},
		{
			-12345678,
			"-12345678",
		},
		{
			0,
			"0",
		},
	}

	for i, tc := range tests {
		v := itoc(tc.n)
		if v != tc.s {
			t.Errorf("tc: %d, oczekiwano: %q, jest: %q", i, tc.s, v)
		}
	}
}

func TestCtoi(t *testing.T) {
	type test struct {
		s string
		n int
	}

	tests := []test{
		{
			"5",
			5,
		},
		{
			"1234",
			1234,
		},
		{
			"  \t\t 123",
			123,
		},
		{
			" 123abc",
			123,
		},
		{
			"abc123",
			0,
		},
		{
			"abc",
			0,
		},
		{
			"-123",
			-123,
		},
		{
			"+123",
			123,
		},
	}

	for i, tc := range tests {
		v := ctoi(tc.s)
		if v != tc.n {
			t.Errorf("tc: %d, oczekiwano: %d, jest %d", i, tc.n, v)
		}
	}
}
