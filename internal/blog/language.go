package blog

import (
	"regexp"
	"strings"

	"github.com/abadojack/whatlanggo"
)

var languageCodeBlock = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~|`[^`]*`")
var languageURL = regexp.MustCompile(`https?://\S+`)

// ContentLanguage describes the article, independently of the reader's UI.
// Short, ambiguous text uses the official blog's English default. No text is
// sent to an external detection service, and existing records need no migration.
func (post Post) ContentLanguage() string {
	prose := languageCodeBlock.ReplaceAllString(post.Title+"\n"+post.Summary+"\n"+post.Content, " ")
	prose = languageURL.ReplaceAllString(prose, " ")
	runes := []rune(strings.TrimSpace(prose))
	if len(runes) > 12000 {
		runes = runes[:12000]
	}
	info := whatlanggo.Detect(string(runes))
	if info.IsReliable() {
		if code := info.Lang.Iso6391(); code != "" {
			return code
		}
	}
	return "en"
}
