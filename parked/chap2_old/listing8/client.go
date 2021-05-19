package client

import (
	"net/http"
)

type AddHeadersMiddleware struct {
	headers map[string]string
}

func (h AddHeadersMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {
	reqCopy := r.Clone(r.Context())
	for k, v := range h.headers {
		reqCopy.Header.Add(k, v)
	}
	return http.DefaultTransport.RoundTrip(reqCopy)
}

func createClient(headers map[string]string) *http.Client {
	h := AddHeadersMiddleware{
		headers: headers,
	}
	client := http.Client{
		Transport: &h,
	}
	return &client
}
