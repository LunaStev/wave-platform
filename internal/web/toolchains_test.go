package web

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func toolchainsTestRouter(t *testing.T, environment, directory string) http.Handler {
	t.Helper()
	return NewRouter(environment, t.TempDir(), directory, "https://wave.example", "test", nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
}

func TestDevelopmentToolchainDownloads(t *testing.T) {
	directory := t.TempDir()
	archive := "llvm/21.1.8/r1/wave-llvm-21.1.8-linux-loong64-r1.tar.xz"
	content := "test archive content"
	checksum := fmt.Sprintf("%x  %s\n", sha256.Sum256([]byte(content)), filepath.Base(archive))
	files := map[string]string{
		"index.json":        `{"schema_version":1,"bundles":[]}`,
		archive:             content,
		archive + ".sha256": checksum,
		".hidden":           "private",
		".staging/sdk":      "unpublished",
	}
	for name, body := range files {
		path := filepath.Join(directory, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	router := toolchainsTestRouter(t, "development", directory)
	for _, check := range []struct {
		name, method, path, byteRange, body string
		status                              int
	}{
		{"catalog", "GET", "index.json", "", files["index.json"], 200},
		{"archive", "GET", archive, "", content, 200},
		{"checksum", "GET", archive + ".sha256", "", checksum, 200},
		{"head", "HEAD", archive, "", "", 200},
		{"range", "GET", archive, "bytes=0-3", "test", 206},
		{"missing", "GET", "llvm/missing.tar.xz", "", "", 404},
		{"hidden", "GET", ".hidden", "", "", 404},
		{"staging", "GET", ".staging/sdk", "", "", 404},
		{"directory", "GET", "llvm/21.1.8/r1/", "", "", 404},
	} {
		t.Run(check.name, func(t *testing.T) {
			request := httptest.NewRequest(check.method, "/downloads/toolchains/"+check.path, nil)
			if check.byteRange != "" {
				request.Header.Set("Range", check.byteRange)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != check.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, check.status, response.Body.String())
			}
			if check.status < 400 && response.Body.String() != check.body {
				t.Fatalf("body = %q, want %q", response.Body.String(), check.body)
			}
			if response.Header().Get("Cache-Control") != "no-cache" || response.Header().Get("X-Robots-Tag") != "noindex" {
				t.Fatalf("unexpected download headers: %v", response.Header())
			}
			if check.name == "catalog" && !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
				t.Fatalf("catalog content type = %q", response.Header().Get("Content-Type"))
			}
			if check.name == "range" && response.Header().Get("Content-Range") != "bytes 0-3/20" {
				t.Fatalf("content range = %q", response.Header().Get("Content-Range"))
			}
		})
	}

	// Production must keep using Caddy even when a public directory is present.
	response := httptest.NewRecorder()
	toolchainsTestRouter(t, "production", directory).ServeHTTP(response, httptest.NewRequest("GET", "/downloads/toolchains/index.json", nil))
	if response.Code != http.StatusNotFound || strings.Contains(response.Body.String(), "schema_version") {
		t.Fatalf("production Go server exposed the catalog: %d %s", response.Code, response.Body.String())
	}
}

func TestDevelopmentToolchainRootIsolation(t *testing.T) {
	parent := t.TempDir()
	directory := filepath.Join(parent, "toolchains")
	router := toolchainsTestRouter(t, "development", directory)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest("GET", "/downloads/toolchains/index.json", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("absent catalog status = %d", response.Code)
	}
	if err := os.Mkdir(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(parent, "private.txt")
	if err := os.WriteFile(secret, []byte("private storage"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"../private.txt", "%2e%2e%2fprivate.txt", "%5c..%5cprivate.txt"} {
		response := httptest.NewRecorder()
		developmentToolchainsHandler(directory).ServeHTTP(response, httptest.NewRequest("GET", "/downloads/toolchains/"+path, nil))
		if response.Code != http.StatusNotFound || strings.Contains(response.Body.String(), "private storage") {
			t.Fatalf("path %q escaped root: %d %s", path, response.Code, response.Body.String())
		}
	}
	t.Run("external symlink", func(t *testing.T) {
		if err := os.Symlink(secret, filepath.Join(directory, "linked.txt")); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest("GET", "/downloads/toolchains/linked.txt", nil))
		if response.Code != http.StatusNotFound || strings.Contains(response.Body.String(), "private storage") {
			t.Fatalf("symlink escaped root: %d %s", response.Code, response.Body.String())
		}
	})
}
