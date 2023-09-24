package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	HTTP_GET_BODY = "<p>Test</p>"
	WANT          = "Test"
)

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	reader := strings.NewReader(HTTP_GET_BODY)
	readCloser := io.NopCloser(reader)
	response := http.Response{
		Body: readCloser,
	}
	return &response, nil
}

func TestGetWebsiteAsString(t *testing.T) {
	client = &mockClient{}
	website := website{
		Name:     "Test",
		Url:      "https://google.com",
		Snapshot: "",
	}
	result, error := getWebsiteAsString(&website)
	if error != nil {
		t.Fatal(error)
	}
	if result != WANT {
		t.Fatalf(`%q, want match for %q`, result, WANT)
	}
}
