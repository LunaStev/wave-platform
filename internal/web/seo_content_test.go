package web

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	blogdomain "github.com/wavefnd/wave-platform/internal/blog"
	communitydomain "github.com/wavefnd/wave-platform/internal/community"
	documentdomain "github.com/wavefnd/wave-platform/internal/document"
	questiondomain "github.com/wavefnd/wave-platform/internal/question"
	"github.com/wavefnd/wave-platform/internal/storage"
)

func seoFrontend(t *testing.T, seo SEOHandler) http.Handler {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(`<!doctype html><html lang="en"><head><!-- wave:seo:start --><!-- wave:seo:end --></head><body><div id="app"></div></body></html>`), 0600); err != nil {
		t.Fatal(err)
	}
	return frontendHandler(dir, seo)
}

func TestPublicHTMLContainsArticleWithoutJavaScript(t *testing.T) {
	service, closeDB := testSEOService(t)
	defer closeDB()
	post := blogdomain.Post{Slug: "rendered", Category: "article", Title: "Compiler article", Summary: "Compiler architecture", Content: "## Parsing\n\nWave parses **source** into a tree.\n\n```wave\nvar n: i32 = 1;\n```\n\n[Documentation](/docs/en)\n\n<script>alert('injected')</script>\n\n[unsafe](javascript:alert(1))", Status: "published", UpdatedAt: time.Now()}
	if err := service.Repository().Upsert(post); err != nil {
		t.Fatal(err)
	}
	post.Slug, post.Status, post.Content = "private", "draft", "SECRET DRAFT TEXT"
	if err := service.Repository().Upsert(post); err != nil {
		t.Fatal(err)
	}
	seo := NewSEOHandler("https://wave.example", nil, service, nil, nil, nil)
	server := seoFrontend(t, seo)
	var first string
	for _, agent := range []string{"Mozilla/5.0", "Googlebot", "bingbot", "DuckDuckBot", "YandexBot"} {
		request := httptest.NewRequest("GET", "/blog/rendered", nil)
		request.Header.Set("User-Agent", agent)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		body := response.Body.String()
		if response.Code != 200 {
			t.Fatalf("status %d", response.Code)
		}
		for _, expected := range []string{`<h1>Compiler article</h1>`, `<strong>source</strong>`, `var n: i32 = 1;`, `href="/docs/en"`, `id="wave-page-data"`} {
			if !strings.Contains(body, expected) {
				t.Fatalf("missing %q", expected)
			}
		}
		if strings.Contains(body, `<script>alert`) || strings.Contains(body, `href="javascript:`) || strings.Contains(body, "SECRET DRAFT TEXT") {
			t.Fatal("unsafe or private content leaked")
		}
		if first != "" && first != body {
			t.Fatal("response varies by crawler")
		}
		first = body
	}
	for _, path := range []string{"/blog/private", "/blog/missing", "/releases/rendered"} {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != 404 || response.Header().Get("X-Robots-Tag") != "noindex, nofollow, noarchive" || strings.Contains(response.Body.String(), "SECRET DRAFT TEXT") || strings.Contains(response.Body.String(), `id="wave-page-data"`) {
			t.Fatalf("invalid missing page: %s, %d", path, response.Code)
		}
	}
	dataHTML := seo.pageData("/blog/rendered")
	encoded := strings.TrimSuffix(strings.SplitN(dataHTML, ">", 2)[1], "</script>")
	var data map[string]string
	if err := json.Unmarshal([]byte(encoded), &data); err != nil {
		t.Fatal(err)
	}
	var decoded blogdomain.Post
	if err := xml.Unmarshal([]byte(data["/api/v1/blog/posts/rendered"]), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Language != "en" || decoded.Content == "" {
		t.Fatalf("bootstrap content missing: %#v", decoded)
	}
}

func TestTranslatedSitemapAndFallbackHTML(t *testing.T) {
	db, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := documentdomain.SeedOfficial(db); err != nil {
		t.Fatal(err)
	}
	seo := NewSEOHandler("https://wave.example", documentdomain.NewRepository(db), nil, communitydomain.NewRepository(db), questiondomain.NewRepository(db), nil)
	response := httptest.NewRecorder()
	seo.Sitemap(response, httptest.NewRequest("GET", "/sitemap.xml", nil))
	var set sitemapURLSet
	if err := xml.Unmarshal(response.Body.Bytes(), &set); err != nil {
		t.Fatal(err)
	}
	path := "/language/explicit-memory-type-model"
	seen := map[string]bool{}
	for _, item := range set.URLs {
		if seen[item.Location] {
			t.Fatalf("duplicate URL %s", item.Location)
		}
		seen[item.Location] = true
		if strings.HasSuffix(item.Location, path) {
			if item.LastModified == "" || len(item.Alternates) != 4 {
				t.Fatalf("alternates or lastmod missing: %#v", item)
			}
			for _, alt := range item.Alternates {
				if alt.Language != "en" && alt.Language != "ko" && alt.Language != "ja" && alt.Language != "x-default" {
					t.Fatalf("untranslated locale advertised: %#v", alt)
				}
			}
		}
	}
	for _, locale := range []string{"en", "ko", "ja"} {
		if !seen["https://wave.example/docs/"+locale+path] {
			t.Fatalf("missing %s", locale)
		}
	}
	if seen["https://wave.example/docs/zh"+path] {
		t.Fatal("fallback URL in sitemap")
	}
	server := seoFrontend(t, seo)
	for _, test := range []struct{ path, language, og string }{{"/docs/ja" + path, "ja", "ja_JP"}, {"/docs/zh" + path, "en", "en_US"}} {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", test.path, nil))
		body := response.Body.String()
		if response.Code != 200 || response.Header().Get("Content-Language") != test.language {
			t.Fatalf("language/status %s: %d %s", test.path, response.Code, response.Header().Get("Content-Language"))
		}
		for _, expected := range []string{`<html lang="` + test.language + `">`, `content="` + test.og + `"`, `<strong>Wave Explicit Memory Type Model</strong>`, `hreflang="ja"`, `hreflang="x-default"`} {
			if !strings.Contains(body, expected) {
				t.Fatalf("missing %q on %s", expected, test.path)
			}
		}
		if strings.Contains(body, `hreflang="zh"`) {
			t.Fatal("fallback advertised as translated")
		}
	}
	for _, path := range []string{"/docs/en/missing", "/docs/ja/new", "/community/thread/missing", "/questions/missing", "/no-such-page"} {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != 404 || strings.Contains(response.Body.String(), "application/ld+json") {
			t.Fatalf("soft 404 at %s: %d", path, response.Code)
		}
	}
	for _, path := range []string{"/community/new", "/questions/new"} {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != 200 || !strings.Contains(response.Header().Get("X-Robots-Tag"), "noindex") {
			t.Fatalf("creation route %s: %d", path, response.Code)
		}
	}
}

func TestStorageFailureDoesNotLookLikeDeletedContent(t *testing.T) {
	db, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repository := documentdomain.NewRepository(db)
	db.Close()
	seo := NewSEOHandler("https://wave.example", repository, nil, nil, nil, nil)
	response := httptest.NewRecorder()
	seoFrontend(t, seo).ServeHTTP(response, httptest.NewRequest("GET", "/docs/en/missing", nil))
	if response.Code != 503 || response.Header().Get("Retry-After") != "60" {
		t.Fatalf("status %d", response.Code)
	}
}

func TestSitemapSplitsLargeCollections(t *testing.T) {
	entries := make([]sitemapURL, sitemapPageSize+1)
	for i := range entries {
		entries[i].Location = "https://wave.example/blog/" + strconv.Itoa(i)
	}
	index := httptest.NewRecorder()
	writeSitemap(index, httptest.NewRequest("GET", "/sitemap.xml", nil), "https://wave.example", entries)
	if !strings.Contains(index.Body.String(), "<sitemapindex") || !strings.Contains(index.Body.String(), "/sitemap.xml?page=2") {
		t.Fatal("missing sitemap index")
	}
	page := httptest.NewRecorder()
	writeSitemap(page, httptest.NewRequest("GET", "/sitemap.xml?page=2", nil), "https://wave.example", entries)
	var set sitemapURLSet
	if err := xml.Unmarshal(page.Body.Bytes(), &set); err != nil {
		t.Fatal(err)
	}
	if len(set.URLs) != 1 || set.URLs[0].Location != entries[sitemapPageSize].Location {
		t.Fatal("invalid sitemap page boundary")
	}
	for _, value := range []string{"0", "-1", "3", "invalid", "99999999999999999999999"} {
		response := httptest.NewRecorder()
		writeSitemap(response, httptest.NewRequest("GET", "/sitemap.xml?page="+value, nil), "https://wave.example", entries)
		if response.Code != 404 {
			t.Fatalf("invalid page %s returned %d", value, response.Code)
		}
	}
}

func TestKoreanBlogMetadataDoesNotDependOnUILanguage(t *testing.T) {
	service, closeDB := testSEOService(t)
	defer closeDB()
	if err := service.Repository().Upsert(blogdomain.Post{Slug: "korean", Title: "Wave 개발 소식", Category: "article", Status: "published", Content: "Wave 프로그래밍 언어의 새로운 소식을 전합니다. 이번 릴리즈에는 컴파일러 개선과 표준 라이브러리 기능이 포함되어 있습니다. 개발에 참여해 주신 모든 분께 감사드립니다."}); err != nil {
		t.Fatal(err)
	}
	seo := NewSEOHandler("https://wave.example", nil, service, nil, nil, nil)
	server := seoFrontend(t, seo)
	for _, locale := range []string{"en-US", "ko-KR", "ja-JP"} {
		request := httptest.NewRequest("GET", "/blog/korean", nil)
		request.Header.Set("Accept-Language", locale)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		if response.Header().Get("Content-Language") != "ko" || !strings.Contains(response.Body.String(), `content="ko_KR"`) || !strings.Contains(response.Body.String(), `<html lang="ko">`) {
			t.Fatalf("language varies with reader %s", locale)
		}
		if graphNode(t, metadataGraph(t, response.Body.String()), "BlogPosting")["inLanguage"] != "ko" {
			t.Fatal("article language mismatch")
		}
	}
}

func TestSafeHTMLSurvivesMarkdownRendering(t *testing.T) {
	markup := renderPublicMarkdown("<details><summary>Read more</summary><p>Useful explanation</p></details>\n\n<img src=\"/media/picture.webp\" onerror=\"alert(1)\">\n\n<iframe src=\"https://untrusted.example\"></iframe>")
	if !strings.Contains(markup, "<details>") || !strings.Contains(markup, "Useful explanation") || !strings.Contains(markup, "/media/picture.webp") {
		t.Fatal("safe article content lost")
	}
	if strings.Contains(markup, "onerror") || strings.Contains(markup, "<iframe") {
		t.Fatal("executable markup remains")
	}
}

func TestMultilingualHeadingAnchorsRemainReadable(t *testing.T) {
	markup := renderPublicMarkdown("## `ptr<T>`\n\n## 明示的なデリファレンス\n\n## 포인터 연산\n\n## 포인터 연산\n")
	for _, expected := range []string{`id="ptrt"`, `ptr&lt;T&gt;`, `id="明示的なデリファレンス"`, `id="포인터-연산"`, `id="포인터-연산-1"`} {
		if !strings.Contains(markup, expected) {
			t.Fatalf("missing %s in %s", expected, markup)
		}
	}
}
