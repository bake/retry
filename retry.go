package retry

import (
	"net/http"
	"time"
)

// New returns a new transport that can retry a request multiple times between
// which it sleeps a given duration.
// The transport argument can be used to chain transports. If it is nil,
// http.DefaultTransport is used.
func New(retries int, backoff time.Duration, transport http.RoundTripper) http.RoundTripper {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return roundTripper{retries, backoff, transport}
}

type roundTripper struct {
	retries   int
	backoff   time.Duration
	transport http.RoundTripper
}

// RoundTrip will try to successfully send a request at most t.retries times.
// If fails the underlying http.RoundTriller retruns an error.
func (t roundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	for i := 0; i < t.retries; i++ {
		res, err = t.transport.RoundTrip(req)
		if err != nil {
			return res, err
		}
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			break
		}
		time.Sleep(t.backoff)
	}
	return res, err
}
