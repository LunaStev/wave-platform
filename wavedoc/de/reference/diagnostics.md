---
translation_set_id: diagnostics
path: reference/diagnostics
locale: de
group: reference
group_order: 5
order: 2
title: Fehlerbehebung: Von der Installation bis zur Ausführung
summary: Isolieren Sie fehlgeschlagene Schritte und grenzen Sie die Ursache mithilfe reproduzierbarer Informationen ein.
---

## Unterscheiden Sie zunächst die Stadien des Scheiterns

|beobachtetes Phänomen|Überprüfen Sie zuerst|nächste Aktion|
| --- | --- | --- |
|wavec Befehl nicht gefunden|PATH und Speicherort der ausführbaren Datei|Mit absolutem Pfad ausführen und PATH einstellen|
|Die zum Ausführen benötigten Dateien können nicht gefunden werden|Fehlen Dateien im Installationsordner?|Entpacken und installieren Sie das gesamte Paket erneut|
|std import Fehlgeschlagen| `wavec print std-path` |Korrespondenz: std installieren oder `--std-root` angeben|
|Quellort und Typfehlerausgabe| `wavec check main.wave` |Beheben Sie den ersten Fehler und überprüfen Sie ihn erneut|
|Build-Fehler für andere OS·CPU-Ziele|Angegebener target und Zielumgebung|[Cross-Build-Einstellungen](/docs/de/whale/build-link-targets) Bestätigen|
|Die Ausführung schlägt nach erfolgreichem Build fehl|Exit-Code, Eingabe, Arbeitsverzeichnis|Ausführungsumgebung und API Fehlerprüfung|

## Ein kleines Diagnosebeispiel

Hier ist das gesamte Programm, das absichtlich falsch ist:

```wave
fun main() {
    var count: i32 = 1;
    println("{}", missing);
}
```

`wavec check main.wave` muss auf den nicht deklarierten Namen missing verweisen. Ändern Sie den Variablennamen in count, überprüfen Sie ihn und führen Sie ihn erneut aus. Konzentrieren Sie sich auf die Datei, den Speicherort und die Ursache und nicht auf den gesamten Diagnosetext. Nachfolgende Fehler können auf den anfänglichen Fehler zurückzuführen sein.

## Wenn eine ausführbare Datei fehlschlägt

In der Linux/macOS-Shell wird der Exit-Code unmittelbar nach der Ausführung als `echo $?` überprüft, und in PowerShell ist er `$LASTEXITCODE`. Eingabefehler und explizites `return 1` sind nicht die gleiche Ursache. Ungültige Laufzeitwerte für Schichtanzahl oder Realkonvertierung können trap verursachen. Schauen Sie sich [Betriebsregeln](/docs/de/language/expressions-and-operators) an.

Relative Dateipfade werden vom Arbeitsverzeichnis der ausführbaren Datei und nicht vom Speicherort der Quelldatei beeinflusst. Behandeln Sie einen Dateilesefehler nicht als Zeichenfolgenlänge 0, sondern überprüfen Sie zuerst den zurückgegebenen Fehler. Netzwerkverbindungsfehler werden durch Adresssuche, Serverwartezeiten, Berechtigungen und Zeitüberschreitungen überprüft.

## Informationen, die zum Melden eines Problems erforderlich sind

1. `wavec --version` Ausgabe und genauer Befehl ausgeführt.
2. target. wird separat von der Host-OS·Architektur angegeben
3. Compiler-Quelle, die mit dem ausgewählten std-Pfad verwendet wird.
4. Minimale Quell-, Eingabe- und erforderliche Dateien zur Reproduktion des Problems.
5. Erwartete Ergebnisse, tatsächliche Ergebnisse, Diagnose- und Exit-Codes.

Passwörter, Token und persönliche Dateiinhalte werden entfernt. Wenn das Problem verschwindet, wenn Sie das Minimalbeispiel reduzieren, ist das letzte entfernte Element der Hinweis. `--error-format=json` ist verfügbar, wenn das Tool Diagnosedaten erfasst.

[Installation](/docs/de/getting-started/install) · [Compiler-Befehl](/docs/de/getting-started/compiler) · [Ziele und Links](/docs/de/whale/build-link-targets)
