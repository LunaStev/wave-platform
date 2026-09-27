---
translation_set_id: console-io-formatting
path: language/console-io-and-formatting
locale: de
group: language
group_order: 2
order: 16
title: Konsoleneingabe, -ausgabe und -formatierung
summary: print, println, input Beschreibt Sätze und Platzhalterregeln.
---

## Eingabe-/Ausgabeanweisung

Wave stellt `print`, `println` und `input` als Eingabe-/Ausgabeanweisungen für die Konsole bereit.

```wave
fun main() {
    var count: i32 = 0;
    input("{}", count);
    print("count = ");
    println("{}", count);
}
```

Jeder Satz endet mit `;`. Das erste Argument muss ein String-Literal sein. Variablen oder berechnete Zeichenfolgen können nicht als format-Argumente verwendet werden.

## Platzhalter

Nur genau zwei Zeichen, `{}`, sind Platzhalter.

```wave
println("name = {}, score = {}", name, score);
```

Die Anzahl der Platzhalter und die Anzahl der darauffolgenden Ausdrücke müssen genau gleich sein. Wenn die Zahlen unterschiedlich sind, handelt es sich um einen Grammatikfehler.

```wave
println("{} {}", one);
// 오류: 자리표시자 2개, 값 1개
println("plain text", one);
// 오류: 자리표시자 없음, 값이 남음
```

Andere Formen von geschweiften Klammern bleiben als einfacher Text übrig. In dieser Grammatik gibt es keine benannten oder nummerierten Platzhalter.

## print und println

`print` druckt den formatierten Text unverändert und `println` fügt einen Zeilenumbruch hinzu.

```wave
print("loading...");
println("done");
```

Formatierungsargumente verwenden Skalarwerte wie Ganzzahlen, Gleitkommazahlen, Zeichenfolgen und Zeiger. Arrays und Strukturen können nicht als Formatierungsargumente verwendet werden.

## input Ziel

`input` speichert den gelesenen Wert im Ziel, daher müssen alle Ausdrücke nach format beschreibbare Speicherorte sein.

```wave
var number: i32 = 0;
input("{}", number);
```

Als Ziele können Variablen, Felder und dereferenzierte Speicherorte verwendet werden. Literale und Berechnungsergebnisse können nicht als Eingabe verwendet werden.

Wenn nicht alle Eingabewerte in den angeforderten Typ konvertiert werden können, wird das Programm mit einem Fehlerstatus beendet.

## Laufzeitgrenze

Diese Anweisungen verwenden Konsoleneingaben und -ausgaben aus der hosted-Umgebung. In einer freistehenden Umgebung muss die vom Kernel oder Gerät bereitgestellte Eingabe/Ausgabe als Funktion oder FFI-Grenze definiert werden.

## Eingabewert und -bereich

Die Eingabe bool akzeptiert nur `0` und `1`. Es interpretiert 2 nicht als true und akzeptiert die Zeichenfolge `true` nicht als dieselbe Eingabe. Die Ganzzahleingabe muss innerhalb des Bereichs der Ziel-Ganzzahlbreite liegen. Basierend auf der Gesamtbreite des Typs werden auch 128-, 256-, 512- und 1024-Bit-Ganzzahlen verarbeitet.

Formatfehler, außerhalb des Bereichs, bevor die erforderliche Eingabe EOF fehlschlägt. Das integrierte input ist keine Funktion, die einen Fehler zurückgibt und erneut eintritt, sondern eine Eingabefunktion, die den Prozess bei einem Fehler beendet. Wenn Sie eine wiederherstellbare Eingabeverarbeitung benötigen, lesen Sie die Bytes mit io und erstellen Sie einen separaten Parser.

[Praxis des Eingaberechners](/docs/de/practice/input-calculator) · [Datei und io](/docs/de/stdlib/files-io)
