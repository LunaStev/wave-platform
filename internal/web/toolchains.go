package web

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

// Development serves the same public subtree as Caddy without requiring a
// second server. OpenRoot keeps file access within that subtree, including
// when a requested file is a symlink. Missing catalogs remain a normal 404.
func developmentToolchainsHandler(directory string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Robots-Tag", "noindex")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("Cache-Control", "no-cache")
		name := strings.TrimPrefix(request.URL.Path, "/downloads/toolchains/")
		if !fs.ValidPath(name) {
			http.NotFound(writer, request)
			return
		}
		for _, segment := range strings.Split(name, "/") {
			if strings.HasPrefix(segment, ".") || strings.Contains(segment, `\`) {
				http.NotFound(writer, request)
				return
			}
		}

		root, err := os.OpenRoot(directory)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				http.NotFound(writer, request)
			} else {
				http.Error(writer, "Toolchain downloads unavailable", http.StatusServiceUnavailable)
			}
			return
		}
		defer root.Close()
		file, err := root.Open(name)
		if err != nil {
			http.NotFound(writer, request)
			return
		}
		defer file.Close()
		info, err := file.Stat()
		if err != nil || !info.Mode().IsRegular() {
			http.NotFound(writer, request)
			return
		}
		http.ServeContent(writer, request, info.Name(), info.ModTime(), file)
	})
}
