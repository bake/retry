# retry

[![GoDoc](https://godoc.org/github.com/bakerolls/retry?status.svg)](https://godoc.org/github.com/bakerolls/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/bakerolls/retry)](https://goreportcard.com/report/github.com/bakerolls/retry)

`retry` is a small implementation of the [`http.RoundTripper`](https://golang.org/pkg/net/http/#RoundTripper) interface, that can be found in [`http.Client`](https://golang.org/pkg/net/http/#Client). It is responsible to make HTTP requests and can be used to cache or retry them.

This package will only retry a request if it did not return an error and the status code is outside of [200-300[.

```go
func main() {
	// Create a new http.Client that will retry requests five times and sleeps
	// one second between each one.
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
```
