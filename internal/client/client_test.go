package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDoSendsBearerAndAccountHeaders(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/mobile/v1/properties" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("missing bearer")
		}
		if r.Header.Get("X-Sunflower-Account-ID") != "42" {
			t.Errorf("missing account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer s.Close()
	c := &Client{BaseURL: s.URL, Token: "secret", AccountID: "42", HTTP: s.Client()}
	status, b, e := c.Do("POST", "/api/mobile/v1/properties", []byte(`{"name":"Cabin"}`), url.Values{"x": []string{"1"}})
	if e != nil || status != 200 || string(b) != `{"ok":true}` {
		t.Fatalf("got %d %s %v", status, b, e)
	}
}
func TestDoReturnsStructuredHTTPError(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, `{"error":"Unauthorized"}`, 401) }))
	defer s.Close()
	c := &Client{BaseURL: s.URL, HTTP: s.Client()}
	_, _, e := c.Do("GET", "/x", nil, nil)
	if e == nil || e.Error() != "API 401: Unauthorized" {
		t.Fatalf("error=%v", e)
	}
}
