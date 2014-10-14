// 2014-10-11 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 4.2 "Sortowanie metodą Shella".
package shell

func sort(a []int) {
	for gap := len(a) / 2; gap > 0; gap /= 2 {
		for i := gap; i < len(a); i++ {
			for j := i - gap; j >= 0; j -= gap {
				jj := j + gap
				if a[j] <= a[jj] {
					break
				} else {
					a[j], a[jj] = a[jj], a[j]
				}
			}
		}
	}
}
