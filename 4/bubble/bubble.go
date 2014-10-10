// 2014-10-08 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 4.1 "Sortowanie bąbelkowe".
package bubble

// sort sortuje a metodą sortowania bąbelkowego (bubble sort).
func sort(a []int) {
	for i := len(a) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
			}
		}
	}
}
