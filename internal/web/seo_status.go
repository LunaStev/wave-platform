package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	documentdomain "github.com/wavefnd/wave-platform/internal/document"
	"github.com/wavefnd/wave-platform/internal/storage"
)

func contentStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if errors.Is(err, storage.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusServiceUnavailable
}

func (handler SEOHandler) publicStatus(request *http.Request) (int, bool) {
	parts := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	if len(parts) < 2 {
		return 0, false
	}
	switch request.URL.Path {
	case "/community/new", "/community/showcase/new", "/lunastev/new", "/questions/new", "/rfcs/new":
		return 0, false
	}
	switch parts[0] {
	case "docs":
		if handler.documents == nil {
			return 0, false
		}
		locale, start := "en", 1
		if supportedDocumentLocale(parts[1]) {
			locale, start = parts[1], 2
		}
		documentPath := strings.Join(parts[start:], "/")
		if documentPath == "" || documentPath == "whale" {
			_, err := handler.documents.Navigation(locale, documentdomain.ProjectForPath(documentPath+"/"))
			return contentStatus(err), true
		}
		_, err := handler.documents.Published(locale, documentPath)
		if errors.Is(err, storage.ErrNotFound) && locale != "en" {
			_, err = handler.documents.Published("en", documentPath)
		}
		return contentStatus(err), true
	case "community", "lunastev":
		if handler.community == nil {
			return 0, false
		}
		if len(parts) == 2 && parts[1] == "showcase" && parts[0] == "community" {
			return http.StatusOK, true
		}
		if len(parts) != 3 || parts[1] != "thread" && parts[1] != "showcase" {
			return http.StatusNotFound, true
		}
		thread, err := handler.community.Thread(parts[2])
		if err != nil {
			return contentStatus(err), true
		}
		personal := thread.SpaceID == "founder-notes" || thread.SpaceID == "development-log"
		if personal != (parts[0] == "lunastev") || (parts[1] == "showcase" && thread.SpaceID != "showcase") {
			return http.StatusNotFound, true
		}
		return http.StatusOK, true
	case "questions":
		if handler.questions == nil {
			return 0, false
		}
		if len(parts) != 2 {
			return http.StatusNotFound, true
		}
		_, err := handler.questions.Question(parts[1])
		return contentStatus(err), true
	case "rfcs":
		if handler.rfcs == nil {
			return 0, false
		}
		if len(parts) != 2 && !(len(parts) == 3 && parts[2] == "edit") {
			return http.StatusNotFound, true
		}
		number, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil || number == 0 {
			return http.StatusNotFound, true
		}
		_, err = handler.rfcs.Repository().Proposal(number)
		return contentStatus(err), true
	case "source":
		if handler.source == nil {
			return 0, false
		}
		if len(parts) != 2 {
			return http.StatusNotFound, true
		}
		_, err := handler.source.Repository(request.Context(), parts[1])
		return contentStatus(err), true
	}
	return 0, false
}
