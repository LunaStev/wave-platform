package web

import (
	"encoding/xml"
	documentdomain "github.com/wavefnd/wave-platform/internal/document"
	"github.com/wavefnd/wave-platform/internal/storage"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStandardLibraryRoutesFallbackAndSitemap(t *testing.T) {
	db, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	repository := documentdomain.NewRepository(db)
	seo := NewSEOHandler("https://wave.example", repository, nil, nil, nil, nil)
	server := seoFrontend(t, seo)
	get := func(path string, status int) string {
		t.Helper()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != status {
			t.Fatalf("%s status=%d", path, response.Code)
		}
		return response.Body.String()
	}
	if !strings.Contains(get("/docs/ko/stdlib", 200), `href="/docs/ko/stdlib" aria-current="page"`) {
		t.Fatal("library catalogue not selected")
	}
	for _, item := range []struct{ locale, path, title string }{
		{"en", "reference/standard-library", "Library overview"},
		{"ko", "reference/standard-library", "표준 라이브러리 안내"},
		{"en", "stdlib/buffer", "Buffer"},
		{"ko", "stdlib/files", "파일"},
		{"en", "language/types", "Wave types"},
		{"en", "whale/overview", "Whale overview"},
	} {
		id := item.locale + "/" + item.path
		content, err := xml.Marshal(documentdomain.Content{Markdown: "## Example\n\n" + item.title})
		if err != nil {
			t.Fatal(err)
		}
		if err := repository.PutRevision(documentdomain.Revision{ID: "r1", DocumentID: id, ContentXML: content}); err != nil {
			t.Fatal(err)
		}
		if err := repository.UpsertDocument(documentdomain.Document{ID: id, TranslationSetID: item.path, Locale: item.locale, Path: item.path, Group: "reference", GroupOrder: 1, Order: 1, Title: item.title, Status: "published", PublishedRevisionID: "r1"}); err != nil {
			t.Fatal(err)
		}
	}
	catalog := seo.metadata("/docs/ko/stdlib", "https://wave.example")
	if len(catalog.Items) != 3 || catalog.SchemaType != "CollectionPage" {
		t.Fatalf("catalog=%#v", catalog)
	}
	for _, item := range catalog.Items {
		if strings.Contains(item.URL, "/language/") || strings.Contains(item.URL, "/whale/") {
			t.Fatalf("cross-tab item: %#v", item)
		}
	}
	for _, alt := range catalog.Alternates {
		if !strings.HasSuffix(alt.URL, "/stdlib") {
			t.Fatalf("cross-tab alternate: %#v", alt)
		}
	}
	if !strings.Contains(get("/docs/ko/reference/standard-library", 200), `href="/docs/ko/stdlib" aria-current="page"`) {
		t.Fatal("legacy URL lost library tab")
	}
	fallback := get("/docs/ko/stdlib/buffer", 200)
	for _, want := range []string{`<html lang="en">`, `rel="canonical" href="https://wave.example/docs/en/stdlib/buffer"`, `href="/docs/ko/stdlib" aria-current="page"`} {
		if !strings.Contains(fallback, want) {
			t.Fatalf("fallback missing %s", want)
		}
	}
	get("/docs/ko/stdlib/missing", 404)
	if got := seo.metadata("/docs/ko", "https://wave.example"); len(got.Items) != 1 {
		t.Fatalf("library leaked into Wave catalogue: %#v", got.Items)
	}
	response := httptest.NewRecorder()
	seo.Sitemap(response, httptest.NewRequest("GET", "/sitemap.xml", nil))
	var set sitemapURLSet
	if err := xml.Unmarshal(response.Body.Bytes(), &set); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, item := range set.URLs {
		if seen[item.Location] {
			t.Fatalf("duplicate %s", item.Location)
		}
		seen[item.Location] = true
	}
	for _, path := range []string{"/docs/ko/stdlib", "/docs/en/stdlib", "/docs/ko/reference/standard-library", "/docs/en/stdlib/buffer"} {
		if !seen["https://wave.example"+path] {
			t.Fatalf("missing %s", path)
		}
	}
	if seen["https://wave.example/docs/ko/stdlib/buffer"] {
		t.Fatal("fallback advertised as translation")
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	get("/docs/ko/stdlib", 503)
}
