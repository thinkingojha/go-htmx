package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

func TestMain(m *testing.M) {
	// Setup templates for testing
	templateDir := filepath.Join("..", "..", "internal", "template")
	utils.ParseTemplates(templateDir)
	m.Run()
}

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	err = HomeHandler(rr, req)
	if err != nil {
		t.Errorf("HomeHandler returned an error: %v", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestExpHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	err = ExpHandler(rr, req)
	if err != nil {
		t.Errorf("ExpHandler returned an error: %v", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestMarkdownHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "GET request without content",
			method:         "GET",
			url:            "/write",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET request with content",
			method:         "GET",
			url:            "/write?content=# Test\nHello world",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			err = MarkdownHandler(rr, req)
			if err != nil {
				t.Errorf("MarkdownHandler returned an error: %v", err)
			}

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestMarkdownHandlerHTMX(t *testing.T) {
	req, err := http.NewRequest("POST", "/write", strings.NewReader("content=# Test\nHello **world**"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()

	err = MarkdownHandler(rr, req)
	if err != nil {
		t.Errorf("MarkdownHandler returned an error: %v", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that it returns HTML content for HTMX request
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %s", contentType)
	}

	// Check that response contains rendered HTML
	body := rr.Body.String()
	if !strings.Contains(body, "<h1") {
		t.Errorf("Expected rendered HTML with h1 tag, got: %s", body)
	}
}
