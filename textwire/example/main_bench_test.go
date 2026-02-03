package main

import (
	"net/http"
	"net/url"
	"testing"
)

type fakeResponseWriter struct{}

func (f fakeResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (f fakeResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (f fakeResponseWriter) WriteHeader(statusCode int) {}

func BenchmarkTestProject(b *testing.B) {
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/"},
		Header: make(http.Header),
	}

	resp := fakeResponseWriter{}

	b.ResetTimer()

	for b.Loop() {
		tpl := startTextwire()
		homeHandler(tpl)(resp, req)
	}
}
