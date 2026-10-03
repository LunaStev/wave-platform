package document

import (
	"reflect"
	"strings"
	"testing"

	"github.com/wavefnd/wave-platform/internal/storage"
	"github.com/wavefnd/wave-platform/wavedoc"
)

func TestOfficialDocumentationUsesStableLanguageVocabulary(t *testing.T) {
	documents, _, err := readOfficialDocuments()
	if err != nil {
		t.Fatal(err)
	}

	forbidden := []string{
		"`let`",
		"`let mut`",
		"current compiler contract",
		"current implementation",
		"this release",
		"removed syntax",
		"code generation treats",
		"lexer recognizes",
		"type-conversion path",
		"현재 컴파일러 계약",
		"현재 구현",
		"이 릴리스",
		"폐지된 문법",
		"코드 생성 단계",
		"타입 변환 경로",
	}
	for _, document := range documents {
		for _, phrase := range forbidden {
			if strings.Contains(document.Markdown, phrase) {
				t.Errorf("%s/%s contains implementation-history wording %q", document.Locale, document.Path, phrase)
			}
		}
	}

	// Assert language syntax without pinning natural-language translations to
	// a particular sentence or punctuation convention.
	requiredByDocument := map[string][]string{
		"language/declarations-and-types":     {"`var`", "`i8`", "`i16`", "`i32`", "`i64`", "`i128`", "`isz`", "`usz`"},
		"language/explicit-memory-type-model": {"**Wave Explicit Memory Type Model**", "`null`", "ptr<i32>", "deref value = deref value + 1;"},
		"language/console-io-and-formatting":  {"println", "input"},
		"language/modules-imports-and-ffi":    {"`./`", "pub fun"},
	}
	for _, document := range documents {
		for _, syntax := range requiredByDocument[document.Path] {
			if !strings.Contains(document.Markdown, syntax) {
				t.Errorf("%s/%s is missing language syntax %q", document.Locale, document.Path, syntax)
			}
		}
	}
}

func TestSeedOfficialPublishesSupportedDocumentationLocales(t *testing.T) {
	database, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })

	sources, _, err := readOfficialDocuments()
	if err != nil {
		t.Fatal(err)
	}
	expectedCounts := make(map[string]int)
	for _, source := range sources {
		expectedCounts[source.Locale]++
	}
	count, err := SeedOfficial(database)
	if err != nil || count != len(sources) {
		t.Fatalf("seed count=%d err=%v", count, err)
	}
	repository := NewRepository(database)
	for _, locale := range wavedoc.SupportedLocales {
		if wavedoc.RequiresCompleteCoverage(locale) && (expectedCounts[locale] == 0 || expectedCounts[locale] != expectedCounts["ko"]) {
			t.Fatalf("%s has %d documents; Korean has %d", locale, expectedCounts[locale], expectedCounts["ko"])
		}
		items, err := repository.Summaries(locale)
		if err != nil || len(items) != expectedCounts[locale] {
			t.Fatalf("%s summaries=%d err=%v", locale, len(items), err)
		}
	}
	for _, source := range sources {
		view, err := repository.Published(source.Locale, source.Path)
		if err != nil {
			t.Fatalf("%s/%s: %v", source.Locale, source.Path, err)
		}
		if view.Title != source.Title || view.Markdown != source.Markdown {
			t.Errorf("%s/%s did not preserve the translated authoring source", source.Locale, source.Path)
		}
	}
	install, err := repository.Published("en", "getting-started/install")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(install.Markdown, "curl -fsSL https://wave-lang.dev/install.sh | bash -s -- latest") {
		t.Fatal("official installer command is missing")
	}
	count, err = SeedOfficial(database)
	if err != nil || count != 0 {
		t.Fatalf("second seed count=%d err=%v", count, err)
	}
}

func TestOfficialInstallDocumentsIncludeWindowsInstaller(t *testing.T) {
	database, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })
	if _, err := SeedOfficial(database); err != nil {
		t.Fatal(err)
	}
	repository := NewRepository(database)
	for _, translation := range repository.PublishedTranslations("getting-started/install") {
		locale := translation.Locale
		install, err := repository.Published(locale, "getting-started/install")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(install.Markdown, "https://wave-lang.dev/install.ps1") || !strings.Contains(install.Markdown, "-Latest") {
			t.Fatalf("%s official Windows installer command is missing", locale)
		}
		if !strings.Contains(install.Markdown, "vex --version") {
			t.Fatalf("%s Vex installation guidance is missing", locale)
		}
	}
}

func TestMarkdownHeadingsExcludeFencedCode(t *testing.T) {
	cases := []struct {
		name string
		code string
	}{
		{"backticks", "```wave\n## 소개\n```"},
		{"tildes", "~~~wave\n## 소개\n~~~"},
		{"shorter backticks", "````wave\n```\n## 소개\n````"},
		{"shorter tildes", "~~~~wave\n~~~\n## 소개\n~~~~"},
		{"tildes inside backticks", "```wave\n~~~\n## 소개\n```"},
		{"backticks inside tildes", "~~~wave\n```\n## 소개\n~~~"},
		{"backtick closing text", "```wave\n```text\n## 소개\n```"},
		{"tilde closing text", "~~~wave\n~~~text\n## 소개\n~~~"},
		{"tab after closing fence", "~~~wave\n~~~\t\n## 소개\n~~~"},
		{"longer closing fence", "~~~wave\n## 소개\n~~~~~"},
		{"indented fence", "   ~~~wave\n## 소개\n  ~~~  "},
	}
	want := []Block{
		{Kind: "heading", Anchor: "소개", Level: 2, Text: "소개"},
		{Kind: "heading", Anchor: "소개-1", Level: 3, Text: "소개"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			markdown := "## 소개\n\n" + tc.code + "\n\n### 소개\n"
			if got := markdownHeadings(markdown); !reflect.DeepEqual(got, want) {
				t.Fatalf("markdownHeadings() = %#v, want %#v", got, want)
			}
		})
	}
}
