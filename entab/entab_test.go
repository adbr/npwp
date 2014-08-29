// 2014-06-27 Adam Bryt

package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"npwp/tab"
)

type entabTest struct {
	input   string
	output  string
	tablist []int
	indent  bool // indent only?
}

var entabTests = []entabTest{
	{
		// string pusty
		"",
		"",
		[]int{},
		false,
	},
	{
		// jeden znak
		"a",
		"a",
		[]int{},
		false,
	},
	{
		// kilka słów
		"aa ą bbbb",
		"aa ą bbbb",
		[]int{},
		false,
	},
	{
		// typowy, najprostszy przypadek
		//234567812345678123456781234567812345678
		"aa      b",
		"aa\tb",
		[]int{},
		false,
	},
	{
		//234567812345678123456781234567812345678
		"aa         b",
		"aa\t   b",
		[]int{},
		false,
	},
	{
		//234567812345678123456781234567812345678
		"aę         b    c",
		"aę\t   b\tc",
		[]int{},
		false,
	},
	{
		"aaa\n bbb\n",
		"aaa\n bbb\n",
		[]int{},
		false,
	},
	{
		// różne ilości spacji przed progiem
		"a       bb      ccc     dddd    eeeee   ffffff  ggggggg hhhhhhhh iiiiiiiii",
		"a\tbb\tccc\tdddd\teeeee\tffffff\tggggggg hhhhhhhh iiiiiiiii",
		[]int{},
		false,
	},
	{
		// kilka wierszy
		"a       bb      \nccc     dddd    \neeeee   \n",
		"a\tbb\t\nccc\tdddd\t\neeeee\t\n",
		[]int{},
		false,
	},
	{
		// jedna spacja przed progiem z nie-spacją - nie zastępuj przez tab
		"aaaaaaa b",
		"aaaaaaa b",
		[]int{},
		false,
	},
	{
		// jedna spacja przed progiem ze spacją - zastąp przez tab
		"aaaaaaa   b",
		"aaaaaaa\t  b",
		[]int{},
		false,
	},
	{
		// czy \n resetuje licznik spacji?
		"a       b  \n" +
			"        c\n",
		"a\tb  \n" +
			"\tc\n",
		[]int{},
		false,
	},
	{
		// nie standardowy rozmiar tabulacji (5 zamiast 8)
		"a    bb   ccc  dddd eeeee fffff",
		"a\tbb\tccc\tdddd eeeee fffff",
		[]int{5},
		false,
	},
	{
		// lista progów tabulacji
		"a  bb      c          d e f",
		"a\tbb\t   c          d e f",
		[]int{3, 8},
		false,
	},
	{
		// dane kończą się spacjami
		"a       b    ",
		"a\tb    ",
		[]int{},
		false,
	},
	{
		// dane kończą się spacjami - eof na progu
		"a       b       ",
		"a\tb       ",
		[]int{8},
		false,
	},
	{
		// zastąp tylko wcięcia - bez wcięć
		"a       b    ",
		"a       b    ",
		[]int{},
		true,
	},
	{
		// zastąp tylko wcięcia - jeden wiersz
		"        a        b             c",
		"\ta        b             c",
		[]int{},
		true,
	},
	{
		// zastąp tylko wcięcia
		"aaa\n" +
			"        a       b       c\n" +
			"                a       b\n" +
			"        aaa               b",
		"aaa\n" +
			"\ta       b       c\n" +
			"\t\ta       b\n" +
			"\taaa               b",
		[]int{},
		true,
	},
}

func TestEntab(t *testing.T) {
	for i, tc := range entabTests {
		tab, err := tab.NewTabulator(tc.tablist)
		if err != nil {
			t.Error(err)
		}

		ibuf := bytes.NewBufferString(tc.input)
		obuf := bytes.NewBuffer([]byte{})

		err = entab(obuf, ibuf, tab, tc.indent)
		if err != nil {
			t.Error(err)
		}

		if obuf.String() != tc.output {
			t.Errorf("tc #%d: oczekiwano: %q jest: %q", i, tc.output, obuf.String())
		}
	}
}

type fileTest struct {
	ifname  string
	efname  string
	tablist []int
	indent  bool // indent only?
}

var fileTests = []fileTest{
	{
		"testfiles/input.txt",
		"testfiles/expected.txt",
		[]int{},
		false,
	},
	{
		"testfiles/input.txt",
		"testfiles/expected-i.txt",
		[]int{},
		true,
	},
}

func TestEntabFiles(t *testing.T) {
	for _, tc := range fileTests {
		data, err := ioutil.ReadFile(tc.ifname)
		if err != nil {
			t.Error(err)
		}

		ibuf := bytes.NewBuffer(data)
		obuf := bytes.NewBuffer([]byte{})
		tab, err := tab.NewTabulator(tc.tablist)
		if err != nil {
			t.Error(err)
		}

		err = entab(obuf, ibuf, tab, tc.indent)
		if err != nil {
			t.Error(err)
		}

		expected, err := ioutil.ReadFile(tc.efname)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(expected, obuf.Bytes()) {
			t.Log(string(expected))
			t.Log(string(obuf.Bytes()))
			t.Errorf("dane wynikowe (entab z pliku %q) i oczekiwane (z pliku %q) są różne",
				tc.ifname, tc.efname)
		}
	}
}
