package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootServesHTMLWithGIFs(t *testing.T) {
	h := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", got)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<h1") {
		t.Fatalf("body = %q, want heading", body)
	}
	if !strings.Contains(body, "<img") {
		t.Fatalf("body = %q, want an image tag", body)
	}
}

func TestIssue5UsesTheAdorbsGIFList(t *testing.T) {
	h := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/?count=3", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := strings.Count(rr.Body.String(), "<img"); got != 3 {
		t.Fatalf("img count = %d, want 3", got)
	}
	if !strings.Contains(rr.Body.String(), "https://s3.amazonaws.com/files.656.mba/mgt656/fall-2026/random-gifs/adorbs/") {
		t.Fatalf("body = %q, want GIF URLs from the adorbs list", rr.Body.String())
	}
}

func TestQueryCountControlsNumberOfGIFs(t *testing.T) {
	h := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/?count=3", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := strings.Count(rr.Body.String(), "<img"); got != 3 {
		t.Fatalf("img count = %d, want 3", got)
	}
}

func TestInvalidCountFallsBackToOneGIF(t *testing.T) {
	h := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/?count=banana", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := strings.Count(rr.Body.String(), "<img"); got != 1 {
		t.Fatalf("img count = %d, want 1", got)
	}
}
