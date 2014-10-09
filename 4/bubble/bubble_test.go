// 2014-10-08 Adam Bryt

package bubble

import (
	"testing"
	"reflect"
)

func TestSort(t *testing.T) {
	type test struct {
		a []int
		b []int
	}

	tests := []test{
		{
			[]int{2, 1, 3, 2},
			[]int{1, 2, 2, 3},
		},
	}

	for i, tc := range tests {
		sort(tc.a)
		if !reflect.DeepEqual(tc.a, tc.b) {
			t.Errorf("tc %d: oczekiwano: %v jest %v", i, tc.b, tc.a)
		}
	}
}
