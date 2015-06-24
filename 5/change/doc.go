// 2015-05-12 Adam Bryt

/*
Narzędzia Programistyczne w Pascalu,
rozdział 5.5 "Zamiana tekstu",
program change.

NAZWA

change - zamienia wzorce w tekście

SPOSÓB UŻYCIA

change <pattern> [<substitution>]

OPIS

Program change czyta wiersze tekstu ze standardowego wejścia,
zamienia wszystkie nie nakładające się fragmenty pasujące do
<pattern> na <substitution> i drukuje je na standardowe wyjście.
Reguły budowy wzorców <pattern> są takie same jak w programie
'find'.  Jeśli argument <substitution> nie istnieje to tekst
pasujący do <pattern> jest usuwany. Jeśli argument <substitution>
zawiera znaki '&' to te znaki są zastępowane fragmentem pasującym
do <pattern>.  Żeby znak '&' pozbawić specjalnego znaczenia,
należy zamiast niego użyć sekwencji '@&'.

UWAGI

Jeśli ostatni wiersz wejściowy nie jest zakończony znakiem '\n',
to na wyjściu zostanie dodany do niego znak '\n'.

Ograniczenie na długość wiersza wejściowego wynosi około 64
KB.  Jeśli wiersz jest za długi to zostanie zgłoszony błąd
'token too long'.

*/
package main
