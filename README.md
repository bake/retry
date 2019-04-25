# retry

[![GoDoc](https://godoc.org/github.com/bake/retry?status.svg)](http://godoc.org/github.com/bake/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/bake/retry)](https://goreportcard.com/report/github.com/bake/retry)

Package retry is a small implementation of the `http.RoundTripper` interface
that can be found in `http.Client`. It is responsible to make HTTP requests
and can be used to cache or retry them.

By default a request will be tried again only if the response code is below
500 (Internal Server Error) and not 429 (Too Many Requests). This behaviour
can be changed with the `WithVerifier` option.

#### Examples

##### New

Create a new HTTP client that will retry requests five times and sleeps
one second between each one. Since the requested URL will always return a
status code of 500 (Internal Server Error), this function will run for five
seconds and the time it takes to call the endpoint five times.

```golang
package main

import (
	"fmt"
	"github.com/bake/retry"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
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

```

##### WithVerifier

WithVerifier can be used to modify the default behaviour that is used to
determine if a request can be repeated. In this example, all responses that
do not succeed (with a status code outside [200-300[) will be tried again.

```golang
package main

import (
	"fmt"
	"github.com/bake/retry"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
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

```
