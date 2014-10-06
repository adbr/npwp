// 2014-08-22 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 2.5 "Parametry dyrektyw",
// program echo.
//
// NAZWA
//
// echo - drukuje argumenty polecenia
//
// SPOSÓB UŻYCIA
//
// echo [arg...]
//
// OPIS
//
// Program echo drukuje na stdout argumenty otrzymane podczas
// wywołania. Argumenty są oddzielane pojedynczą spacją, a na końcu
// jest drukowany znak newline ('\n'). Jeśli nie podano argumentów to
// jest drukowany tylko znak newline.
//
// PRZYKŁADY
//
// Polecenie:
//
//	$ echo Hello world!
//
// wypisuje na stdout wiersz:
//
// 	Hello world!
//
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
