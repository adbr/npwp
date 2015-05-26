// Narzędzia Programistyczne w Pascalu,
// rozdział 5 "Wzorce tekstowe",
// program find.
//
// NAZWA
//
// find - wyszukuje wiersze tekstu pasujące do wzorca
//
// SPOSÓB UŻYCIA
//
// find wzorzec [<file] [>file2]
//
// OPIS
//
// Program find czyta wiersze tekstu z stdin i drukuje na stdout te
// wiersze, które zawierają fragment pasujący do wzorca.
//
// Wzorzec jest konkatenacją następujących elementów:
//
//	c       znak c (może być UTF8)
//	?       dowolny znak oprócz '\n'
//	%       początek wiersza
//	$       koniec wiersza (przed znakiem '\n')
//	[...]   klasa znaków (dowolny znak z wymienionych)
//	[^...]  dopełnienie klasy znaków (dowolny znak z wyjątkiem wymienionych)
//	*       domknięcie (zero lub więcej wystąpień poprzedniego elementu wzorca)
//	@c      wyróżnik (przywraca pierwotne znaczenie znaku c, np @%)
//
// Znaki ?%$[]*@ są metaznakami i mają specjalne znaczenie we wzorcu.
// To specjalne znaczenie zanika w następujących przypadkach:
//
//	po znaku @
//	wewnątrz [...] (z wyjątkiem @])
//	% nie na początku
//	$ nie na końcu
//	* na początku
//
// Klasa znaków zawiera zero lub więcej następujących elementów w nawiasach []:
//
//	c	dowolny znak (może być utf8)
//	c1-c2	przedział alfanumerycznych znaków ASCII (a-z lub A-Z lub 0-9)
//	^	jeśli występuje na początku, po znaku [, oznacza dopełnienie
//		(negację) klasy znaków (np. [^ab] oznacza dowolny znak różny od
//		a i b
//
// Znaki wyróżnione w opisie klasy znaków tracą swoje specjalne znaczenie gdy
// są poprzedzone znakiem @, oraz w przypadkach:
//
//	^	nie na początku
//	-	na początku lub na końcu
//
// Wyróżnik, czyli znak @, odbiera specjalne znaczenie metaznakom, a ponadto:
//
//	@@	oznacza znak @
//	@c	oznacza znak c (dla dowolnego znaku c)
//	@n	oznacza znak nowego wiersza '\n'
//	@t	oznacza znak tabulacji '\t'
//
// PRZYKŁADY
//
// Wydrukowanie wierszy zawierających komentarz zaczynający się od początku wiersza:
//
//	cat *.go | ./find "%//?*"
//
package main
