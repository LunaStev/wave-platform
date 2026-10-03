package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestToolchainsPageMetadataAndStatus(t *testing.T) {
	handler := NewSEOHandler("https://wave-lang.dev", nil, nil, nil, nil, nil)
	if status := handler.StatusCode(httptest.NewRequest("GET", "/toolchains", nil)); status != http.StatusOK {
		t.Fatalf("toolchain page status = %d", status)
	}
	if status := handler.StatusCode(httptest.NewRequest("GET", "/toolchains/missing", nil)); status != http.StatusNotFound {
		t.Fatalf("unknown toolchain page status = %d", status)
	}
	meta := handler.metadata("/toolchains", "https://wave-lang.dev")
	if !strings.Contains(meta.Title, "LLVM Toolchains") || meta.SchemaType != "CollectionPage" {
		t.Fatalf("unexpected metadata: %+v", meta)
	}
}
