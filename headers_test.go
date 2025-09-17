package opengraph

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/otiai10/mint"
)

func TestFetch_WithBotProtection(t *testing.T) {
	// Create a test server that mimics bot protection behavior
	// by checking User-Agent and Accept headers
	protectedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		accept := r.Header.Get("Accept")

		// Simulate bot protection: reject requests without browser-like headers
		if userAgent == "" || userAgent == "Go-http-client/1.1" || userAgent == "Go-http-client/2.0" {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("<html><body>Access denied: Bot detected</body></html>"))
			return
		}

		if accept == "" || !strings.Contains(accept, "text/html") {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("<html><body>Access denied: Invalid Accept header</body></html>"))
			return
		}

		// Return valid HTML with OpenGraph tags if headers are correct
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<meta property="og:title" content="Protected Content">
	<meta property="og:description" content="This content requires browser-like headers">
	<meta property="og:type" content="article">
	<meta property="og:url" content="` + r.URL.String() + `">
</head>
<body>
	<h1>Protected Content</h1>
</body>
</html>`))
	}))
	defer protectedServer.Close()

	When(t, "fetching with default headers (should fail currently)", func(t *testing.T) {
		// This test demonstrates the issue reported in #33
		// We expect to get the OpenGraph data, but with default headers we get nothing
		og, err := Fetch(protectedServer.URL)

		// This SHOULD work but currently fails due to lack of browser-like headers
		// After the fix, these assertions should pass
		Expect(t, err).ToBe(nil)
		Expect(t, og.Title).ToBe("Protected Content")  // This will FAIL with current implementation
		Expect(t, og.Description).ToBe("This content requires browser-like headers")  // This will also FAIL
	})

	When(t, "fetching with custom browser-like headers", func(t *testing.T) {
		// This demonstrates the workaround users need to implement
		client := &http.Client{
			Transport: &headerTransport{
				headers: map[string]string{
					"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
					"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
					"Accept-Language": "en-US,en;q=0.9",
				},
			},
		}

		og := New(protectedServer.URL)
		og.Intent.HTTPClient = client
		err := og.Fetch()

		// With proper headers, this should succeed
		Expect(t, err).ToBe(nil)
		Expect(t, og.Title).ToBe("Protected Content")
		Expect(t, og.Description).ToBe("This content requires browser-like headers")
	})
}

// Helper transport to add custom headers
type headerTransport struct {
	headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return http.DefaultTransport.RoundTrip(req)
}