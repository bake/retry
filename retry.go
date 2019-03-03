// Package retry is a small implementation of the http.RoundTripper interface
// that can be found in http.Client. It is responsible to make HTTP requests
// and can be used to cache or retry them.
//
// By default a request will be tried again only if the response code is below
// 500 (Internal Server Error) and not 429 (Too Many Requests). This behavious
// can be changed with the WithVerifier option.
package retry

import (
	"net/http"
	"time"
)

// WithTransport replaces the http.DefaultTransport.
func WithTransport(t http.RoundTripper) func(*roundTripper) {
	return func(rt *roundTripper) {
		rt.transport = t
	}
}

// WithVerifier replaces the default verifier. If the given function returns
// true and there are retries left, another request will be sent.
func WithVerifier(fn func(*http.Response) bool) func(*roundTripper) {
	return func(rt *roundTripper) {
		rt.verifier = fn
	}
}

func defaultVerifier(res *http.Response) bool {
	if res.StatusCode < 500 {
		return false
	}
	if res.StatusCode == http.StatusTooManyRequests {
		return false
	}
	return true
}

// New returns a new http.RoundTripper that can retry a request multiple times
// between which it sleeps a given duration. The options arguments can be used
// to replace the http.DefaultTransport or the function that determines if a
// request should be retried based on the corresponding *http.Response.
func New(retries int, backoff time.Duration, options ...func(*roundTripper)) http.RoundTripper {
	rt := &roundTripper{retries, backoff, http.DefaultTransport, defaultVerifier}
	for _, option := range options {
		option(rt)
	}
	return rt
}

// roundTripper implements the http.RoundTripper interface.
type roundTripper struct {
	retries   int
	backoff   time.Duration
	transport http.RoundTripper
	verifier  func(*http.Response) bool
}

// RoundTrip will try to successfully send a request at most t.retries times.
// It fails when the underlying http.RoundTriller returns an error.
func (t *roundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	for i := 0; i < t.retries; i++ {
		res, err = t.transport.RoundTrip(req)
		if err != nil || !t.verifier(res) {
			break
		}
		time.Sleep(t.backoff)
	}
	return res, err
}
