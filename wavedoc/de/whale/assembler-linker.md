---
translation_set_id: whale-assembler-linker
path: whale/assembler-linker
locale: de
group: whale
group_order: 1
order: 11
title: Assemblierung und statisches Linken
summary: Beschreibt Operandenkodierung, Abschnittsplatzierung, Symbolbindung und Ausführungseinstiegspunkte.
---

## Baugruppen und Objekte

Der Whale-Assembler wandelt die AMD64-Anweisung in Maschinencode-Bytes und Verschiebungsinformationen um. Das Objekt ELF64 enthält diese Informationen und Abschnitte/Symbole. Die Assemblierung erfolgt durch die eigene Implementierung von Whale und erfordert keinen externen Assembler.

Verschiebbare Objekte können Referenzen haben, deren endgültige Adresse noch nicht bekannt ist. Der Prozess zum Auflösen dieser Adresse ist ein Link. Eine erfolgreiche Assemblierung sollte nicht darauf schließen lassen, dass alle externen Symbole aufgelöst wurden oder dass eine ausführbare Datei erstellt wurde.

## Bauen Sie die Funktion zusammen

Speichern Sie den folgenden Code als `answer.asm`.

```asm
section .text
global answer

answer:
    mov eax, 42
    ret
```

```sh
whale asm --amd64 answer.asm -o answer.o
```

Der Inhalt von `.text` ist `b8 2a 00 00 00 c3`, und auf `mov eax, 42` folgt `ret`. Das Objekt ELF64 legt `answer` frei. Es handelt sich um eine aufrufbare Funktion ohne Prozessstartcode und nicht um eine ausführbare Datei. Die Anweisungen in diesem Beispiel können vom aktuellen Assembler verarbeitet werden.

## Objekte konstruieren mit Rust API

Hier ist ein vollständiges Beispiel für das Schreiben derselben Funktionsbytes in die Kiste `object`. Gibt Ausgabeziele und globale Symbole an.

```rust
use object::{ObjectFile, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility, Target};

fn main() {
    let mut object = ObjectFile::with_target(Target::X86_64WhaleLinux.object_target());
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xb8, 0x2a, 0x00, 0x00, 0x00, 0xc3];
    object.symbols.push(ObjectSymbol {
        name: "answer".into(),
        section_index: Some(text),
        value: 0,
        size: 6,
        binding: SymbolBinding::Global,
        visibility: SymbolVisibility::Default,
    });
    let elf = object.write().unwrap();
    assert_eq!(&elf[..7], b"\x7fELF\x02\x01\x01");
    assert_eq!(u16::from_le_bytes([elf[18], elf[19]]), 62);
    std::fs::write("answer.o", elf).unwrap();
}
```

Assertion prüft die ELF-Klasse, die Bytereihenfolge und den machine-Bezeichner. `value: 0` ist `.text` innerhalb von offset und `size: 6` ist die Bytegröße des Symbols. Ungültige Abschnittsverweise oder -bereiche sind Serialisierungsfehler. Lehnt andere machine- oder Byte-Bestellungen ab, ohne sie als AMD64 zu markieren.

## Interpretieren Sie die Symbole zweier Objekte

Derzeit bietet die Kiste `linker` eine Symbolinterpretation. Das folgende ausführbare Beispiel definiert ein lokales `helper` für zwei Objekte und prüft dann, ob ein Fehler auftritt, wenn derselbe Name zweimal veröffentlicht wird.

```rust
use linker::core::symbol_table::{SymbolKey, SymbolTable};
use object::{
    ObjectFile, ObjectFormat, ObjectSymbol, SectionKind, SymbolBinding, SymbolVisibility,
};

fn input(binding: SymbolBinding) -> ObjectFile {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 1);
    object.sections[text].data = vec![0xc3];
    object.symbols.push(ObjectSymbol {
        name: "helper".into(),
        section_index: Some(text),
        value: 0,
        size: 1,
        binding,
        visibility: SymbolVisibility::Default,
    });
    object
}

fn main() {
    let mut symbols = SymbolTable::new();
    symbols
        .resolve(&[input(SymbolBinding::Local), input(SymbolBinding::Local)])
        .unwrap();
    for object_index in 0..2 {
        let key = SymbolKey::Local {
            object_index,
            name: "helper".into(),
        };
        assert_eq!(symbols.symbols[&key].object_index, Some(object_index));
    }
    let error = symbols
        .resolve(&[input(SymbolBinding::Global), input(SymbolBinding::Global)])
        .unwrap_err();
    assert_eq!(error, "Duplicate global symbol: helper");
    println!("{error}");
}
```

Ausgabe:

```text
Duplicate global symbol: helper
```

Die beiden lokalen Definitionen haben separate Schlüssel mit unterschiedlichem `object_index`. Die beiden globalen Definitionen stehen im Widerspruch. Dieses Beispiel interpretiert direkt die Symbole des Objekts, das Sie im Speicher erstellt haben. `.o` Das Lesen von Dateien, das Anwenden einer Verschiebung oder die Ausgabe ausführbarer Dateien wird nicht durchgeführt. Von den unten beschriebenen vollständigen Definitionsauswahlrichtlinien wurden weak-Prioritäten usw. noch nicht implementiert.

## Literale und Speicheroperanden

Literale behalten ihre Breite und Beschilderung, bis der Codierungsbereich der eigentlichen Anweisung untersucht wird. Werte außerhalb des zulässigen Bereichs sollten zu einem Fehler führen und nicht stillschweigend abgeschnitten werden.

Eine explizite Größenangabe ist erforderlich, wenn die Speicherbreite nicht aus anderen Informationen im Befehl ermittelt werden kann. Beispielsweise kann ein Registeroperand die Breite bestimmen, dies kann jedoch mehrdeutig sein, wenn nur Speicher und ein unmittelbarer Wert vorhanden sind. Der Assembler sollte nicht zufällig mehrdeutige Breiten erraten.

Der aus einem Symbol bestehende Speicheroperand in AMD64 ist grundsätzlich RIP-relative. Wählen Sie die Adressierungsmethode als explizit rel/abs aus. Das unbekannte escape in einem Literal ist ein Fehler.

## Abschnitte und Ausrichtung

|Abschnittsinhalt|Ausrichtungsverhalten|
| --- | --- |
|Code|NOP Befehl einfügen|
|initialisierte Daten|0 Bytes einfügen|
| BSS |Erhöhen Sie die Größe des logischen Speichers, ohne eine Datei hinzuzufügen payload|

Es wird zwischen Dateigröße und Speichergröße unterschieden. BSS reserviert Speicher, erfordert jedoch nicht die gleiche Größe von 0 Bytes, um in der Objektdatei gespeichert zu werden. Benutzerdefinierte Abschnitte verfügen über Eigenschaften und Symbole, die Bindungs- und Typinformationen verwalten.

### Payload Reserve BSS ohne Zuordnung

Speichern Sie Folgendes als `buffer.asm`. Logik BSS reserviert 1 TiB, weist aber 1 TiB beim Zusammenbau nicht zu oder schreibt es nicht.

```asm
section .bss
global buffer
buffer:
    resb 1099511627776
buffer_end:

section .text
    ret
```

```sh
whale asm --amd64 buffer.asm -o buffer.o
```

Der Wert innerhalb des Abschnitts von `buffer` ist 0 und der Wert von `buffer_end` ist 1099511627776. Der `.bss`-Header ist `SHT_NOBITS` mit dieser Größe und es gibt keine Datei payload. Wenn Sie anschließend zu `.bss` zurückkehren, verwenden Sie weiterhin die Logik offset. Datenanweisungen mit Nullen erhöhen ebenfalls die logische Größe. Ein Anfangswert ungleich Null, eine innerhalb von BSS anzuwendende Verschiebung und Befehle innerhalb von BSS werden abgelehnt.

Die gleiche Unterscheidung kann für Objekte und Linker verwendet werden API.

```rust
use linker::core::layout::Layout;
use object::{ObjectFile, ObjectFormat, SectionKind};

fn main() {
    let mut object = ObjectFile::new(ObjectFormat::ELF64);
    let text = object.add_section(".text", SectionKind::Text, 16);
    object.sections[text].data = vec![0xc3];
    let bss = object.add_section(".bss", SectionKind::Bss, 16);
    object.sections[bss].zero_fill = 1 << 40;
    assert!(object.sections[bss].data.is_empty());

    let elf = object.write_with_limit(4096).unwrap();
    assert!(elf.len() < 4096);
    let layout = Layout::compute(&[object], 0x1000).unwrap();
    assert_eq!(layout.file_size, 1);
    assert_eq!(layout.memory_size, (1 << 40) + 16);
    assert_eq!(layout.sections[bss].memory_address, 0x1010);
    assert_eq!(layout.sections[bss].file_size, 0);
}
```

```text
section  memory address  file bytes       memory bytes
.text    0x1000          1                1
.bss     0x1010          0                1099511627776
```

`Section::zero_fill` ist die Anzahl der zusätzlichen BSS Bytes, die nicht in der Datei gespeichert sind. Die getestete Speichergröße beträgt `data.len() + zero_fill`. Der vorhandene mit Nullen aufgefüllte BSS `data` wird ebenfalls akzeptiert, seine Zuweisung kann jedoch mit einem leeren `data`-Vektor vermieden werden. Andere Abschnitte als BSS müssen `zero_fill == 0` sein. Ein physischer Speicherplatz von 0 bedeutet nicht den Initialisierungsstatus IR.

`Layout::compute` gibt `Result` zurück, das die Korrespondenz zwischen Eingabeobjekt und Abschnitt, Ausrichtung, Datei offset, Speicheradresse und beide Größen aller Abschnitte aufzeichnet. Adress- und Ausrichtungsarithmetik werden überprüft, die Eingabereihenfolge bleibt erhalten und die Objektausrichtung 0 wird als uneingeschränkt behandelt (Ausrichtung 1). BSS bewegt den Dateicursor nicht. Dies ist die payload-Bereitstellung mit dem ausführbaren Header load segment mit Zugriffsrechten. Die Durchführung der Verschiebung ist eine separate Aufgabe.

ELF writer lehnt overflow und Kürzung beim Reduzieren der Feldbreite ab. Erweiterte Abschnittsnummern werden nicht unterstützt und die Gesamtzahl der Kopfzeilen einschließlich der erstellten Tabellen- und Neuanordnungsabschnitte muss kleiner als `0xff00` sein. Der Standardgrenzwert für die serialisierte Ausgabe beträgt 256 MiB. Sie können ein Byte-Limit einschließlich der Tabelle padding· mit `ObjectFile::write_with_limit` oder `write_elf_with_limit` angeben. Das Limit umfasst nicht die Speichergröße BSS, die nicht in einer Datei gespeichert ist. Die Ausgabegröße wird überprüft, bevor der endgültige Byte-Vektor zugewiesen wird.

## Symbolidentifikation

Funktionen und Variablen verwenden innerhalb von IR unterschiedliche Bezeichner. Externe Verbindungen verwenden `link_name`, das vom Frontend angegeben wird. Whale benennt eines der widersprüchlichen öffentlichen Symbole nicht automatisch um und behält seinen angegebenen Namen bei.

Die IR-Deklarationstabelle zeichnet typisierte `FunctionId`-Referenzen und explizite Funktionswerte `link_name` auf. Die folgenden Assembler-/Objekt-/Linker-APIs bleiben separate Schnittstellen: Die native IR-Emission überträgt diese Identitäten noch nicht bis zum endgültigen Link. Die IR-Aufrufüberprüfung allein stellt diese End-to-End-Eigenschaft nicht her.

Daher ist es möglich, dass sowohl eine interne Funktion als auch eine Variable den Namen `item` haben, es wäre jedoch ein Fehler, beide mit demselben externen Namen bereitzustellen. Durch die Trennung des internen Namensraums wird nicht automatisch auch der externe Namensraum getrennt.

Der Geltungsbereich eines lokalen Objektsymbols ist sein Eingabeobjekt. Globale Symbole sind an der Interpretation zwischen Objekten beteiligt. Bestätigte Funktions-/Datenkonflikte sind Fehler. Das Symbol NOTYPE gewährleistet die Kompatibilität mit Eingaben, die keinen spezifischeren Typ bereitstellen. Die bloße Tatsache, dass es keinen Typ hat, bedeutet nicht, dass es sich um eine Funktion oder Daten handelt.

## Definition auswählen

|Definition oder Referenz|Ergebnis|
| --- | --- |
|Strong und strong|Doppelter Definitionsfehler|
|Strong und weak|Strong Definition auswählen|
|Weak und weak|Wählen Sie die erste Definition in der Eingabereihenfolge aus|
|Siehe strong ungelöst|Linkfehler|
|Siehe ungelöstes weak|Fehler wird im anfänglichen statischen Profil nicht unterstützt|

Wenn mehrere weak-Definitionen vorhanden sind, wirkt sich die Reihenfolge ihrer Eingabe auf die Ergebnisse aus. Sie müssen die übergebene Eingabereihenfolge für deterministische Links konsistent verwenden.

## Statische ausführbare Ausgabe

Das statische Profil native generiert ELF ET_EXEC und gibt den Einstiegspunkt an. Der Einstiegspunkt wird nicht aus dem Funktionsnamen `main` abgeleitet. Außerdem wird der Startcode, der diese Funktion aufruft, nicht automatisch eingefügt.

Abschnitte werden nicht automatisch entfernt, identischer Code wird nicht zusammengeführt und Symbole werden nicht entfernt. Die Dateiplatzierung erfordert eine separate Abrechnung der tatsächlich gespeicherten Bytes und des zur Laufzeit reservierten Speichers.

Der vollständige Pfad zur Erstellung der statischen ausführbaren Datei für CLI ist noch nicht verfügbar. `whale asm` erstellt ein verschiebbares Objekt und `whale object` verpackt Rohbytes in ein Objekt. Informationen zur Verfügbarkeit finden Sie unter [Übersicht über die Toolchain](overview), zu den Anforderungen unter ABI und [AMD64 Ziel](amd64-target).


## Wählbare Wave Datensatzserialisierung

Linux Whale Builds für x86_64 können jetzt die Wave-Implementierung für feste ELF64 Header-Abschnittssymbol-RELA-Datensätze verwenden. Eine Standardimplementierung Rust wird ebenfalls bereitgestellt. Die Auswahl des Ausgabeziels, die Objektüberprüfung, die Platzierung, die Symbolinterpretation und die Pufferzuweisung werden von Rust übernommen. Durch Auswahl von Wave werden keine unterstützten Architekturen oder ein vollständiger Linker hinzugefügt.

Installieren Sie Rust, LLVM 21 Entwicklungsbibliotheken, C-Linker und `ar` und erstellen Sie dann den ausgewählten Pfad aus dem Whale-Repository.

```sh
git clone https://github.com/wavefnd/Wave.git /tmp/whale-wave-bootstrap
git -C /tmp/whale-wave-bootstrap checkout --detach 8a465e30aeea4b817d925cdd0e8d08c1bb029c9a
python3 tools/build_wave_elf.py --wave-source /tmp/whale-wave-bootstrap --out-dir /tmp/whale-wave-elf
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo build --locked --all-features
WHALE_WAVE_ELF_DIR=/tmp/whale-wave-elf cargo test --locked --workspace --all-features
```

Das Skript prüft, ob revision fixiert ist, lehnt die Änderung an tracked ab, erstellt den Wave-Compiler, erstellt ein Wave-Objekt mit LLVM und bündelt es mit archive für die statische Verknüpfung. `WHALE_WAVE_ELF_DIR` sollte dieses archive haben. Das explizite Anfordern eines ungültigen archive oder eines nicht unterstützten Hosts ist ein Build-Fehler. Normale Builds ohne angegebene Variablen erfordern weder die Compiler `--all-features` noch Wave. Der Wave-Compiler ist zur Laufzeit auch nicht für die verknüpfte ausführbare Datei Whale erforderlich. Der Pfad zum Cross-Compilieren von Whale selbst mit bootstrap wird noch nicht unterstützt.

Speichern Sie beispielsweise Folgendes als `return.asm` und fügen Sie es in die generierte Binärdatei ein:

```asm
section .text
global entry
entry:
    ret
```

```sh
target/debug/whale asm --amd64 return.asm -o return.o
```

```text
ELF class: ELF64
machine: AMD64 (62)
.text bytes: c3
```

Der Datensatz ABI trägt den Typ, den Feldzeiger u64, die Anzahl der Felder, den Ausgabezeiger und die Kapazität. Der Puffer gehört dem Aufrufer Rust. Der Zuweisungsbesitz oder der Ausdruck Rust enum/String/Vec wird nicht über Grenzen hinweg weitergegeben. Die Routine Wave prüft vor dem Schreiben Anzahl, Kapazität und Feldbreite. Erfolg gibt Status 0 zurück, ungültiger Typ/Zeiger/Kapazität gibt 1 zurück und Feld overflow gibt 2 zurück. Der Zeiger muss auf aktive Puffer zeigen, die die richtige Größe haben und sich nicht überlappen. Der rohe Zeiger C allein kann diese Bedingung nicht beweisen; wrapper erfüllt es. Der Test vergleicht den gesamten ELF, einschließlich BSS und signed Verschiebung addend, mit dem Rust Pfad. Dies ist eine teilweise Wave-Implementierung, die Wave/LLVM bootstrap verwendet und nicht vollständig selbst gehostet ist.
