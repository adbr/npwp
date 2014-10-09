// 2014-10-08 Adam Bryt

package bubble

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
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
	}

	for i, tc := range tests {
		sort(tc.a)
		if !reflect.DeepEqual(tc.a, tc.b) {
			t.Errorf("tc %d: oczekiwano: %v jest %v", i, tc.b, tc.a)
		}
	}
}
