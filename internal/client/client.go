package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/theinventor/sunflower-tasks-cli/internal/config"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL, Token, AccountID string
	HTTP                      *http.Client
}

func New(profile string) (*Client, error) {
	_, p, e := config.Resolve(profile)
	if e != nil {
		return nil, e
	}
	return &Client{BaseURL: strings.TrimRight(p.APIURL, "/"), Token: p.APIKey, AccountID: p.AccountID, HTTP: &http.Client{Timeout: 45 * time.Second}}, nil
}
func (c *Client) Do(method, path string, body []byte, query url.Values) (int, []byte, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := c.BaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	var r io.Reader
	if len(body) > 0 {
		r = bytes.NewReader(body)
	}
	req, e := http.NewRequest(method, u, r)
	if e != nil {
		return 0, nil, e
	}
	req.Header.Set("Accept", "application/json")
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	if c.AccountID != "" {
		req.Header.Set("X-Sunflower-Account-ID", c.AccountID)
	}
	resp, e := c.HTTP.Do(req)
	if e != nil {
		return 0, nil, e
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return resp.StatusCode, b, &HTTPError{Status: resp.StatusCode, Body: b}
	}
	return resp.StatusCode, b, nil
}

type HTTPError struct {
	Status int
	Body   []byte
}

func (e *HTTPError) Error() string {
	var x struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(e.Body, &x) == nil && x.Error != "" {
		return fmt.Sprintf("API %d: %s", e.Status, x.Error)
	}
	return fmt.Sprintf("API %d: %s", e.Status, strings.TrimSpace(string(e.Body)))
}
