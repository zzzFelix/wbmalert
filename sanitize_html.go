package wbmalert

import (
	"strings"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

func removeAllWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func stripTags(s string) string {
	p := bluemonday.StripTagsPolicy()
	html := p.Sanitize(s)
	return html
}

func sanitizeHtml(s string) string {
	content := stripTags(s)
	content = removeAllWhitespace(content)
	return content
}
