package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ServiceProxy struct {
	client *http.Client
}

func NewServiceProxy() *ServiceProxy {
	return &ServiceProxy{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// ProxyTo creates a handler that proxies requests to the target service
func (p *ServiceProxy) ProxyTo(targetURL string, stripPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse target URL
		target, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Build new path
		path := r.URL.Path
		if stripPrefix != "" {
			path = strings.TrimPrefix(path, stripPrefix)
		}

		// Build proxy URL
		proxyURL := target.ResolveReference(&url.URL{
			Path:     path,
			RawQuery: r.URL.RawQuery,
		})

		// Create proxy request
		proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, proxyURL.String(), r.Body)
		if err != nil {
			http.Error(w, `{"error":"failed to create request"}`, http.StatusInternalServerError)
			return
		}

		// Copy headers
		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Add forwarding headers
		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
		proxyReq.Header.Set("X-Forwarded-Host", r.Host)
		proxyReq.Header.Set("X-Forwarded-Proto", "http")
		if r.TLS != nil {
			proxyReq.Header.Set("X-Forwarded-Proto", "https")
		}

		// Send request
		start := time.Now()
		resp, err := p.client.Do(proxyReq)
		duration := time.Since(start)

		if err != nil {
			log.Printf("⚠️ Proxy error to %s: %v (took %v)", proxyURL.Host, err, duration)
			http.Error(w, `{"error":"service unavailable"}`, http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		log.Printf("→ %s %s → %s (%d) [%v]", r.Method, r.URL.Path, proxyURL.Host, resp.StatusCode, duration)

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Write status code
		w.WriteHeader(resp.StatusCode)

		// Copy response body
		io.Copy(w, resp.Body)
	}
}
