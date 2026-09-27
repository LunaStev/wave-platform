---
translation_set_id: memory-model
path: language/explicit-memory-type-model
locale: de
group: language
group_order: 2
order: 9
title: 9. Zeiger, modifizierende Werte und Lebensdauern
summary: Lernen Sie das Adressieren, Dereferenzieren, Ändern eines Werts durch einen Zeiger und baumelnde Zeiger.
---

## Unterscheiden Sie zwischen Wert und Lagerort

Die Ganzzahl 42 und die Adresse, an der die Ganzzahl gespeichert ist, sind unterschiedliche Werte. Der Zeiger zeigt auf den Speicherort. Durch die Übergabe einer Adresse kann die Funktion den Speicher des Aufrufers lesen oder ändern.

In diesem Kapitel geht es um die Übernahme von Adressen, die Dereferenzierung, das Ändern des ursprünglichen Werts, die Zeigerarithmetik und die Lebensdauer. Die dynamische Zuordnung wird im nächsten Kapitel behandelt. Beginnen Sie mit Adressen lokaler Variablen und Array-Elemente.

## Adresserfassung und Dereferenzierung

<!-- wave-example: book-pointer-first -->
```wave
fun main() {
    var value: i32 = 42;
    var address: ptr<i32> = &value;

    println("value={}", value);
    println("through pointer={}", deref address);

    deref address = 99;
    println("changed={}", value);
}
```

Ausführungsergebnis:

```text
value=42
through pointer=42
changed=99
```

`&value` ruft die Adresse ab, `deref address` liest oder schreibt den Wert dieser Adresse und so weiter. Anstatt 99 in address zu speichern, wurde 99 in die Ganzzahl geschrieben, auf die address zeigt. Die Variable address selbst zeigt weiterhin auf value.

`ptr<i32>` ist ein Zeigertyp für den Zugriff auf den Speicher i32. Der Typ zeichnet keine Länge auf und bietet keine automatische Freigabe.

## Den Zeiger selbst ändern

<!-- wave-example: book-pointer-reassign -->
```wave
fun main() {
    var first: i32 = 10;
    var second: i32 = 20;
    var selected: ptr<i32> = &first;

    selected = &second;
    deref selected = 25;

    println("{} {}", first, second);
}
```

Ausführungsergebnis:

```text
10 25
```

`selected = &second` speichert eine andere Adresse in einer Zeigervariablen. Der Wert von first ändert sich nicht. Wenn Sie dann den Wert als deref schreiben, ändert sich second. Durch die Trennung von „Adressänderung“ und „Wertänderung durch Adresse“ in separate Sätze wird Verwirrung vermieden.

## Wenn Sie eine Funktion ausführen, wird das Original geändert

<!-- wave-example: book-pointer-increment -->
```wave
fun increment(value: ptr<i32>) {
    deref value = deref value + 1;
}

fun main() {
    var count: i32 = 4;

    increment(&count);
    increment(&count);

    println("{}", count);
}
```

Ausführungsergebnis:

```text
6
```

Die Funktion erhält die Adresse von count und nicht ihren Wert 4. Durch die Änderung dieses Speichers ändert sich auch count im Aufrufer. Die Funktion erfordert eine gültige, beschreibbare i32-Adresse. Das Bestehen von null verstößt gegen diese Anforderung.

Jede Funktion definiert, ob sie null akzeptiert. Ist dies nicht der Fall, muss der Anrufer eine gültige Adresse angeben. Wenn dies der Fall ist, muss die Funktion einen Pfad enthalten, der null verarbeitet.

## Funktion zur Verarbeitung von null

<!-- wave-example: book-pointer-null -->
```wave
fun try_increment(value: ptr<i32>) -> bool {
    if (value == null) {
        return false;
    }

    deref value = deref value + 1;
    return true;
}

fun main() {
    var count: i32 = 7;

    if (!try_increment(null)) {
        println("no value");
    }

    if (try_increment(&count)) {
        println("count={}", count);
    }
}
```

Ausführungsergebnis:

```text
no value
count=8
```

Eine null-Prüfung behandelt nur das Fehlen einer Adresse. Das Konvertieren einer beliebigen Nicht-null-Zahl in einen Zeiger erzeugt keinen gültigen Speicher. Lesen und Schreiben erfordern außerdem eine gültige Lebensdauer, ausreichende Größe, korrekte Ausrichtung und die entsprechenden Zugriffsberechtigungen.

## Array-Adressen und elementweise Zeigerarithmetik

<!-- wave-example: book-pointer-array -->
```wave
fun main() {
    var values: array<i32, 3> = [10, 20, 30];
    var first: ptr<i32> = &values[0];
    var second: ptr<i32> = first + 1;

    println("{}", deref first);
    println("{}", deref second);
    println("{}", deref first[2]);
}
```

Ausführungsergebnis:

```text
10
20
30
```

Durch Hinzufügen von 1 zum Zeiger wird dieser um ein Element des Zieltyps verschoben. Das nächste Element in i32 und das nächste Element in u8 haben eine unterschiedliche Anzahl an Schiebebytes. Wenn Sie `first + 1` erneut mit der Schriftgröße multiplizieren und addieren, wird es an eine unerwünschte Position verschoben.

Die Zeigerindizierung muss ebenfalls innerhalb des gültigen Bereichs erfolgen. first merkt sich die Array-Länge 3 selbst nicht und verwendet daher bei der Übergabe des Bereichs an die Funktion ein Formular, das sowohl einen Zeiger als auch eine Länge empfängt.

## Übergabe des Lesebereichs an eine Funktion

<!-- wave-example: book-pointer-range -->
```wave
fun sum(values: ptr<i32>, count: i32) -> i32 {
    var total: i32 = 0;

    for (var index: i32 = 0; index < count; index += 1) {
        total += deref values[index];
    }

    return total;
}

fun main() {
    var values: array<i32, 4> = [2, 4, 6, 8];

    println("first two={}", sum(&values[0], 2));
    println("all={}", sum(&values[0], 4));
}
```

Ausführungsergebnis:

```text
first two=6
all=20
```

Die Einheit von count ist die Anzahl der Elemente. Für diese Funktion muss der Anrufer über count lesbares i32 verfügen. Die Übergabe einer Länge, die größer als das tatsächliche Array ist, verstößt gegen den Vertrag. Auch wenn Sie die gleiche ptr<u8>·i64-Kombination API verwenden, sollten Sie in der Dokumentation prüfen, ob die Länge in Bytes oder in Elementen angegeben ist.

## Lebensdauer: Wie lange ist die Adresse gültig?

Lokale Variablen werden innerhalb der Lebensdauer des Aufrufs und Blocks verwendet. Wenn Sie die Adresse einer lokalen Variablen innerhalb einer Funktion zurückgeben, damit der Aufrufer sie später lesen kann, hat dieser Speicherplatz möglicherweise das Ende seiner Lebensdauer erreicht.

Hier ist ein schlechtes Design, das nicht umgesetzt werden sollte:

```wave
fun invalid_address() -> ptr<i32> {
    var local: i32 = 42;

    return &local;
}
```

Wenn ein Wert erforderlich ist, wird i32 zurückgegeben. Wenn in den vom Aufrufer bereitgestellten Speicherplatz geschrieben werden muss, verwendet es einen Zeiger als Eingabe. Wenn Sie separaten Speicher benötigen, der über den Anruf hinaus erhalten bleibt, weisen Sie ihn explizit zu und übertragen Sie die Verantwortung für die Freigabe.

## Zwei Zeiger, die auf denselben Speicherplatz zeigen

Durch das Kopieren eines Zeigers wird ein anderer Name erstellt, der auf dieselbe Adresse zeigt. Speicher wird nicht dupliziert.

<!-- wave-example: book-pointer-alias -->
```wave
fun main() {
    var value: i32 = 1;
    var first: ptr<i32> = &value;
    var second: ptr<i32> = first;

    deref second = 9;

    println("{} {}", value, deref first);
}
```

Ausführungsergebnis:

```text
9 9
```

Die durch second geänderten Ergebnisse sind auch durch first sichtbar. Sobald der Quellspeicher freigegeben ist, werden beide Zeiger unbrauchbar. Das Zuweisen von null zu einer Zeigervariablen ändert nicht automatisch die anderen Kopien.

## Übung: Tauschen Sie zwei ganze Zahlen aus

Schreiben Sie eine Funktion, die zwei i32-Adressen akzeptiert und deren Werte austauscht. Der erste Wert muss in einer temporären Variablen gespeichert werden, bevor er überschrieben wird. Prüfen Sie, ob der Wert auch dann erhalten bleibt, wenn Sie dieselbe Adresse zweimal übergeben.

### Vollständige Lösung

<!-- wave-example: book-pointer-swap -->
```wave
fun swap(left: ptr<i32>, right: ptr<i32>) {
    var saved: i32 = deref left;

    deref left = deref right;
    deref right = saved;
}

fun main() {
    var first: i32 = 3;
    var second: i32 = 8;

    swap(&first, &second);
    println("{} {}", first, second);

    swap(&first, &first);
    println("{}", first);
}
```

Ausführungsergebnis:

```text
8 3
8
```

Diese Funktion erfordert außerdem, dass beide Adressen auf einen gültigen, beschreibbaren Ganzzahlspeicher verweisen. Um null zu verarbeiten, fügen Sie ein Ergebnis hinzu, das Erfolg oder Misserfolg angibt, wie in try_increment.


## **Wave Explicit Memory Type Model**

Das Zeigerdesign von Wave basiert auf **Wave Explicit Memory Type Model**. Dieses Modell definiert Zeiger und Arrays als explizite Speichertypen auf Sprachebene und nicht als syntaktische Tricks oder Bibliotheksabstraktionen.

`ptr<T>` ist ein Typ, der auf die Speicheradresse zeigt, die den `T`-Wert speichert, und `array<T, N>` ist ein Speichertyp mit fester Länge, der `N`-Werte von `T` nacheinander speichert. Daher wird die Struktur von Zeigern und Arrays ebenso offengelegt wie in Funktionsargumenten, Rückgabewerten, Strukturfeldern und anderen Typen.

## null

```wave
var buffer: ptr<u8> = null;
if (buffer == null) {
    println("no buffer");
}
```

`null` ist ein Zeigerwert, der nicht auf eine gültige Speicheradresse zeigt. `null` kann nur dem Typ `ptr<T>` zugewiesen werden und kann nicht als Ganzzahl, Boolescher Wert oder Array-Wert verwendet werden.

Eine Zuordnungs- oder Suchfunktion kann `null` zurückgeben, wenn sie kein Ergebnis liefert. Suchen Sie nach `null`, bevor Sie ein solches Ergebnis dereferenzieren. Die Dereferenzierung eines `null`-Zeigers führt nicht zum Zugriff auf gültigen Speicher.

## Zeigerkonvertierung

Wenn Sie die Adresse oder eine andere Zeigerdarstellung ändern müssen, verwenden Sie `as`.

```wave
var raw: i64 = 0;
var p: ptr<u8> = raw as ptr<u8>;
```

Verwenden Sie Konvertierungen zwischen Ganzzahlen und Zeigern nur an Grenzen auf niedriger Ebene und berücksichtigen Sie die Adressbreite der Zielplattform und ABI.
