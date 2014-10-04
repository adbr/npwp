// 2014-10-03 Adam Bryt

// Pakiet header implementuje header plików w archiwum.
package header

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	hdrmark = "-h-" // string oznaczający początek nagłówka
)

// Header zawiera informacje o pliku zawarte w jego nagłówku.
type Header struct {
	Mark string // string identyfikujący początek nagłówka
	Name string // nazwa pliku
	Size int64  // rozmiar pliku w bajtach
}

// New tworzy i zwraca header z informacjami o pliku fname.
func New(fname string) (*Header, error) {
	fi, err := os.Stat(fname)
	if err != nil {
		return nil, err
	}

	h := &Header{
		hdrmark,
		fname,
		fi.Size(),
	}
	return h, nil
}

// Parse parsuje header w postaci stringu s (jak w archiwum) i zwraca
// obiekt Header.
func Parse(s string) (*Header, error) {
	a := strings.Fields(s)
	if len(a) != 3 {
		return nil, errors.New("header: wrong number of fields")
	}

	if a[0] != hdrmark {
		return nil, errors.New("header: wrong magic string")
	}

	size, err := strconv.ParseInt(a[2], 10, 64)
	if err != nil {
		return nil, errors.New("header: wrong file size: " + err.Error())
	}
	if size < 0 {
		return nil, fmt.Errorf("header: file size <0: %d", size)
	}

	h := &Header{
		a[0],
		a[1],
		size,
	}
	return h, nil
}

// String zwraca header w postaci takiej jak w archiwum - czyli string:
// "-h- name size\n", gdzie name jest nazwą pliku a size jest rozmiarem pliku
// w bajtach.
func (h *Header) String() string {
	return fmt.Sprintf("%s %s %d\n", h.Mark, h.Name, h.Size)
}

// Read czyta header z r.
func Read(r *bufio.Reader) (*Header, error) {
	line, err := r.ReadString('\n')
	if err == io.EOF {
		if len(line) > 0 {
			return nil, errors.New("header: eof before newline")
		}
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}

	hdr, err := Parse(line)
	if err != nil {
		return nil, err
	}
	return hdr, nil
}

// Write zapisuje header do w.
func Write(w *bufio.Writer, h *Header) error {
	_, err := w.WriteString(h.String())
	return err
}
