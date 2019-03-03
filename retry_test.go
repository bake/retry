package retry_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bakerolls/retry"
)

// Create a new http.Client that will retry requests five times and sleeps
// one second between each one. Since the requested URL will always return a
// status code of 500 (Internal Server Error), this function will run for five
// seconds and the time it takes to call the endpoint five times.
func ExampleNew() {
	// Create a *http.Client that will retry requests five times and sleeps one
	// second between each retried request.
	client := &http.Client{
		Transport: retry.New(5, time.Second),
	}
	res, err := client.Get("https://httpbin.org/status/500")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	fmt.Printf("%s\n", res.Status)
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		log.Fatal(err)
	}
}

// WithVerifier can be used to modify the default verifier that is used to
// determine if a request can be tried again.
func ExampleWithVerifier() {
	// Only retry if the responses status code is below 500. This is similar to
	// the default behaviour, but ignores 429 (Too Many Requests).
	// With the default verifier this function would take about five seconds.
	verifier := func(res *http.Response) bool {
		return res.StatusCode < 500
	}
	client := &http.Client{
		Transport: retry.New(5, time.Second, retry.WithVerifier(verifier)),
	}
	res, err := client.Get("https://httpbin.org/status/429")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	fmt.Printf("%s\n", res.Status)
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		log.Fatal(err)
	}
}
