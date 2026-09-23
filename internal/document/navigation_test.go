package document

import (
	"bytes"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/wavefnd/wave-platform/internal/storage"
)

func TestNavigationScopesAndMergesPublishedDocuments(t *testing.T) {
	db, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	repository := NewRepository(db)
	for _, item := range []Document{
		{Locale: "en", Path: "language/types", Title: "Wave", GroupOrder: 1, Order: 1},
		{Locale: "en", Path: "toolchain/whale-overview", Title: "Legacy", GroupOrder: 4, Order: 1},
		{Locale: "en", Path: "whale/symbols", Title: "Symbols", GroupOrder: 5, Order: 4},
		{Locale: "ko", Path: "whale/symbols", Title: "심볼", GroupOrder: 5, Order: 4},
		{Locale: "en", Path: "whale/alignment", Title: "Alignment", GroupOrder: 5, Order: 4},
		{Locale: "en", Path: "whale/overview", Title: "Overview", GroupOrder: 5, Order: 1},
		{Locale: "en", Path: "whale/other", Title: "Other", GroupOrder: 6, Order: 1},
		{Locale: "en", Path: "whale/draft", Title: "Private", Status: "draft"},
	} {
		item.ID = item.Locale + "/" + item.Path
		if item.Status == "" {
			item.Status = "published"
		}
		// Group/title deliberately do not decide project membership.
		item.Group = "reference"
		if err := repository.UpsertDocument(item); err != nil {
			t.Fatal(err)
		}
	}
	whale, err := repository.Navigation("ko", "whale")
	if err != nil {
		t.Fatal(err)
	}
	paths := make([]string, 0, len(whale))
	for _, item := range whale {
		paths = append(paths, item.Path)
	}
	if want := []string{"whale/overview", "whale/alignment", "whale/symbols", "whale/other"}; !reflect.DeepEqual(paths, want) {
		t.Fatalf("paths=%v", paths)
	}
	if whale[2].Title != "심볼" || whale[2].Locale != "ko" || whale[0].Locale != "en" {
		t.Fatalf("incorrect translation merge: %#v", whale)
	}
	wave, err := repository.Navigation("ko", "wave")
	if err != nil || len(wave) != 2 || wave[1].Path != "toolchain/whale-overview" {
		t.Fatalf("Wave navigation=%v err=%v", wave, err)
	}
	for path, want := range map[string]string{"whale/overview": "whale", "whales/overview": "wave", "toolchain/whale-overview": "wave"} {
		if got := ProjectForPath(path); got != want {
			t.Fatalf("project(%q)=%s", path, got)
		}
	}
}

func TestMarkdownDiscoveryAcceptsNewWhaleFilesWithoutRegistry(t *testing.T) {
	source := fstest.MapFS{"en/whale/overview.md": &fstest.MapFile{Data: []byte("---\ntranslation_set_id: whale-overview\npath: whale/overview\nlocale: en\ngroup: whale\ngroup_order: 5\norder: 1\ntitle: Overview\nsummary: Fixture only.\n---\n\n## Overview\n\nFixture body.\n")}}
	first, digest, err := readOfficialDocumentsFrom(source)
	if err != nil || len(first) != 1 {
		t.Fatalf("discovery=%v err=%v", first, err)
	}
	source["ko/whale/symbols.md"] = &fstest.MapFile{Data: []byte("---\ntranslation_set_id: whale-symbols\npath: whale/symbols\nlocale: ko\ngroup: whale\ngroup_order: 5\norder: 4\ntitle: 심볼\nsummary: 테스트 문서.\n---\n\n## 심볼\n\n본문.\n")}
	second, nextDigest, err := readOfficialDocumentsFrom(source)
	if err != nil || len(second) != 2 {
		t.Fatalf("discovery=%v err=%v", second, err)
	}
	if bytes.Equal(digest, nextDigest) {
		t.Fatal("new Markdown file must invalidate the startup import marker")
	}
	if second[1].Path != "whale/symbols" || second[1].Locale != "ko" || second[1].GroupOrder != 5 || second[1].Order != 4 {
		t.Fatalf("metadata lost: %#v", second[1])
	}
}
