// 2014-08-29 Adam Bryt

// Pakiet util zawiera różne funkcje usługowe.
package util

// itoc zwraca stringową reprezentację liczby całkowitej n.
func itoc(n int) string {
	s := ""

	if n < 0 {
		s = "-"
		n = -n
	}

	if n > 10 {
		s += itoc(n / 10)
	}

	c := (n % 10) + '0'
	s += string(c)

	return s
}

// ctoi zwraca liczbę int, której stringowa reprezentacja znajduje sie w s.
func ctoi(s string) int {
	i := 0

	// pomiń spacje i tabulacje
	for (i < len(s)) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}

	sign := 1
	if i < len(s) && s[i] == '-' {
		sign = -1
	}

	// pomiń znak + lub -
	if i< len(s) && (s[i] == '-' || s[i] == '+') {
		i++
	}

	n := 0
	for i < len(s) && isdigit(s[i]) {
		n = n*10 + int(s[i]-'0')
		i++
	}

	return n * sign
}

// isdigit zwraca true jeśli znak c jest cyfrą.
func isdigit(c byte) bool {
	if '0' <= c && c <= '9' {
		return true
	} else {
		return false
	}
}
