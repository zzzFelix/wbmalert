package main

import (
	"testing"
)

func TestRemoveByRegex(t *testing.T) {
	want := "<p></p>"
	result := removeByRegex("<p>last update: 01.01.1970 - 00:00</p>", "last update: \\d\\d.\\d\\d.\\d\\d\\d\\d - \\d\\d:\\d\\d")
	if result != want {
		t.Fatalf(`Expected "%s", received "%s"`, want, result)
	}
}

func TestRemoveNothing(t *testing.T) {
	want := "<p>last update: 01.01.1970 - 00:00</p>"
	result := removeByRegex(want, "")
	if result != want {
		t.Fatalf(`Expected "%s", received "%s"`, want, result)
	}
}
