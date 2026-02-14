package main

import (
	"net/http"
	"net/url"
	"testing"
)

type fakeResponseWriter struct {
	dest *[]byte
}

func (f fakeResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (f fakeResponseWriter) Write(p []byte) (int, error) {
	*f.dest = append(*f.dest, p...)
	return len(p), nil
}

func (f fakeResponseWriter) WriteHeader(statusCode int) {}

func BenchmarkTestProject(b *testing.B) {
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/"},
		Header: make(http.Header),
	}

	var result []byte
	resp := fakeResponseWriter{dest: &result}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		tpl := startTextwire()
		homeHandler(tpl)(resp, req)
	}

	_ = len(result)
}
