package apiclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type APIClient struct {
	BaseUri string
	Headers map[string]string
	Timeout time.Duration
	Verify  bool
	Client  *http.Client
}

func NewAPIClient(baseUri string, headers map[string]string, timeout time.Duration, verify bool) *APIClient {
	cli := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: verify},
		},
	}
	return &APIClient{Headers: headers, BaseUri: baseUri, Timeout: timeout, Verify: verify, Client: cli}
}

func (c *APIClient) Wrap() *APIClient {
	transport := http.DefaultTransport
	if c.Client.Transport != nil {
		transport = c.Client.Transport
	}

	c.Client = &http.Client{
		Transport: &wrappedTransport{
			headers:   c.Headers,
			baseUri:   c.BaseUri,
			transport: transport,
		},
	}

	return c
}

type wrappedTransport struct {
	headers   map[string]string
	baseUri   string
	transport http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (t *wrappedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}

	if t.baseUri != "" {
		rel, _ := url.Parse(t.baseUri + req.URL.String())
		req.URL = req.URL.ResolveReference(rel)
	}

	return t.transport.RoundTrip(req)
}

func cloneRequest(r *http.Request) *http.Request {
	r2 := &http.Request{}
	*r2 = *r
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

func (c *APIClient) GET(path string, opts map[string]string) (*http.Response, error) {
	opt := toUrlValues(opts)
	req, err := http.NewRequest("GET", addParams(path, opt), nil)
	if err != nil {
		return nil, err
	}

	return c.Client.Do(req)

}

func (c *APIClient) POST(path string, data interface{}) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return c.Client.Do(req)

}

func (c *APIClient) PUT(path string, data interface{}) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return c.Client.Do(req)

}

func (c *APIClient) DELETE(path string, opts map[string]string) (*http.Response, error) {
	opt := toUrlValues(opts)
	req, err := http.NewRequest("DELETE", addParams(path, opt), nil)
	if err != nil {
		return nil, err
	}

	return c.Client.Do(req)

}

func addParams(url_ string, params url.Values) string {
	if len(params) == 0 {
		return url_
	}

	if !strings.Contains(url_, "?") {
		url_ += "?"
	}

	if strings.HasSuffix(url_, "?") || strings.HasSuffix(url_, "&") {
		url_ += params.Encode()
	} else {
		url_ += "&" + params.Encode()
	}

	return url_
}

func toUrlValues(t map[string]string) url.Values {
	rst := make(url.Values)
	for k, v := range t {
		rst.Add(k, v)
	}
	return rst
}
