---
translation_set_id: whale-ir-reference
path: whale/ir-reference
locale: de
group: whale
group_order: 1
order: 6
title: Whale-IR-Referenz
summary: Beschreibt Typen, Bezeichner, Funktionsgültigkeit, Auswertungsreihenfolge und Austauschformat.
---

## Module und Bezeichner

Ein Modul besteht aus Zielinformationen, globalen Definitionen und Funktionen. Werte haben explizite Typen. Das Frontend löst den Namen, den Typ, die Überladung und die Generika der Quellsprache auf und generiert typed IR.

Funktionen und globale Variablen verwenden unterschiedliche interne Namensräume. Daher können Funktionen und Variablen denselben Namen haben. Der interne Bezeichner unterscheidet sich vom externen Verbindungsnamen `link_name` und der externe Name wird vom Frontend angegeben. Whale löst externe Konflikte nicht durch die automatische Generierung neuer Namen. Bitte beachten Sie [Symbole und Links](assembler-linker).

Jede Wertdefinition verfügt über einen Bezeichner. Definitionen können nicht dupliziert werden und Typmetadaten müssen mit dem in der Definition angegebenen Typ übereinstimmen. Der Name allein identifiziert keine Definitionen, selbst wenn Deklarationen mit demselben Namen einander verdecken.

## IR Konfiguration und Auslesen

Unten finden Sie ein vollständiges Rust-Beispiel, das die Funktion mit der Kiste `ir` erstellt und überprüft und dann typed IR ausgibt.

```rust
use ir::{ModuleBuilder, Target, Type};

fn main() {
    let target = Target::lookup("x86_64-whale-linux").unwrap();
    let mut module = ModuleBuilder::new(target.name(), target.data_layout());
    let mut function = module.begin_function("answer", vec![], Type::I32);
    let left = function.const_i32(40);
    let right = function.const_i32(2);
    let answer = function.add(Type::I32, left, right);
    function.ret(Some(answer));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    print!("{}", ir::print_module(&module));
}
```

Der Drucker gibt IR aus:

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "answer": whale () -> i32, linkage internal

  fn @answer() -> i32, id @f0 {
  entry:
    %v0: i32 = const i32 40
    %v1: i32 = const i32 2
    %v2: i32 = add i32 %v0, %v1
    ret i32 %v2
  }

}
```

`%v0` und `%v1` sind Definitionen der Konstante i32. Verwenden Sie `%v2` definiert durch `add` als Rückgabewert von i32 der Funktion. Wenn Sie den Rückgabewert in Bool ändern, führt dies zu einem Überprüfungsfehler, da er nicht mit der Funktionssignatur übereinstimmt. Auch wenn beide Operanden Konstanten sind, behält O0 die Additionsanweisung bei.

Bei diesem Code handelt es sich um eine tatsächliche Druckerausgabe und nicht um eine Eingabedatei, die an einen Textparser übergeben wird. Textparsing und Ausführung von IR werden noch nicht unterstützt, derzeit kann dieses Modul als Rust builder konfiguriert werden.

## Typ

|Typ|Bedeutung|
| --- | --- |
| `bool` |Logischer Wert false oder true|
| `i1`, `i8`, `i16`, `i32`, `i64`, `i128` |Ganzzahl mit Vorzeichen und angegebener Bitbreite|
| `u1`, `u8`, `u16`, `u32`, `u64`, `u128` |Ganzzahl ohne Vorzeichen der angegebenen Bitbreite|
| `f16`, `f32`, `f64` |Gleitkommawert mit angegebener Bitbreite|
| `ptr<T>` |T Zeiger auf Typwert|
| `fnptr<signature>` | Aufrufbarer Zeiger mit genauen Parameter-/Ergebnistypen und Aufrufkonvention |
| `array<T, N>` |N Elemente des gleichen Typs|
| `struct{T, ...}` |Geordnetes Strukturfeld|
| `tuple<T, ...>` |geordnete Tupelelemente|
| `void` |Keine Ergebnisse|

`bool`, `i1` und `u1` sind unterschiedliche Typen. signed `i1` steht für −1 und 0 und unsigned `u1` steht für 0 und 1. Die Ganzzahl 1 ist keine implizite logische Bedingung. Bedingter Zweig, die Bedingung von Select, `trap_if` erfordert den Operanden Bool.

Die Speichergröße wird nicht allein durch die Anzahl der Bits im Wert bestimmt und folgt [Ziellayout](memory-model). Der Wert von `i1` beträgt beispielsweise 1 Bit, belegt aber mindestens 1 Byte im Speicher.

## Funktionen und Aufrufe

Die Funktion gibt alle Parameter, Ergebnistypen, Aufrufkonventionen und linkage an. Direkte und indirekte Anrufe müssen mit der Signatur des Anrufers übereinstimmen. Der Anruf void hat kein Ergebnis ID. Ein Aufruf von nonvoid behält die Ergebnisdefinition bei, auch wenn O0 das Ergebnis nicht verwendet.

Die Rückgabe muss mit dem Ergebnistyp der Funktion übereinstimmen. Die Rückgabe void enthält keinen Wert, und die Rückgabe nonvoid enthält einen Wert des deklarierten Ergebnistyps.

### Deklarationen, Identitäten und Anrufe

`Module.declarations` zeichnet den `FunctionId` jeder Funktion, den Namen, die vollständige Signatur, die Verknüpfung und den Namen des externen Links auf. Eine Definition bezieht sich auf diese Identität; Parameter- und Rückgabetypen müssen mit seiner Deklaration übereinstimmen. Identische wiederholte Deklarationen werden über `declare_function` in dieselbe ID aufgelöst; Konflikte und doppelte Definitionen sind Fehler. Eine interne Deklaration benötigt einen Hauptteil im Modul. Eine externe Deklaration kann bis zur Verknüpfung ungelöst sein oder einen exportierten Textkörper haben. Interne Funktionen haben kein `link_name`; Externe Funktionen erfordern einen expliziten, nicht leeren Namen ohne NUL. Zwei unterschiedliche Funktionsdeklarationen können nicht denselben externen Namen beanspruchen. Globale und Funktionen verwenden weiterhin separate interne Namespaces.

Registrieren Sie Deklarationen vor dem Erstellen von Körpern mit `begin_declared_function`, um Weiterleitungsaufrufe und Rekursion zu unterstützen. `begin_function` bleibt eine Annehmlichkeit für eine neue interne Whale-Funktion. Die überprüften APIs `declare_function`, `begin_declared_function`, `function_addr`, `null_function` und `call` geben `Result` zurück; Bei einem abgelehnten Aufruf wird weder eine Anweisung angehängt noch ihre Ergebnis-ID zugewiesen.

Das folgende vollständige Rust-Programm deklariert eine externe Funktion, übernimmt deren typisierte Adresse und gibt sowohl direkte als auch indirekte Aufrufe aus:

```rust
use ir::{Callee, CallingConvention, DataLayout, FunctionSignature, Linkage, ModuleBuilder, Type};

fn main() {
    let mut module = ModuleBuilder::new("x86_64-whale-linux", DataLayout::default_64bit_le());
    let signature = FunctionSignature {
        params: vec![Type::I32], ret: Type::I32,
        convention: CallingConvention::SysV64, variadic: false,
    };
    let identity = module.declare_function(
        "identity", signature, Linkage::External, Some("identity_i32".into()),
    ).unwrap();
    let mut function = module.begin_function("answer", vec![], Type::I32);
    let input = function.const_i32(42);
    let callback = function.function_addr(identity).unwrap();
    // The direct call's result remains defined even though it is unused.
    function.call(Callee::Direct(identity), vec![input]).unwrap();
    let result = function.call(Callee::Indirect(callback), vec![input]).unwrap().unwrap();
    function.ret(Some(result));
    function.finish();
    let module = module.finish();
    ir::verify_module(&module).unwrap();
    print!("{}", ir::print_module(&module));
}
```

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "identity": sysv64 (i32) -> i32, linkage external, link_name "identity_i32"
  declare @f1 "answer": whale () -> i32, linkage internal

  fn @answer() -> i32, id @f1 {
  entry:
    %v0: i32 = const i32 42
    %v1: fnptr<sysv64 (i32) -> i32> = function_addr @f0
    %v2: i32 = call sysv64 i32 @f0(%v0)
    %v3: i32 = call sysv64 i32 indirect %v1(%v0)
    ret i32 %v3
  }

}
```

`Callee::Direct(FunctionId)` wird über die Deklarationstabelle aufgelöst; `Callee::Indirect(ValueId)` erfordert einen `Type::FnPtr(FunctionSignature)`-Wert. Die Signatur umfasst alle Parametertypen, den Ergebnistyp und `CallingConvention::{Whale, SysV64}`. Es bleibt durch Kopien, Speicherung, Parameter, Rückgaben, Phi und Auswahl erhalten. Datenzeiger und Ganzzahlwerte sind nicht aufrufbar. Umwandlungen mit Funktionszeigertypen werden abgelehnt. Durch das Ändern einer Typanmerkung kann eine aufrufbare Signatur nicht geändert werden. Ein Funktionszeiger verfügt über einen 64-Bit-Adressspeicher auf diesem Ziel; Dadurch werden selbst keine Runtime-Shadow-Metadaten implementiert.

Arität, genaue Argument-/Ergebnistypen, Vorhandensein der Ergebnis-ID und Aufrufkonvention müssen übereinstimmen. Es gibt keine impliziten Konvertierungen. Der indirekte Angerufene muss den Anruf genauso dominieren wie seine Argumente. `variadic: true`, Parameter vom Typ void und SysV64 Aggregatparameter/Ergebnissignaturen werden abgelehnt. Whale aggregierte Signaturen können in IR dargestellt werden; Für keine der beiden Konventionen sind die native ABI-Klassifizierung und die Ausgabe von Maschinenaufrufen noch verfügbar.

`null_function(signature)` stellt einen typisierten Nullfunktionszeiger dar. Der Aufruf erfolgt durch ein gut eingegebenes IR mit einem erforderlichen Laufzeit-Trap vor der Eingabe eines Angerufenen. Ein Ziel ungleich Null, das ungültig, abgelaufen oder mit der überprüften Signatur nicht kompatibel ist, muss ebenfalls abgefangen werden. Diese Laufzeitprüfungen und die Verwaltung der Lebensdauer ausländischer Rückrufe warten auf die Interpreter-/native Ausführungsschicht. Ein erfolgreicher Verifizierer bedeutet nicht, dass beliebige externe Adressen sicher sind.

### AST-Aufrufformen

Dies sind Ausdrucksfragmente innerhalb eines AST-Format-2-Programms:

```json
{"Call":{"callee":{"Direct":"increment"},"args":[{"Lit":{"Int":{"bits":32,"signed":true,"value":"41"}}}]}}
```

```json
{"Call":{"callee":{"Indirect":{"FunctionRef":"increment"}},"args":[{"Lit":{"Int":{"bits":32,"signed":true,"value":"41"}}}]}}
```

`Direct` und `FunctionRef` verwenden den Funktionsnamensraum, auch wenn eine Variable denselben Namen hat. `Indirect` wertet zuerst seinen Ausdruck und dann die Argumente von links nach rechts aus. Ein void-Aufruf ist als `ExprStmt` gültig, jedoch nicht als Variableninitialisierer, Argument, Operand oder Rückgabewert. Aufrufe und Funktionsverweise sind keine numerischen Konstantenausdrücke zur Kompilierungszeit. `NullFunction` nimmt ein Signaturobjekt mit den Feldern `params`, `ret`, `convention` und `variadic`.

Das [vollständige JSON-Beispiel](https://github.com/wavefnd/Whale/blob/master/ir/tests/fixtures/ast-v2-calls.json) speichert einen Rückruf und ruft ihn vor einem externen Anruf auf. Senken Sie es ab mit:

```sh
cargo run --locked --features socket-cli -- ir lower ir/tests/fixtures/ast-v2-calls.json
```

Sein [erwarteter IR](https://github.com/wavefnd/Whale/blob/master/ir/tests/fixtures/calls-v2.wir) wird in den Tieferlegungstests überprüft. Funktionsidentität und Linknamen werden an der Grenze IR dargestellt; Ihre Erhaltung durch native Objektgenerierung und -verknüpfung ist immer noch eine separate Arbeit.

## Verfügbarkeit von Blöcken und Werten

Jeder Block hat eine eindeutige Kennung und genau einen terminator. Verzweigungsziele müssen zur gleichen Funktion gehören. Der Eintrittsblock muss vorhanden sein und darf keine Vorderkante und phi haben. Wenn Sie eine Schleife erstellen, verzweigen Sie vom Eingangsblock zu einem separaten Schleifenkopf.

Die Definition des Werts im ausführbaren Pfad sollte seinen typischen Verwendungspunkt bestimmen. Das bedeutet, dass alle Pfade vom Einstiegspunkt bis zum Verwendungspunkt diese Definition durchlaufen müssen. Im selben Block müssen Definitionen vor Verwendungen stehen. Die Speicherreihenfolge der Blöcke bestimmt nicht die Dominanz.

Gehen Sie davon aus, dass der Einstiegspunkt zu left oder right verzweigt und dann bei join zusammentrifft. Werte, die nur in left definiert sind, können nicht als allgemeine Werte in join verwendet werden. Dies liegt daran, dass der durch right verlaufende Pfad undefiniert ist. Die Werte jedes vorangehenden Blocks müssen zu phi kombiniert werden, das als Eingabe empfangen wird.

Auch nicht erreichbare Blöcke bleiben im Modul erhalten. Der Verifizierer überprüft kontinuierlich die Kennung, den Typ, den Operanden und die Verzweigungsstruktur des Blocks. Die Definition eines nicht erreichbaren Blocks kann keinen Mehrwert für die normale Nutzung eines erreichbaren Pfads liefern.

## Phi Befehl

phi wird vor allen regulären Befehlen im Block platziert. Für jeden unterschiedlichen Vorgängerblock ist genau eine Eingabe erforderlich. Der Eingabewert muss vom Typ phi sein und am Ende des entsprechenden Vorgängersatzes verfügbar sein.

Auch wenn in einem vorhergehenden Block mehrere Kanten vorhanden sind, gibt es nur eine Eingabe. Schleife phi kann sich auf den Wert eines Blocks beziehen, der später in der Modulspeicherreihenfolge erscheint, solange es sich um einen Wert handelt, der an der Wiederholungskante berechnet wird. Fehlende, doppelte oder irrelevante vorangehende Blöcke und falsch eingegebene Eingaben sind Validierungsfehler.

## Bewertung und Auswahl

Whale AST wertet Aufrufziele und Unterausdrücke von links aus, wo sie angegeben sind. Das Frontend drückt die Kurzschlussauswertung als Kontrollflusszweig aus.

Select wählt einen der bereits berechneten Werte aus. Die Berechnung beider Eingaben wird nicht ausgelassen. Selbst wenn Sie beispielsweise einen sicheren Wert wählen, können Sie das Auftreten von trap bei der Berechnung anderer Eingaben nicht vermeiden. Berechnungen, die nur auf einem bestimmten Pfad ausgeführt werden müssen, sollten in einen bedingten Block eingefügt werden.

## Verifizierungsabteilung trap

Ungültiges IR wird während der Überprüfungsphase abgelehnt. Verstöße gegen die Ausführungsbedingung werden mit dem definierten trap behandelt, und `undef` und `poison` sind keine akzeptablen Werte. Missbrauch von builder, doppelte Definition, Hinzufügung eines zweiten terminator sollte als struktureller Fehler zurückgegeben werden.

trap enthält den Grund·Quellenort·IR ID und stoppt dann die Ausführung. Durch Ausführen von native wird das Programm beendet und der Interpreter API gibt den Fehler Trap zurück. Behält frühere Nebenwirkungen bei, garantiert jedoch nicht den Puffer flush·Destruktoraufruf·Stack unwinding.

Diese Garantie gilt für verifizierte IR- und Trace-Speicher. Die externe C·Rohadresse·Inline-Assembly verfügt über einen separaten Vertrag und erkennt nicht immer Verstöße außerhalb seiner Grenzen. Bitte beachten Sie [Speichermodell](memory-model).

## Austauschformat und Textdarstellung

AST und typed IR verwenden das entsprechende format version und das gemeinsame semantics version. Der Leser sollte den unversionierten/unbekannten Versions-/Feld-/Funktions-/Duplikatschlüssel JSON ablehnen. Konstrukteure sollten nicht davon ausgehen, dass nicht unterstützte Eigenschaften stillschweigend ignoriert werden.

Ganzzahlen werden als Bitbreite·signedness·String-Zahlen übergeben. Gleitkommakonstanten werden als Breite und genaue Bitfolge übergeben. Der Text round-trip in IR muss den Namen·ID·Typ·Konstante·Sequenz·Eigenschaft·Metadaten beibehalten. Leerzeichen und Kommentarplatzierung unterliegen nicht der Aufbewahrung.

Sie können die folgenden Skalarverträge AST JSON verwenden. Die Ausgabe typed IR enthält Versionsinformationen, ein vollständiger Roundtrip-Austausch mit dem Textparser wird jedoch noch nicht unterstützt.


### Angegebene Version AST JSON

Speichern Sie Folgendes als `program.json`. Alle vier Umschlagfelder sind Pflichtfelder. `program` enthält die erforderlichen Arrays `declarations`, `globals` und `functions`, die möglicherweise leer sind. Funktionsname, Parameter, Rückgabetyp, Hauptteil, `convention` und `linkage` sind erforderlich. `link_name` kann für interne Funktionen fehlen/null sein und muss eine nicht leere Zeichenfolge ohne NUL für externe Funktionen sein. Jede Aufzählung verwendet entweder ihren Einheitennamen oder ein einzelnes Variantenschlüsselobjekt. Einheitenvarianten akzeptieren auch ein Nullwertobjekt, z. B. `{"Void":null}`; Der Encoder gibt den Gerätenamen `"Void"` aus. `VarDecl.init` kann fehlen oder null sein; weitere erforderliche Felder müssen vorhanden sein.

```json
{
  "format_version": 2,
  "semantics_version": 1,
  "features": [],
  "program": {
    "globals": [],
    "functions": [
      {
        "name": "answer",
        "parameters": [],
        "return_type": {
          "Int": {
            "bits": 128,
            "signed": false
          }
        },
        "body": [
          {
            "Return": {
              "Lit": {
                "Int": {
                  "bits": 128,
                  "signed": false,
                  "value": "340282366920938463463374607431768211455"
                }
              }
            }
          }
        ],
        "convention": "Whale",
        "linkage": "Internal",
        "link_name": null
      }
    ],
    "declarations": []
  }
}
```

```sh
cargo run --locked --features socket-cli -- ir lower program.json
```

```text
module {
  format_version 2
  semantics_version 1
  target "x86_64-whale-linux"
  datalayout { ptr=64, endian=little }

  declare @f0 "answer": whale () -> u128, linkage internal

  fn @answer() -> u128, id @f0 {
  entry:
    %v0: u128 = const u128 340282366920938463463374607431768211455
    ret u128 %v0
  }

}
```

Die Ganzzahl `value` ist eine Dezimalzeichenfolge. Signed Eine Zahl wird nach dem optionalen Minus einer Ganzzahl verwendet und Leerzeichen, Pluszeichen, Exponenten und Trennzeichen sind nicht zulässig. Der zulässige Bereich wird durch die angegebene Breite und signedness bestimmt. Das obige `u128::MAX` bleibt unverändert bis JSON und lowering erhalten. Negative unsigned-Werte oder Werte außerhalb des Bereichs sind keine wrap und sind Fehler. Der Wert Float verwendet eine hexadezimale Bitfolge mit der genauen Breite, die in [Numerische Operationen](numeric-operations) beschrieben ist.

`format_version` ist 2 für dieses AST-Format; `semantics_version` ist 1. `features` muss ein leeres Array sein. Unbekannte Felder, Versionen, Funktionen, doppelte rohe JSON-Schlüssel (einschließlich maskierter äquivalenter Schlüssel) und nachgestellte Werte sind Fehler, auch mit `--no-verify`. Der Bibliothekseinstiegspunkt ist `ir::lower_ast::interchange::decode`; `encode` gibt den Umschlag aus. `decode` ist standardmäßig auf ein Quellbyte-Limit von 8 MiB eingestellt; `decode_with_limit` akzeptiert ein Anruflimit. JSON Die Verschachtelung ist begrenzt. Verwenden Sie diesen Rohdecoder, anstatt ihn in eine generische Karte zu analysieren, die bereits doppelte Schlüssel verwerfen könnte.

[Das vollständige JSON Schema](https://github.com/wavefnd/Whale/blob/master/ir/schema/ast-v2.schema.json) spezifiziert Formen, erforderliche Felder und Varianten. Zusätzlich kommen Bereichs-/Typprüfungen und die Erkennung doppelter Schlüssel zum Einsatz. Die Teilmenge der Skalarreduzierung umfasst Literale, Variablen/Konstanten, Add/Sub/Mul, Vergleiche, Zuweisungen, If/While, Return und Break/Continue. Funktionsreferenzen, direkte Aufrufe und indirekte Aufrufe werden unterstützt; Aggregatausdrücke werden nicht unterstützt. `Opaque` ist im Schema darstellbar, wird jedoch durch Absenken nicht unterstützt.

Bei der Migration müssen alte reine Programmnutzlasten umschlossen und numerische JSON-Literale durch dezimale Ganzzahlzeichenfolgen oder Gleitkommabitzeichenfolgen ersetzt werden. Alte unversionierte Payloads werden abgelehnt. Nutzlasten von Format 1 müssen in Format 2 migriert werden: Fügen Sie `program.declarations` (ein leeres Array, wenn es nicht verwendet wird) und explizites `convention`/`linkage` in Definitionen hinzu. AST und eingegebene IR Versionsnummern sind unabhängig; beide sind jetzt 2, mit Semantikversion 1.

### Abgelehnte Eingabe und CLI Wiederherstellung

Speichern Sie die folgende vollständige Eingabe als `invalid.json`.

```json
{
  "format_version": 99,
  "semantics_version": 1,
  "features": [],
  "program": {
    "globals": [],
    "functions": [],
    "declarations": []
  }
}
```

```sh
cargo run --locked --features socket-cli -- ir lower invalid.json -o rejected.wir
```

```text
Failed to parse socket JSON: unsupported AST format_version 99; expected 2
```

Der Befehl wird mit einem Status ungleich Null beendet und erstellt keine neue Ausgabe und überschreibt keine vorhandenen Dateien. Typkonflikte schlagen ebenfalls fehl, bevor die Ausgabe veröffentlicht wird. Ohne `socket-cli` erstellte Binärdateien werden mit Status 2 beendet und geben einen Wiederherstellungsbefehl aus, der `--features socket-cli` enthält.
