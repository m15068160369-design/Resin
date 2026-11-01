package proxy

import (
	"net/http/httptest"
	"testing"
)

// [patch lou] All proxy error responses must advertise Connection: close so
// clients drop the connection instead of pooling it after an early (body
// unread) error response.
func TestWriteProxyErrorSetsConnectionClose(t *testing.T) {
	cases := []*ProxyError{
		ErrAuthRequired,
		ErrAuthFailed,
		ErrURLParseError,
		ErrInvalidProtocol,
		ErrInvalidHost,
		ErrPlatformNotFound,
		ErrAccountRejected,
		ErrNoAvailableNodes,
		ErrUpstreamConnectFailed,
		ErrUpstreamTimeout,
		ErrUpstreamRequestFailed,
		ErrInternalError,
	}
	for _, pe := range cases {
		rec := httptest.NewRecorder()
		writeProxyError(rec, pe)
		resp := rec.Result()
		if got := resp.Header.Get("Connection"); got != "close" {
			t.Errorf("ProxyError %s (HTTP %d): Connection header = %q, want %q",
				pe.ResinError, pe.HTTPCode, got, "close")
		}
	}
}

// The 407 path must still add Proxy-Authenticate alongside Connection: close.
func TestWriteProxyErrorAuthRequiredKeepsAuthenticateHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	writeProxyError(rec, ErrAuthRequired)
	resp := rec.Result()
	if got := resp.Header.Get("Proxy-Authenticate"); got == "" {
		t.Error("Proxy-Authenticate header missing on 407 response")
	}
	if got := resp.Header.Get("Connection"); got != "close" {
		t.Errorf("Connection header = %q, want close", got)
	}
}
