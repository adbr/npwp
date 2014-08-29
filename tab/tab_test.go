// 2014-06-27 Adam Bryt

package tab

import (
	"testing"
)

func TestNewTabulator(t *testing.T) {
	// Tabulator z domyślną wartością progów tabulacji
	tab, err := NewTabulator([]int{})
	if err != nil {
		t.Error(err)
	}
	if len(tab) != 1 {
		t.Error("domyślny tabulator powinien mieć jeden element:", tab)
	}
	if tab[0] != 8 {
		t.Error("domyślnym rozmiarem tabulacji powinno być 8:", tab[0])
	}

	// Tabulator z jedną wartością (rozmiar tabulacji)
	tab, err = NewTabulator([]int{4})
	if err != nil {
		t.Error(err)
	}
	if len(tab) != 1 {
		t.Error("tabulator powinien mieć jeden element:", tab)
	}
	if tab[0] != 4 {
		t.Error("rozmiarem tabulacji powinno być 4:", tab[0])
	}

	// Tabulator z kilkoma progami tabulacji
	tab, err = NewTabulator([]int{4, 8, 20})
	if err != nil {
		t.Error(err)
	}
	if len(tab) != 3 {
		t.Error("tabulator powinien mieć 3 elementy:", tab)
	}
	if tab[1] != 8 {
		t.Error("drugi próg tabulacji powinien mieć wartość 8:", tab[1])
	}

	// Zła wartość progu tabulatora - powinien być błąd
	tab, err = NewTabulator([]int{4, -8, 20})
	if err == nil {
		t.Error("konstruktor nie zgłosił błędu: ujemna wartość progu tabulacji")
	}

	// Progi tabulacji nie posortowane - powinien być błąd
	tab, err = NewTabulator([]int{8, 4, 20})
	if err == nil {
		t.Error("konstruktor nie zgłosił błędu: nie posortwane progi tabulacji")
	}
	tab, err = NewTabulator([]int{4, 4, 20})
	if err != nil {
		t.Error(err)
	}
}

// Czy w kolumnie col jest próg tabulacji
type colTab struct {
	col int  // nr kolumny
	tab bool // czy jest próg tabulacji
}

func TestIsTab(t *testing.T) {
	// Domyślny tabulator: 8
	tab, err := NewTabulator([]int{})
	if err != nil {
		t.Error(err)
	}

	tests := []colTab{
		{-2, false},
		{0, true},
		{1, false},
		{2, false},
		{3, false},
		{4, false},
		{5, false},
		{6, false},
		{7, false},
		{8, true},
		{9, false},
		{15, false},
		{16, true},
		{17, false},
	}
	for _, tc := range tests {
		v := tab.IsTab(tc.col)
		if v != tc.tab {
			t.Errorf("próg tabulacji w kolumnie %d: powinno być %v, jest %v",
				tc.col, tc.tab, v)
		}
	}

	// Tabulator o odstępie 3
	tab, err = NewTabulator([]int{3})
	if err != nil {
		t.Error(err)
	}

	tests = []colTab{
		{-2, false},
		{0, true},
		{1, false},
		{2, false},
		{3, true},
		{4, false},
		{5, false},
		{6, true},
		{7, false},
		{8, false},
		{9, true},
		{10, false},
		{12, true},
		{17, false},
	}
	for _, tc := range tests {
		v := tab.IsTab(tc.col)
		if v != tc.tab {
			t.Errorf("próg tabulacji w kolumnie %d: powinno być %v, jest %v",
				tc.col, tc.tab, v)
		}
	}

	// Tabulator z listą pozycji tabulatora
	tab, err = NewTabulator([]int{4, 6, 8})
	if err != nil {
		t.Error(err)
	}

	tests = []colTab{
		{-2, false},
		{0, false},
		{1, false},
		{2, false},
		{3, false},
		{4, true},
		{5, false},
		{6, true},
		{7, false},
		{8, true},
		{9, false},
		{10, false},
		{12, false},
		{99, false},
	}
	for _, tc := range tests {
		v := tab.IsTab(tc.col)
		if v != tc.tab {
			t.Errorf("próg tabulacji w kolumnie %d: powinno być %v, jest %v",
				tc.col, tc.tab, v)
		}
	}
}

type beyondTest struct {
	tablist []int
	col     int  // nr kolumny
	exp     bool // czy kolumna col jest poza ostatnim progiem
}

var beyondTests = []beyondTest{
	{
		[]int{},
		0, false,
	},
	{
		[]int{},
		1, false,
	},
	{
		[]int{},
		8, false,
	},
	{
		[]int{},
		9, false,
	},
	{
		[]int{},
		12, false,
	},
	{
		[]int{4},
		2, false,
	},
	{
		[]int{4},
		5, false,
	},
	{
		[]int{4},
		22, false,
	},
	{
		[]int{3, 5, 8},
		1, false,
	},
	{
		[]int{3, 5, 8},
		8, false,
	},
	{
		[]int{3, 5, 8},
		9, true,
	},
	{
		[]int{3, 5, 8},
		22, true,
	},
}

func TestBeyond(t *testing.T) {
	for i, tc := range beyondTests {
		tab, err := NewTabulator(tc.tablist)
		if err != nil {
			t.Error(err)
		}
		v := tab.BeyondLastTab(tc.col)
		if v != tc.exp {
			t.Errorf("tc: #%d: tab: %v, col: %d: jest %v, oczekiwano %v",
				i, tab, tc.col, v, tc.exp)
		}
	}
}

type tablist struct {
	input  string // argument parseTablist
	output []int  // oczekiwany rezultat parseTablist
	error  bool   // czy funkcja powinna zwrócić błąd
}

var tablists = []tablist{
	{"", []int{}, true}, // musi być co najmniej jeden element
	{"1", []int{1}, false},
	{"1,2", []int{1, 2}, false},
	{"1,2,4", []int{1, 2, 4}, false},
	{"0,-2,4", []int{0, -2, 4}, false}, // 0 i wartość ujemna
	{"1,2.2,4", []int{}, true},         // wartość zmiennoprzecinkowa
	{"1, 2, 3", []int{}, true},         // odstępy w liście
	{"4,3,-1", []int{4, 3, -1}, false}, // porządek nie ma znaczenia
	{"1,w,2", []int{}, true},           // znak nie będący cyfrą
}

func TestParseTablist(t *testing.T) {
	for _, tc := range tablists {
		v, err := ParseTablist(tc.input)
		if tc.error == false && err != nil { // nieoczekiwany błąd
			t.Error(err)
			continue
		}
		if tc.error == true { // powinien wystąpić błąd
			if err != nil {
				continue // nie ma sensu dalsze sprawdzanie
			} else {
				t.Error("nie wystąpił oczekiwany błąd dla:", tc)
			}
		}

		if len(v) != len(tc.output) {
			t.Errorf("zła długość tablicy: jest %v, oczekiwano %v",
				v, tc.output)
		}

		// porównaj poszczególne elementy tablicy
		for i, e := range v {
			if e != tc.output[i] {
				t.Errorf("elementy o indeksie %d są różne: jest %v, oczekiwano %v",
					i, v, tc.output)
			}
		}
	}
}
