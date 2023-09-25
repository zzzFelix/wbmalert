package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	HTML_TAG_START = 60 // Unicode `<`
	HTML_TAG_END   = 62 // Unicode `>`
)

func removeAllWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// See https://stackoverflow.com/a/64701836
func stripTags(s string) string {
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	in := false
	start := 0
	end := 0

	for i, c := range s {
		if (i+1) == len(s) && end >= start {
			builder.WriteString(s[end:])
		}

		if c != HTML_TAG_START && c != HTML_TAG_END {
			continue
		}

		if c == HTML_TAG_START {
			if !in {
				start = i
			}
			in = true

			builder.WriteString(s[end:start])
			continue
		}

		in = false
		end = i + 1
	}
	s = builder.String()
	return s
}

func sanitizeHtml(s string) string {
	content := stripTags(s)
	content = removeAllWhitespace(content)
	return content
}
