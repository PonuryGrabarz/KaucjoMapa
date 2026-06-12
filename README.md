# KaucjoMapa

Aplikacja webowa do wyszukiwania najbliższych punktów, w których można oddać pojemniki zwrotne — butelkomatów, sklepów oraz innych miejsc przyjmujących opakowania kaucyjne.

Projekt ma pomagać użytkownikom szybko znaleźć najbliższy punkt zwrotu, sprawdzić jego aktualny status oraz współtworzyć bazę miejsc dzięki zgłoszeniom i potwierdzeniom społeczności.

> Nazwa **KaucjoFinder** jest robocza i może zostać zmieniona przed publikacją projektu.

---

## Opis projektu

KaucjoFinder to aplikacja webowa, która umożliwia znalezienie w okolicy najbliższego miejsca, gdzie można oddać pojemniki zwrotne, takie jak butelki, puszki lub inne opakowania objęte systemem kaucyjnym.

Aplikacja pokazuje punkty zwrotu na mapie, pozwala wyznaczyć trasę do wybranego miejsca oraz umożliwia społecznościowe aktualizowanie informacji o punktach, np. czy butelkomat jest czynny, zapełniony albo tymczasowo niedostępny.

Celem projektu jest stworzenie praktycznego narzędzia, które oszczędza czas użytkowników i wspiera recykling oraz wygodniejsze korzystanie z systemu kaucyjnego.

---

## Główne funkcjonalności

### Mapa punktów zwrotu

* Wyświetlanie butelkomatów, sklepów i innych punktów przyjmujących pojemniki zwrotne.
* Pokazywanie punktów znajdujących się najbliżej aktualnej lokalizacji użytkownika.
* Możliwość ręcznego wyszukania lokalizacji.
* Widok mapy z markerami punktów zwrotu.
* Lista najbliższych punktów posortowana według odległości.

### Nawigacja do punktu

* Wyznaczanie trasy do wybranego punktu zwrotu.
* Możliwość otwarcia trasy w zewnętrznej aplikacji mapowej, np. Google Maps, Apple Maps lub OpenStreetMap.
* Informacja o szacowanej odległości i czasie dotarcia.

### Status punktu zwrotu

Użytkownik może sprawdzić aktualny status miejsca, np.:

* czynny,
* nieczynny,
* zapełniony,
* dostępny,
* chwilowo niedostępny,
* wymaga serwisu,
* brak możliwości zwrotu w danym momencie.

Statusy mogą być dodawane i aktualizowane przez użytkowników.

### Dodawanie nowych miejsc

* Użytkownik może zgłosić nowy punkt zwrotu.
* Formularz dodawania miejsca może zawierać:

  * nazwę miejsca,
  * typ punktu,
  * adres,
  * lokalizację na mapie,
  * godziny otwarcia,
  * obsługiwane typy opakowań,
  * zdjęcie punktu,
  * dodatkowy opis.

### Społecznościowa weryfikacja informacji

* Użytkownicy mogą potwierdzać lub odrzucać zgłoszone statusy.
* Użytkownicy mogą głosować, czy dane miejsce istnieje i czy informacje są poprawne.
* Informacje często potwierdzane przez społeczność mogą otrzymywać wyższy poziom wiarygodności.
* Zgłoszenia odrzucone przez wielu użytkowników mogą trafiać do moderacji.

---

## Proponowane dodatkowe funkcjonalności

### Konta użytkowników

* Możliwość tworzenia konta i logowania.
* Historia zgłoszeń użytkownika.
* Lista dodanych punktów.
* Lista potwierdzonych statusów.
* System reputacji dla aktywnych użytkowników.

### System reputacji

* Punkty za dodawanie poprawnych miejsc.
* Punkty za potwierdzanie statusów.
* Wyższa wiarygodność zgłoszeń od użytkowników z dobrą historią.
* Odznaki, np. „Lokalny ekspert”, „Aktywny reporter”, „Eko pomocnik”.

### Zdjęcia punktów

* Dodawanie zdjęć butelkomatu lub wejścia do sklepu.
* Zdjęcie pomaga szybciej odnaleźć punkt w terenie.
* Możliwość oznaczania zdjęć jako aktualne lub nieaktualne.

### Komentarze i wskazówki

* Krótkie komentarze użytkowników przy punkcie.
* Przykładowe wskazówki:

  * „Butelkomat znajduje się przy wejściu od parkingu”.
  * „Punkt przyjmuje tylko butelki z etykietą”.
  * „W weekendy często jest zapełniony”.

### Powiadomienia

* Powiadomienie, gdy wybrany punkt ponownie będzie dostępny.
* Powiadomienie, gdy status ulubionego punktu się zmieni.
* Możliwość obserwowania punktów w swojej okolicy.

### Ulubione punkty

* Dodawanie często używanych punktów do ulubionych.
* Szybki dostęp do najbliższego sprawdzonego punktu.
* Informacja o ostatnim znanym statusie ulubionych miejsc.

### Panel administracyjny

* Moderacja zgłoszeń.
* Zarządzanie punktami zwrotu.
* Usuwanie duplikatów.
* Blokowanie fałszywych lub spamerskich zgłoszeń.
* Podgląd statystyk aktywności.

---

## Statusy punktów zwrotu

| Status          | Opis                                                          |
| --------------- | ------------------------------------------------------------- |
| Dostępny        | Punkt działa i można oddać pojemniki.                         |
| Zapełniony      | Butelkomat jest pełny i nie przyjmuje kolejnych opakowań.     |
| Nieczynny       | Punkt jest obecnie nieczynny.                                 |
| Awaria          | Automat lub punkt ma problem techniczny.                      |
| Brak informacji | Nie ma aktualnych danych o statusie.                          |
| Do weryfikacji  | Status został zgłoszony, ale nie został jeszcze potwierdzony. |

Każdy status powinien posiadać datę ostatniej aktualizacji oraz informację o liczbie potwierdzeń i odrzuceń.

---

## Przykładowy przepływ użytkownika

1. Użytkownik otwiera aplikację.
2. Aplikacja prosi o dostęp do lokalizacji lub pozwala wpisać adres ręcznie.
3. Użytkownik widzi mapę z najbliższymi punktami zwrotu.
4. Wybiera najbliższy dostępny punkt.
5. Sprawdza status, godziny otwarcia i przyjmowane opakowania.
6. Uruchamia nawigację do wybranego miejsca.
7. Po dotarciu aktualizuje status, np. „dostępny” albo „zapełniony”.
8. Inni użytkownicy mogą potwierdzić lub odrzucić tę informację.

---

## Roadmapa

### MVP

* [ ] Widok mapy z punktami zwrotu
* [ ] Pobieranie lokalizacji użytkownika
* [ ] Lista najbliższych punktów
* [ ] Szczegóły punktu
* [ ] Wyznaczanie trasy do punktu
* [ ] Dodawanie nowego punktu
* [ ] Aktualizacja statusu punktu
* [ ] Potwierdzanie i odrzucanie statusów
* [ ] Filtrowanie punktów po typie i statusie

### Wersja 1.0

* [ ] Konta użytkowników
* [ ] Historia zgłoszeń
* [ ] System reputacji
* [ ] Zdjęcia punktów
* [ ] Panel administracyjny
* [ ] Moderacja zgłoszeń
* [ ] Ulubione punkty

### Wersja 2.0

* [ ] Powiadomienia o zmianie statusu
* [ ] Ranking punktów
* [ ] Integracja z zewnętrznymi operatorami
* [ ] Tryb PWA
* [ ] Import danych z plików lub API
* [ ] Statystyki ekologiczne użytkownika

---

## Status projektu

Projekt jest na etapie koncepcji / MVP.

Najważniejszym celem pierwszej wersji jest umożliwienie użytkownikowi szybkiego znalezienia najbliższego działającego punktu zwrotu oraz aktualizowania informacji przez społeczność.
