# Calculator Punctaj UBB

Calculator pentru punctajul de la examenul de admitere la Facultatea de Matematică și Informatică UBB.

## Descriere

Acest program calculează punctajul pentru testele grilă de la admiterea la Facultatea de Matematică și Informatică UBB, conform algoritmului oficial de notare. Programul suportă atât formatul vechi (30 întrebări) cât și formatul nou (24 întrebări).

## Algoritm de Notare

- Fiecare întrebare poate avea între 1 și 3 răspunsuri corecte
- Fiecare întrebare are un punctaj „p" care este:
  - 3.75 puncte pentru formatul cu 24 întrebări
  - 3 puncte pentru formatul cu 30 întrebări
- Pentru fiecare întrebare:
  - Dacă un răspuns corect este bifat: +p/t puncte (unde t = numărul total de răspunsuri corecte)
  - Dacă un răspuns incorect este bifat: -0.66*p/t puncte
  - Punctajul minim per întrebare este 0
  - Punctajul maxim per întrebare este p
- Se adaugă 10 puncte din oficiu

## Funcționalități

- Interfață interactivă în terminal
- Suport pentru citirea răspunsurilor corecte din fișier
- Navigare ușoară între întrebări
- Raport detaliat per întrebare la final
- Validare pentru numărul de răspunsuri corecte

## Utilizare

### Mod Interactiv
```bash
calcubb
```

### Citire din Fișier
```bash
calcubb -f raspunsuri.txt
```

### Verificare Versiune
```bash
calcubb -v
```

### Format Fișier Răspunsuri
```
24          # numărul de întrebări
0001        # răspunsul D este corect
0110        # răspunsurile B și C sunt corecte
...         # etc.
```

### Comenzi în Program
- **Tab/N**: Următoarea întrebare
- **Shift+Tab/P**: Întrebarea anterioară
- **Săgeți**: Navigare între opțiuni
- **Space/A**: Selectare/deselectare răspuns
- **Enter**: Calculează punctajul
- **Q**: Ieșire
- **R**: Reset (în ecranul de rezultat)

## Instalare

```bash
go install github.com/drclcomputers/calcubb
```

Sau descărcați executabilul direct din secțiunea Releases.

## Cerințe Sistem

- Windows/Linux/MacOS
- Terminal cu suport pentru TUI (Text User Interface)

## Licență


Acest proiect este licențiat sub termenii [Licenței MIT](LICENSE).
