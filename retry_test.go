package retry_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bake/retry"
)

// Create a new HTTP client that will retry requests five times and sleeps
// one second between each one. Since the requested URL will always return a
// status code of 500 (Internal Server Error), this function will run for five
// seconds and the time it takes to call the endpoint five times.
func ExampleNew() {
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

// WithVerifier can be used to modify the default behaviour that is used to
// determine if a request can be repeated. In this example, all responses that
// do not succeed (with a status code outside [200-300[) will be tried again.
func ExampleWithVerifier() {
	verifier := func(res *http.Response) bool {
		return 200 > res.StatusCode || res.StatusCode >= 300
	}
	client := &http.Client{
		Transport: retry.New(5, time.Second, retry.WithVerifier(verifier)),
	}
	res, err := client.Get("https://httpbin.org/status/418")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	fmt.Printf("%s\n", res.Status)
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		log.Fatal(err)
	}
}
