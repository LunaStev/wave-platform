package web

import (
	"net/http/httptest"
	"testing"
)

func TestMovedDocumentationRedirects(t *testing.T) {
	seo := NewSEOHandler("https://wave.example", nil, nil, nil, nil, nil)
	server := seoFrontend(t, seo)
	for _, item := range []struct{ from, to string }{
		{"/docs/ko/learn/functions", "/docs/ko/language/functions-and-generics"},
		{"/docs/en/toolchain/build-link-targets", "/docs/en/whale/build-link-targets"},
		{"/docs/ja/getting-started/design-goals", "/docs/ja/getting-started/overview"},
	} {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", item.from, nil))
		if response.Code != 308 || response.Header().Get("Location") != item.to {
			t.Errorf("%s: status=%d location=%s", item.from, response.Code, response.Header().Get("Location"))
		}
	}
}
