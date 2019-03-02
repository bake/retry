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
// one second between each one.
//
// Since the requested URL will always return a status code outside of
// [200, 300[, this function will run for five seconds and the time it takes to
// call the endpoint five times.
func ExampleNew() {
	client := &http.Client{
		Transport: retry.New(5, time.Second, nil),
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
