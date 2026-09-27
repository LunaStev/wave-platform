package document

import (
	"errors"
	"strings"
	"testing"

	"github.com/wavefnd/wave-platform/internal/storage"
	"github.com/wavefnd/wave-platform/wavedoc"
)

func TestRestructureArchivesRetiredOfficialPages(t *testing.T) {
	db, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	repository := NewRepository(db)
	for path := range wavedoc.DocumentRedirects {
		id := "official/ko/" + path
		if err := repository.UpsertDocument(Document{ID: id, Locale: "ko", Path: path, Title: "Old", Status: "published"}); err != nil {
			t.Fatal(err)
		}
	}
	if err := repository.UpsertDocument(Document{ID: "community/guide", Locale: "ko", Path: "custom/guide", Title: "Custom", Status: "published"}); err != nil {
		t.Fatal(err)
	}
	if _, err := SeedOfficial(db); err != nil {
		t.Fatal(err)
	}
	for old, target := range wavedoc.DocumentRedirects {
		if _, err := repository.Published("ko", old); !errors.Is(err, storage.ErrNotFound) {
			t.Fatalf("old page %s remains published: %v", old, err)
		}
		if _, err := repository.Published("ko", target); err != nil {
			t.Fatalf("missing destination %s: %v", target, err)
		}
	}
	custom, err := repository.Document("community/guide")
	if err != nil || custom.Status != "published" {
		t.Fatalf("custom content changed: %+v %v", custom, err)
	}
	wave, err := repository.Navigation("ko", "wave")
	if err != nil {
		t.Fatal(err)
	}
	for _, page := range wave {
		if page.Group == "learn" || page.Group == "toolchain" || strings.Contains(page.Path, "design-goals") {
			t.Fatalf("retired navigation remains: %+v", page)
		}
	}
	whale, err := repository.Navigation("ko", "whale")
	if err != nil {
		t.Fatal(err)
	}
	found := map[string]bool{}
	for _, page := range whale {
		found[page.Path] = true
	}
	for _, path := range []string{"whale/build-link-targets", "whale/ecosystem", "whale/vex-package-manager", "whale/whale-cli"} {
		if !found[path] {
			t.Errorf("missing moved toolchain page %s", path)
		}
	}
	if count, err := SeedOfficial(db); err != nil || count != 0 {
		t.Fatalf("second seed=%d %v", count, err)
	}
}
