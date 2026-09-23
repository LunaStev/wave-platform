package web

import (
	"encoding/xml"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	documentdomain "github.com/wavefnd/wave-platform/internal/document"
	"github.com/wavefnd/wave-platform/internal/storage"
	"github.com/wavefnd/wave-platform/internal/web/handler"
)

func TestWhaleDocumentationRoutesAndSEO(t *testing.T) {
	db, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	repository := documentdomain.NewRepository(db)
	seo := NewSEOHandler("https://wave.example", repository, nil, nil, nil, nil)
	server := seoFrontend(t, seo)
	publish := func(locale, path, title string) {
		t.Helper()
		id := locale + "/" + path
		content, err := xml.Marshal(documentdomain.Content{Markdown: "## Fixture\n\n" + title})
		if err != nil {
			t.Fatal(err)
		}
		if err := repository.PutRevision(documentdomain.Revision{ID: "r1", DocumentID: id, ContentXML: content}); err != nil {
			t.Fatal(err)
		}
		if err := repository.UpsertDocument(documentdomain.Document{ID: id, TranslationSetID: path, Locale: locale, Path: path, Group: "reference", GroupOrder: 5, Order: 1, Title: title, Summary: "Test fixture", Status: "published", PublishedRevisionID: "r1", UpdatedAt: time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)}); err != nil {
			t.Fatal(err)
		}
	}
	get := func(path string, status int) string {
		t.Helper()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != status {
			t.Fatalf("%s status=%d want=%d", path, response.Code, status)
		}
		return response.Body.String()
	}
	publish("en", "getting-started/overview", "Wave overview fixture")
	publish("en", "toolchain/whale-legacy", "Wave toolchain fixture")
	// A real empty catalogue, not a missing document or a copy of Wave's list.
	empty := get("/docs/ko/whale", 200)
	if !strings.Contains(empty, "No documents have been published") || !strings.Contains(empty, `href="/docs/ko/whale" aria-current="page"`) {
		t.Fatal("missing empty state or selected Whale tab")
	}
	if got := seo.metadata("/docs/ko/whale", "https://wave.example"); len(got.Items) != 0 || got.SchemaType != "CollectionPage" {
		t.Fatalf("empty catalogue metadata=%#v", got)
	}
	get("/docs/ko/whale/overview", 404)

	publish("en", "whale/overview", "Whale overview fixture")
	publish("ko", "whale/overview", "Whale 소개 테스트")
	publish("en", "whale/symbols", "Whale symbols fixture")
	for _, locale := range []string{"en", "ko", "ja"} {
		for _, project := range []string{"wave", "whale"} {
			path := "/docs/" + locale
			if project == "whale" {
				path += "/whale"
			}
			get(path, 200)
			metadata := seo.metadata(path, "https://wave.example")
			if len(metadata.Items) != 2 {
				t.Fatalf("%s items=%#v", path, metadata.Items)
			}
			for _, item := range metadata.Items {
				if strings.Contains(item.URL, "/whale/") != (project == "whale") {
					t.Fatalf("cross-project item at %s: %#v", path, item)
				}
			}
			for _, alt := range metadata.Alternates {
				if strings.HasSuffix(alt.URL, "/whale") != (project == "whale") {
					t.Fatalf("cross-project alternate: %#v", alt)
				}
			}
		}
	}
	translated := get("/docs/ko/whale/overview", 200)
	if !strings.Contains(translated, `<html lang="ko">`) || !strings.Contains(translated, "Whale 소개 테스트") {
		t.Fatal("translated page missing")
	}
	fallback := get("/docs/ko/whale/symbols", 200)
	for _, expected := range []string{`<html lang="en">`, `rel="canonical" href="https://wave.example/docs/en/whale/symbols"`, `href="/docs/ko/whale" aria-current="page"`, `Whale symbols fixture · Whale Documentation`} {
		if !strings.Contains(fallback, expected) {
			t.Fatalf("fallback missing %q", expected)
		}
	}
	if strings.Contains(fallback, `hreflang="ko"`) {
		t.Fatal("fallback advertised as translated")
	}
	get("/docs/ko/whale/missing", 404)
	get("/docs/en/toolchain/whale-legacy", 200)
	get("/docs/en/getting-started/overview", 200)

	// Existing wildcard API discovers full Whale paths, including metadata order.
	api := handler.DocumentsHandler{Repository: repository}
	request := httptest.NewRequest("GET", "/api/v1/documents/whale/overview?locale=ko", nil)
	request.SetPathValue("path", "whale/overview")
	response := httptest.NewRecorder()
	api.Get(response, request)
	var view documentdomain.View
	if err := xml.Unmarshal(response.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if response.Code != 200 || view.Path != "whale/overview" || view.GroupOrder != 5 || len(view.Translations) != 2 {
		t.Fatalf("API view=%#v status=%d", view, response.Code)
	}

	response = httptest.NewRecorder()
	seo.Sitemap(response, httptest.NewRequest("GET", "/sitemap.xml", nil))
	var set sitemapURLSet
	if err := xml.Unmarshal(response.Body.Bytes(), &set); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, item := range set.URLs {
		if seen[item.Location] {
			t.Fatalf("duplicate sitemap URL %s", item.Location)
		}
		seen[item.Location] = true
	}
	for _, path := range []string{"/docs/en/whale", "/docs/ko/whale", "/docs/en/whale/overview", "/docs/ko/whale/overview", "/docs/en/whale/symbols"} {
		if !seen["https://wave.example"+path] {
			t.Fatalf("missing sitemap URL %s", path)
		}
	}
	for _, path := range []string{"/docs/ja/whale", "/docs/ko/whale/symbols", "/docs/ja/whale/overview"} {
		if seen["https://wave.example"+path] {
			t.Fatalf("untranslated sitemap URL %s", path)
		}
	}
}

func TestWhaleCatalogueDoesNotHideStorageFailures(t *testing.T) {
	db, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repository := documentdomain.NewRepository(db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	server := seoFrontend(t, NewSEOHandler("https://wave.example", repository, nil, nil, nil, nil))
	for _, path := range []string{"/docs/ko/whale", "/docs/ko/whale/overview"} {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != 503 {
			t.Fatalf("%s status=%d", path, response.Code)
		}
	}
}
