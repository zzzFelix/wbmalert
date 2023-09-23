package main

import (
	"testing"
)

func TestRemoveSingleWhitespace(t *testing.T) {
	input := "Ber lin"
	want := "Berlin"
	result := removeAllWhitespace(input)
	if want != result {
		t.Fatalf(`%q, want match for %q`, result, want)
	}
}

func TestRemoveLeadingWhitespace(t *testing.T) {
	input := " Berlin"
	want := "Berlin"
	result := removeAllWhitespace(input)
	if want != result {
		t.Fatalf(`%q, want match for %q`, result, want)
	}
}

func TestRemoveTrailingWhitespace(t *testing.T) {
	input := " Berlin"
	want := "Berlin"
	result := removeAllWhitespace(input)
	if want != result {
		t.Fatalf(`%q, want match for %q`, result, want)
	}
}

func TestRemoveLineBreakWhitespace(t *testing.T) {
	input := "Ber\nlin"
	want := "Berlin"
	result := removeAllWhitespace(input)
	if want != result {
		t.Fatalf(`%q, want match for %q`, result, want)
	}
}

func TestRemoveAllWhitespace(t *testing.T) {
	input := " B\n\te\vr \fli\rn "
	want := "Berlin"
	result := removeAllWhitespace(input)
	if want != result {
		t.Fatalf(`%q, want match for %q`, result, want)
	}
}

func TestStripLinkTag(t *testing.T) {
	input := `<a href="https://belin.de">Berlin</a>`
	want := "Berlin"
	result := stripTags(input)
	if want != result {
		t.Fatalf(`%q, want match for %q`, result, want)
	}
}

func TestSanitizeHTML(t *testing.T) {
	input := "<a href='https://belin.de'>Berlin</a>\n<p> Berlin </p>\r<h1> Ber lin </h1>"
	want := "BerlinBerlinBerlin"
	result := sanitizeHtml(input)
	if want != result {
		t.Fatalf(`%q, want match for %q`, result, want)
	}
}
