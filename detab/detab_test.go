// 2014-06-21 Adam Bryt

package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"npwp/tab"
)

type detabTest struct {
	input   string
	output  string
	tablist []int
	indent  bool // indent only?
	pad     rune
}

var detabTests = []detabTest{
	{
		// typowy, prosty przypadek
		"a\tb",
		"a.......b",
		[]int{},
		false,
		'.',
	},
	{
		// pusty string
		"",
		"",
		[]int{},
		false,
		'.',
	},
	{
		// jeden znak
		"a",
		"a",
		[]int{},
		false,
		'.',
	},
	{
		// kilka słów z polskimi znakami
		"abc dddd ąąą ęę ooo",
		"abc dddd ąąą ęę ooo",
		[]int{},
		false,
		'.',
	},
	{
		// jeden znak tab
		"\t",
		"........",
		[]int{},
		false,
		'.',
	},
	{
		// jeden znak tab w jednym wierszu
		"\t\n",
		"........\n",
		[]int{},
		false,
		'.',
	},
	{
		// jeden pusty wiersz
		"\n",
		"\n",
		[]int{},
		false,
		'.',
	},
	{
		// dwa wiersze bez znaków tab
		"aa ąąą\nbbb\n",
		"aa ąąą\nbbb\n",
		[]int{},
		false,
		'.',
	},
	{
		// kilka znaków tab rozdzielonych innymi znakami
		"a\tbb\tccc\tdddd\teeeee\tffffff\tggggggg\thhhhhhhh\ta",
		"a.......bb......ccc.....dddd....eeeee...ffffff..ggggggg.hhhhhhhh........a",
		[]int{},
		false,
		'.',
	},
	{
		// znaki tab i polskie znaki
		"ąęś\taaa\tbb",
		"ąęś.....aaa.....bb",
		[]int{},
		false,
		'.',
	},
	{
		// sąsiadujące znaki tab
		"a\t\tb\tc",
		"a...............b.......c",
		[]int{},
		false,
		'.',
	},
	{
		// znaki tab i spacje obok siebie
		"\t  \t",
		"........  ......",
		[]int{},
		false,
		'.',
	},
	{
		// kilka wierszy (zerowanie licznika kolumn)
		"a\tb\tccc\na\taa\tbbb\n",
		"a.......b.......ccc\na.......aa......bbb\n",
		[]int{},
		false,
		'.',
	},
	{
		// zamień tylko wcięcia
		"\tb\tccc\na\taa\tbbb\n\t\ta\n",
		"........b\tccc\na\taa\tbbb\n................a\n",
		[]int{},
		true,
		'.',
	},
	{
		// inny rozmiar tabulacji: 5
		"a\tąą\t\tbbb\n\tb\tcc\tddd\teeee\tfffff\tgggggg\th",
		"a....ąą........bbb\n.....b....cc...ddd..eeee.fffff.....gggggg....h",
		[]int{5},
		false,
		'.',
	},
	{
		// lista progów tabulacji: 3, 6, 8
		"a\tbb\tc\tdd\te\t\tf\n\ta\tb\tc\td",
		"a..bb.c.dd.e..f\n...a..b.c.d",
		[]int{3, 6, 8},
		false,
		'.',
	},
}

func TestDetab(t *testing.T) {
	for i, tc := range detabTests {
		tab, err := tab.NewTabulator(tc.tablist)
		if err != nil {
			t.Error(err)
		}
		ibuf := bytes.NewBufferString(tc.input)
		obuf := bytes.NewBuffer([]byte{})

		err = detab(obuf, ibuf, tab, tc.indent, tc.pad)
		if err != nil {
			t.Error(err)
		}
		if obuf.String() != tc.output {
			t.Errorf("tc #%d: oczekiwano: %q jest: %q\n", i, tc.output, obuf.String())
		}
	}
}

type fileTest struct {
	ifname  string
	efname  string
	tablist []int
}

var fileTests = []fileTest{
	{
		"testfiles/input.txt",
		"testfiles/expected.txt",
		[]int{},
	},
	{
		"testfiles/input.txt",
		"testfiles/expected-5.txt",
		[]int{5},
	},
	{
		"testfiles/input.txt",
		"testfiles/expected-3,5,9.txt",
		[]int{3, 5, 9},
	},
}

func TestDetabFiles(t *testing.T) {
	for _, tc := range fileTests {
		tab, err := tab.NewTabulator(tc.tablist)
		if err != nil {
			t.Error(err)
		}

		input, err := ioutil.ReadFile(tc.ifname)
		if err != nil {
			t.Error(err)
		}

		ibuf := bytes.NewBuffer(input)
		obuf := bytes.NewBuffer([]byte{})

		const space = ' '
		err = detab(obuf, ibuf, tab, false, space)
		if err != nil {
			t.Error(err)
		}

		expected, err := ioutil.ReadFile(tc.efname)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(expected, obuf.Bytes()) {
			t.Errorf("dane wynikowe (detab z pliku %q) i oczekiwane (z pliku %q) są różne",
				tc.ifname, tc.efname)
		}
	}
}
