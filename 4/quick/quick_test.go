// 2014-10-23 Adam Bryt

package quick

import (
	"reflect"
	"testing"
)

func TestQuick(t *testing.T) {
	type test struct {
		a []int
		b []int
	}

	tests := []test{
		{
			[]int{},
			[]int{},
		},
		{
			[]int{1},
			[]int{1},
		},
		{
			[]int{2, 1},
			[]int{1, 2},
		},
		{
			[]int{4, -2, 0},
			[]int{-2, 0, 4},
		},
		{
			[]int{2, 1, 3, 2},
			[]int{1, 2, 2, 3},
		},
		{
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
		},
		{
			[]int{2, 3, 6, 0, 1, 3, 3, 5, 6, 7, -2},
			[]int{-2, 0, 1, 2, 3, 3, 3, 5, 6, 6, 7},
		},
	}

	for i, tc := range tests {
		quick(tc.a)
		if !reflect.DeepEqual(tc.a, tc.b) {
			t.Errorf("tc %d: oczekiwano: %v jest %v", i, tc.b, tc.a)
		}
	}
}
