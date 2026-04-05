package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const baseTitle = "Go Tutorial Sample App"

func newTestStaticHandler(t *testing.T) *StaticHandler {
	t.Helper()
	return NewStaticHandler()
}

func TestStaticPagesHome(t *testing.T) {
	h := newTestStaticHandler(t)

	req := httptest.NewRequest("GET", "/static_pages/home", nil)
	rec := httptest.NewRecorder()
	h.Home(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Home: expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := "<title>Home | " + baseTitle + "</title>"
	if !strings.Contains(rec.Body.String(), expected) {
		t.Errorf("Home: expected title %q in body", expected)
	}
}

func TestStaticPagesHelp(t *testing.T) {
	h := newTestStaticHandler(t)

	req := httptest.NewRequest("GET", "/static_pages/help", nil)
	rec := httptest.NewRecorder()
	h.Help(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Help: expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := "<title>Help | " + baseTitle + "</title>"
	if !strings.Contains(rec.Body.String(), expected) {
		t.Errorf("Help: expected title %q in body", expected)
	}
}

func TestStaticPagesAbout(t *testing.T) {
	h := newTestStaticHandler(t)

	req := httptest.NewRequest("GET", "/static_pages/about", nil)
	rec := httptest.NewRecorder()
	h.About(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("About: expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := "<title>About | " + baseTitle + "</title>"
	if !strings.Contains(rec.Body.String(), expected) {
		t.Errorf("About: expected title %q in body", expected)
	}
}

func TestStaticPagesRoot(t *testing.T) {
	h := newTestStaticHandler(t)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.Home)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Root: expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
