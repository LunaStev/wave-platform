// Package wavedoc owns the source files for Wave's official documentation.
// Keeping the Markdown at the repository root makes it usable independently
// from the platform server while still allowing the server to embed it.
package wavedoc

import (
	"embed"
	"encoding/json"
)

type Locale struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Script    string `json:"script"`
	OpenGraph string `json:"openGraph"`
	Complete  bool   `json:"complete"`
}

// LocaleData is shared by the server, browser, and translation checker.
//
//go:embed locales.json
var LocaleData []byte

var Locales = func() []Locale {
	var locales []Locale
	if err := json.Unmarshal(LocaleData, &locales); err != nil {
		panic(err)
	}
	return locales
}()

var SupportedLocales = func() []string {
	var locales []string
	for _, locale := range Locales {
		locales = append(locales, locale.ID)
	}
	return locales
}()

// RequiresCompleteCoverage keeps existing complete locales in sync with Korean.
func RequiresCompleteCoverage(id string) bool {
	for _, locale := range Locales {
		if locale.ID == id {
			return locale.Complete
		}
	}
	return false
}

// Content contains every translated Markdown document below this directory.
//
//go:embed */*/*.md
var Content embed.FS

func SupportsLocale(locale string) bool {
	for _, supported := range SupportedLocales {
		if locale == supported {
			return true
		}
	}
	return false
}

// RedirectData is shared with the browser so moved documentation has one URL map.
//
//go:embed redirects.json
var RedirectData []byte

var DocumentRedirects = func() map[string]string {
	var paths map[string]string
	if err := json.Unmarshal(RedirectData, &paths); err != nil {
		panic(err)
	}
	return paths
}()

func CanonicalDocumentPath(path string) string {
	if target, ok := DocumentRedirects[path]; ok {
		return target
	}
	return path
}
