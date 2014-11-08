// 2014-10-23 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 4.4 "Sortowanie szybkie".
package quick

func quick(a []int) {
	if len(a) < 2 {
		return
	}
	xi := len(a) - 1 // indeks pivota
	x := a[xi]       // wartość pivota
	i, j := 0, len(a)-1
	// zamieniaj miejscami elementy większe i mniejsze od pivota
	for {
		if i == j {
			break
		}
		for {
			if i == j || a[i] > x {
				break
			}
			i++
		}
		for {
			if j == i || a[j] < x {
				break
			}
			j--
		}
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	// wstaw pivot w miejsce spotkania i j
	a[i], a[xi] = a[xi], a[i]
	quick(a[:i])
	quick(a[i:])
}
