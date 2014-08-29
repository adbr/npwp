// 2014-06-26 Adam Bryt

// Pakiet tab zawiera typ Tabulator i funkcje usługowe. Tabulator
// zawiera informacje i pozycjach tabulatora. Pakiet jest używany przez
// programy detab i entab.
package tab

import (
	"errors"
	"strconv"
	"strings"
)

// Rozmiar standardowego, domyślnego tabulatora.
const defaultTabSize = 8

// Tabulator zawiera informacje o pozycjach tabulatora, czyli listę
// kolumn, w których są ustawione progi tabulacji. Jeśli zawiera tylko
// jedną wartość, to jest ona traktowana jako odstęp między punktami
// tabulacji.
type Tabulator []int

// NewTabulator tworzy i inicjuje nowy Tabulator o progach tabulacji
// ustawionych w kolumnach cols. Jeśli lista pozycji tabulatora cols
// jest pusta, to progi tabulatora ustawia się w co 8 kolumnie. Jeśli
// lista pozycji tabulatora zawiera jeden element, to traktuje się go
// jako odstęp między progami tabulatora (np gdy cols[0] == 4 to progi
// tabulatora są ustawione w kolumnach 4, 8, 12, ...). Jeśli lista
// pozycji tabulatora cols zawiera więcej niż jeden element, to są to
// numery kolumn, w których są ustawione progi tabulacji.
// Przeprowadzana jest walidacja pozycji tabulatora - muszą być nie
// ujemne i posortowane rosnąco - jeśli nie są to jest zwracany error.
func NewTabulator(cols []int) (Tabulator, error) {
	if len(cols) == 0 {
		return Tabulator{defaultTabSize}, nil
	}

	err := validateTablist(cols)
	if err != nil {
		return Tabulator{}, err
	}
	return Tabulator(cols), nil
}

// IsTab zwraca true jeśli w kolumnie col jest ustawiony próg
// tabulatora.
func (t Tabulator) IsTab(col int) bool {
	switch {
	case len(t) == 1:
		return (col % t[0]) == 0
	case len(t) > 1:
		for _, c := range t {
			if c == col {
				return true
			}
		}
		return false
	default:
		panic("tabulator nie został ustawiony")
	}
}

// BeyondLastTab zwraca true jeśli kolumna col znajduje się poza
// ostatnim progiem tabulacji.
func (t Tabulator) BeyondLastTab(col int) bool {
	if len(t) == 1 {
		return false
	}

	if col > t[len(t)-1] {
		return true
	} else {
		return false
	}
}

// validateTablist sprawdza czy elementy listy pozycji tabulatora są nie
// ujemne i czy są posortowane rosnąco. Zwraca error != nil gdy te
// warunki nie są spełnione.
func validateTablist(cols []int) error {
	// Czy pozycje tabulatora są nie ujemne.
	for _, c := range cols {
		if c < 0 {
			return errors.New("pozycja tabulatora nie może być ujemna")
		}
	}

	// Czy pozycje tabulatora są posortowane rosnąco.
	for i := 0; i < len(cols)-1; i++ {
		if cols[i] > cols[i+1] {
			return errors.New("pozycje tabulatora muszą być posortowane rosnąco")
		}
	}

	return nil
}

// ParseTablist parsuje string s zawierający listę pozycji tabulacji
// postaci "a,b,c" (lista liczb całkowitych oddzielonych przecinkami) i
// zwraca tablicę liczb całkowitych.
func ParseTablist(s string) ([]int, error) {
	ss := strings.Split(s, ",")
	lst := []int{}
	for _, a := range ss {
		i, err := strconv.Atoi(a)
		if err != nil {
			return lst, err
		}
		lst = append(lst, i)
	}
	return lst, nil
}
