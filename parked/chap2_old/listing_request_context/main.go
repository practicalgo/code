// Listing 2.1
package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	req, err, cancel := createHTTPRequest("GET", url, nil)
	defer cancel()
	if err != nil {
		return nil, err
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func createHTTPClient() *http.Client {
	client := http.Client{Timeout: 15 * time.Second}
	return &client
}

func createHTTPRequest(method string, url string, body io.Reader) (*http.Request, error, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	r, err := http.NewRequestWithContext(ctx, method, url, body)
	return r, err, cancel
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Must specify a HTTP URL to get data from")
		os.Exit(1)
	}
	client := createHTTPClient()
	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%#v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s\n", body)
}
