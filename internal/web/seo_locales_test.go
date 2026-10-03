package web

import (
	"encoding/xml"
	"net/http/httptest"
	"strings"
	"testing"

	documentdomain "github.com/wavefnd/wave-platform/internal/document"
	"github.com/wavefnd/wave-platform/internal/storage"
)

func TestIncrementalLocalesUseEnglishUntilTranslationIsPublished(t *testing.T) {
	db, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	repository := documentdomain.NewRepository(db)
	publish := func(locale, path string) {
		t.Helper()
		id := locale + "/" + path
		content, err := xml.Marshal(documentdomain.Content{Markdown: "## Fixture\n\n" + locale})
		if err != nil {
			t.Fatal(err)
		}
		if err := repository.PutRevision(documentdomain.Revision{ID: "r1", DocumentID: id, ContentXML: content}); err != nil {
			t.Fatal(err)
		}
		if err := repository.UpsertDocument(documentdomain.Document{ID: id, TranslationSetID: path, Locale: locale, Path: path, Title: locale + " fixture", Summary: "Fixture", Group: "language", Status: "published", PublishedRevisionID: "r1"}); err != nil {
			t.Fatal(err)
		}
	}
	paths := []string{"language/fixture", "stdlib/fixture", "whale/fixture"}
	for _, path := range paths {
		publish("en", path)
	}
	seo := NewSEOHandler("https://wave.example", repository, nil, nil, nil, nil)
	server := seoFrontend(t, seo)
	assertSitemap := func(wantTranslated bool) {
		t.Helper()
		response := httptest.NewRecorder()
		seo.Sitemap(response, httptest.NewRequest("GET", "/sitemap.xml", nil))
		if response.Code != 200 {
			t.Fatalf("sitemap status=%d", response.Code)
		}
		for _, locale := range []string{"pt", "fr", "pl", "nl", "tr", "it"} {
			for _, path := range paths {
				present := strings.Contains(response.Body.String(), "https://wave.example/docs/"+locale+"/"+path)
				if present != wantTranslated {
					t.Fatalf("sitemap %s/%s present=%v want=%v", locale, path, present, wantTranslated)
				}
			}
		}
	}
	assertSitemap(false)
	for _, locale := range []string{"pt", "fr", "pl", "nl", "tr", "it"} {
		for _, path := range paths {
			requestPath := "/docs/" + locale + "/" + path
			response := httptest.NewRecorder()
			server.ServeHTTP(response, httptest.NewRequest("GET", requestPath, nil))
			if response.Code != 200 || response.Header().Get("Content-Language") != "en" {
				t.Fatalf("%s: status=%d language=%s", requestPath, response.Code, response.Header().Get("Content-Language"))
			}
			metadata := seo.metadata(requestPath, "https://wave.example")
			if metadata.Canonical != "https://wave.example/docs/en/"+path {
				t.Fatalf("incorrect fallback canonical: %#v", metadata)
			}
			for _, alternate := range metadata.Alternates {
				if alternate.Language == locale {
					t.Fatalf("untranslated alternate: %#v", alternate)
				}
			}
			catalog := "/docs/" + locale
			project := documentdomain.ProjectForPath(path)
			if project != "wave" {
				catalog += "/" + project
			}
			response = httptest.NewRecorder()
			server.ServeHTTP(response, httptest.NewRequest("GET", catalog, nil))
			if response.Code != 200 || len(seo.metadata(catalog, "https://wave.example").Items) != 1 {
				t.Fatalf("missing English catalog for %s", catalog)
			}
			publish(locale, path)
			response = httptest.NewRecorder()
			server.ServeHTTP(response, httptest.NewRequest("GET", requestPath, nil))
			if response.Code != 200 || response.Header().Get("Content-Language") != locale {
				t.Fatalf("translation did not replace fallback: %s", requestPath)
			}
			if !strings.Contains(response.Body.String(), `rel="canonical" href="https://wave.example`+requestPath+`"`) {
				t.Fatalf("translated canonical missing: %s", requestPath)
			}
		}
	}
	assertSitemap(true)
}
