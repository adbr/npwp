// 2015-06-24 Adam Bryt

/*
Narzędzia Programistyczne w Pascalu,
rozdział 6 "Redagowanie tekstów",
program edit.

NAZWA

edit - edytor plików tekstowych

SPOSÓB UŻYCIA

edit [plik]

OPIS

Program edit jest interakcyjnym edytorem tekstowym, który czyta
wiersze z wejścia i wypisuje odpowiednie informacje na wyjście, w
zależności od dyrektywy. Jego działanie polega na czytaniu na żądanie
plików tekstowych do wewnętrznego bufora (który może być duży);
zawartość bufora może być następnie wyświetlana i modyfikowana przy
użyciu różnych dyrektyw, wreszcie na żądanie zapisana w części lub
całości na pliki tekstowe. Bufor jest zorganizowany w postaci ciągu
wierszy, numerowanych od 1; wiersze są automatycznie przenumerowywane
w miarę przybywania lub ubywania tekstu. Wyszukiwań i wymian
kontekstowych można dokonywać podając wzorce tekstowe, zgodnie z
regułami przyjętymi dla programu find.  Wymiany polegają na
zastępowaniu tekstu według tych samych reguł, co w programie change.

Numery wierszy są budowane z następujących elementów:

	n		liczba dziesiętna
	.		wiersz bieżący (kropka)
	$		wiersz ostatni
	/wzorzec/	wyszukiwanie kontekstowe wprzód
	\wzorzec\	wyszukiwanie kontekstowe wstecz

Elementy te mogą tworzyć wyrażenia zawierające + i -, na przykład:

	.+1		wiersz plus jeden
	$-5		piąty wiersz od końca

Numery wierszy mogą być oddzielane przecinkiem lub średnikiem; średnik
ustawia wiersz bieżący na ostatni z obliczonych, przed rozpoczęciem
obliczania następnego. Dyrektywy mogą być poprzedzone dowolną liczbą
numerów wierszy (oprócz dyrektywy e, f i q, dla których nie można
podać żadnego). W miarę potrzeby, będzie użyty ostatni jeden lub dwa
numery.  Jeśli podany jest jeden numer, a potrzebne są dwa, to jeden
będzie użyty dwa razy. Gdy żadne numery nie są podane, to są
przyjmowane następujące wartości domyślne:

	(.)	wiersz bieżący
	(.+1)	wiersz następny po bieżącym
	(.,.)	wiersz bieżący użyty dwa razy
	(1,$)	wszystkie wiersze

A oto dyrektywy i ich domyślne numery wierszy, w porządku
alfabetycznym:

	(.)a	dopisz tekst po wierszu (dalej następuje tekst)
	(.,.)c	wymień tekst (dalej następuje tekst)
	(.,.)dp	usuń tekst
	e plik	redaguj plik, kasując poprzednią zawartość bufora;
		zapamiętaj nazwę pliku
	f plik	wydrukuj i zapamiętaj nazwę pliku
	(.)i	wstaw tekst przed wierszem (dalej następuje tekst)
	(.,.)m w3 p	przenieś tekst za wiersz w3
	(.,.)p	drukuj tekst
	q	wyjdź z edytora
	(.)r plik	czytaj plik i dołącz za wierszem
	(.,.)s/wz/nowy/gp	zastąp występowanie wzorca wz tekstem
				nowy (g powoduje zastąpienie wszystkich
				wystąpień w wierszu)
	(1,$)w plik	zapisz plik (niczego nie zmieniając)
	(.)=p	drukuj numer wiersza
	(.+1)	drukuj jeden wiersz

Dodając końcówkę p, spowodujemy wydrukowanie ostatniego z
przetwarzanych wierszy. Wierszem bieżącym staje się zawsze ostatni z
przetwarzanych wierszy, z wyjątkiem dyrektyw f, w i =, które go nie
zmieniają.

Tekst wprowadzony po a, c lub i należy zakończyć wierszem zawierającym
tylko kropkę.

Globalne przedrostki powodują powtarzanie wykonywania dyrektyw dla
każdego wiersza, który zawiera wystąpienie podanego wzorca (g) lub go
nie zawiera (x):

	(1,$)g/wzorzec/dyrektywa
	(1,$)x/wzorzec/dyrektywa

Dyrektywa musi być różna od a, c, i, q i może być, jak zwykle,
poprzedzona numerami wierszy. Przed wykonaniem dyrektywy, wierszem
bieżącym staje się wiersz, w którym znaleziono dopasowanie.

Jeśli dla dyrektywy podano parametr plik, to edytor zachowuje się tak,
jak gdyby wcześniej była wykonana dyrektywa e plik. Pierwsza z
podanych nazw plików jest pamiętana, tak więc gdy w kolejnych
wywołaniach dyrektyw e, f, r lub w nie zostanie podana nazwa pliku,
będzie mogła być użyta nazwa zapamiętana. Nazwa pliku podana z
dyrektywą e lub f zastępuje każdą zapamiętaną nazwę.

PRZYKŁADY

Nie przesadzajmy!

*/
package main
