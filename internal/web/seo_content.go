package web

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html"
	"net/url"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	documentdomain "github.com/wavefnd/wave-platform/internal/document"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldhtml "github.com/yuin/goldmark/renderer/html"
)

var documentationLocales = []string{"en", "ko", "ja", "zh", "es", "de", "ru", "id", "vi"}

var publicMarkdown = goldmark.New(goldmark.WithExtensions(extension.GFM), goldmark.WithParserOptions(parser.WithAutoHeadingID()), goldmark.WithRendererOptions(goldhtml.WithUnsafe()))
var publicHTMLPolicy = func() *bluemonday.Policy {
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("id").OnElements("h1", "h2", "h3", "h4", "h5", "h6")
	policy.AllowElements("details", "summary")
	policy.AllowAttrs("open").OnElements("details")
	policy.RequireNoFollowOnLinks(false)
	return policy
}()

func renderPublicMarkdown(source string) string {
	var output bytes.Buffer
	// Sanitize after Markdown conversion, preserving safe author markup while
	// removing executable HTML, event handlers, and dangerous URL schemes.
	context := parser.NewContext(parser.WithIDs(&publicHeadingIDs{used: make(map[string]bool)}))
	if err := publicMarkdown.Convert([]byte(source), &output, parser.WithContext(context)); err != nil {
		return "<p>" + html.EscapeString(source) + "</p>"
	}
	return publicHTMLPolicy.Sanitize(output.String())
}

// Preserve the multilingual anchors used by the document navigation. Goldmark's
// default generator drops non-ASCII characters from heading IDs.
type publicHeadingIDs struct{ used map[string]bool }

func (ids *publicHeadingIDs) Put(value []byte) { ids.used[string(value)] = true }

func (ids *publicHeadingIDs) Generate(value []byte, _ ast.NodeKind) []byte {
	var slug strings.Builder
	for _, character := range strings.ToLower(strings.TrimSpace(string(value))) {
		switch {
		case character == ' ':
			slug.WriteByte('-')
		case character == '-' || character == '_' || character >= '0' && character <= '9' || character >= 'a' && character <= 'z' || character >= 0x80:
			slug.WriteRune(character)
		}
	}
	base := slug.String()
	result := base
	for suffix := 1; ids.used[result]; suffix++ {
		result = base + "-" + strconv.Itoa(suffix)
	}
	ids.used[result] = true
	return []byte(result)
}

func openGraphLocale(language string) string {
	locales := map[string]string{"en": "en_US", "ko": "ko_KR", "ja": "ja_JP", "zh": "zh_CN", "es": "es_ES", "de": "de_DE", "ru": "ru_RU", "id": "id_ID", "vi": "vi_VN"}
	if locale := locales[language]; locale != "" {
		return locale
	}
	return language
}

// Only published translations are advertised. An English fallback is never
// presented as a translation in another language.
func (handler SEOHandler) documentAlternates(base, documentPath string) []seoAlternate {
	if handler.documents == nil {
		return nil
	}
	var result []seoAlternate
	locales := documentationLocales
	if documentPath != "" {
		locales = nil
		for _, translated := range handler.documents.PublishedTranslations(documentPath) {
			locales = append(locales, translated.Locale)
		}
	}
	for _, locale := range locales {
		location := handler.location(base, "docs", locale, documentPath)
		result = append(result, seoAlternate{Rel: "alternate", Language: locale, URL: location})
		if locale == "en" {
			result = append(result, seoAlternate{Rel: "alternate", Language: "x-default", URL: location})
		}
	}
	return result
}

func (handler SEOHandler) documentCatalogAlternates(base, project string) []seoAlternate {
	result := handler.documentAlternates(base, "")
	if project == "whale" {
		for i := range result {
			result[i].URL += "/whale"
		}
	}
	return result
}

// This is visible, useful HTML for every visitor, including browsers without
// JavaScript. Vue replaces the app contents when it mounts; no bot detection or
// hidden crawler-only copy is used.
func (handler SEOHandler) htmlContent(metadata pageMetadata) string {
	if strings.HasPrefix(metadata.Robots, "noindex") && metadata.Title != "Page not found · Wave" {
		return ""
	}
	if metadata.Markdown == "" && len(metadata.Items) == 0 && metadata.Title != "Page not found · Wave" && metadata.SchemaType != "CollectionPage" {
		return ""
	}
	escape := html.EscapeString
	var body strings.Builder
	body.WriteString(`<main class="public-reader" lang="` + escape(metadata.Language) + `"><nav aria-label="Main navigation"><a href="/">Wave</a><a href="/blog">Blog</a><a href="/releases">Releases</a><a href="/docs/en">Documentation</a></nav>`)
	if metadata.DocumentProject != "" {
		body.WriteString(`<nav aria-label="Documentation project">`)
		for _, project := range []string{"wave", "whale"} {
			path, name := "/docs/"+metadata.DocumentLocale, "Wave"
			if project == "whale" {
				path, name = path+"/whale", "Whale"
			}
			current := ""
			if project == metadata.DocumentProject {
				current = ` aria-current="page"`
			}
			body.WriteString(`<a href="` + escape(path) + `"` + current + `>` + name + `</a>`)
		}
		body.WriteString(`</nav>`)
		if metadata.SchemaType == "CollectionPage" && len(metadata.Items) == 0 {
			body.WriteString(`<p>No documents have been published for this project yet.</p>`)
		}
	}
	if len(metadata.Alternates) > 0 {
		body.WriteString(`<nav aria-label="Document language">`)
		for _, alt := range metadata.Alternates {
			if alt.Language != "x-default" {
				body.WriteString(`<a href="` + escape(alt.URL) + `" hreflang="` + escape(alt.Language) + `">` + escape(alt.Language) + `</a>`)
			}
		}
		body.WriteString(`</nav>`)
	}
	title := metadata.Headline
	if title == "" {
		title = metadata.Title
	}
	body.WriteString(`<article><header><h1>` + escape(title) + `</h1><p>` + escape(metadata.Description) + `</p>`)
	if metadata.AuthorName != "" {
		body.WriteString(`<p>` + escape(metadata.AuthorName) + `</p>`)
	}
	if metadata.PublishedAt != "" {
		body.WriteString(`<time datetime="` + escape(metadata.PublishedAt) + `">` + escape(metadata.PublishedAt) + `</time>`)
	}
	body.WriteString(`</header><div class="markdown-content">` + renderPublicMarkdown(metadata.Markdown) + `</div>`)
	if len(metadata.Items) > 0 {
		body.WriteString(`<ul class="public-reader-index">`)
		for _, item := range metadata.Items {
			body.WriteString(`<li><a href="` + escape(item.URL) + `">` + escape(item.Name) + `</a></li>`)
		}
		body.WriteString(`</ul>`)
	}
	if len(metadata.Comments) > 0 {
		body.WriteString(`<section aria-label="Comments"><h2>Comments</h2>`)
		for _, comment := range metadata.Comments {
			body.WriteString(`<article><p><a href="` + escape(comment.AuthorURL) + `">` + escape(comment.AuthorName) + `</a></p><p>` + escape(comment.Text) + `</p></article>`)
		}
		body.WriteString(`</section>`)
	}
	body.WriteString(`</article></main>`)
	return body.String()
}

// Seed only anonymous public read responses. This avoids an empty loading
// screen and duplicate network requests on the first client render.
func (handler SEOHandler) pageData(requestPath string) string {
	segments := strings.Split(strings.Trim(requestPath, "/"), "/")
	data := map[string]string{}
	add := func(key string, value any) {
		if encoded, err := xml.Marshal(value); err == nil {
			data[key] = string(encoded)
		}
	}
	if len(segments) >= 2 && segments[0] == "docs" && supportedDocumentLocale(segments[1]) && handler.documents != nil {
		locale := segments[1]
		for _, listLocale := range []string{locale, "en"} {
			if items, err := handler.documents.Summaries(listLocale); err == nil {
				add("/api/v1/documents?locale="+listLocale, struct {
					XMLName xml.Name                 `xml:"documents"`
					Items   []documentdomain.Summary `xml:"document"`
				}{Items: items})
			}
		}
		if len(segments) > 2 {
			documentPath := strings.Join(segments[2:], "/")
			if view, err := handler.documents.Published(locale, documentPath); err == nil {
				parts := append([]string(nil), segments[2:]...)
				for i := range parts {
					parts[i] = url.PathEscape(parts[i])
				}
				add("/api/v1/documents/"+strings.Join(parts, "/")+"?locale="+locale, view)
			}
		}
	}
	if len(segments) == 2 && (segments[0] == "blog" || segments[0] == "releases") && segments[1] != "editor" && handler.blog != nil {
		if post, err := handler.blog.Repository().Post(segments[1], false); err == nil && (segments[0] == "blog" || post.Category == "release") {
			add("/api/v1/blog/posts/"+url.PathEscape(post.Slug), post)
		}
	}
	if len(data) == 0 {
		return ""
	}
	encoded, _ := json.Marshal(data) // escapes HTML, including closing script tags
	return `<script type="application/json" id="wave-page-data">` + string(encoded) + `</script>`
}
