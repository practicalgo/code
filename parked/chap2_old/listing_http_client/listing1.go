package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// ./listing1 http HEAD https://godoc.org/net/http
var usageString = "Usage: ./listing1  HEAD URL"

type HTTPClientConfig struct {
	RequestURL string
	Method     string
}

type Config struct {
	Protocol   string
	HTTPConfig HTTPClientConfig
}

func validateInput(args []string) (Config, error) {
	c := Config{}

	protocol := args[0]
	operation := args[1]
	requestUrl := args[2]

	// validation
	if protocol != "http" {
		return c, errors.New("Invalid protocol")
	}

	if operation != "HEAD" && operation != "GET" {
		return c, errors.New("Invalid operation")
	}
	_, err := url.ParseRequestURI(requestUrl)
	if err != nil {
		return c, errors.New("Invalid URL")
	}
	c.Protocol = protocol
	c.HTTPConfig = HTTPClientConfig{
		Method:     operation,
		RequestURL: requestUrl,
	}
	return c, nil
}

func invokeHttpCmd(c Config) (string, error) {
	// Create HTTP client
	client := &http.Client{}
	// Create a request
	req, err := http.NewRequest(c.HTTPConfig.Method, c.HTTPConfig.RequestURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	fmt.Printf("%#v\n", resp.Status)
	fmt.Printf("%#v\n", resp.StatusCode)
	fmt.Printf("%#v\n", resp.Proto)

	var data string
	for k, v := range resp.Header {
		headerValue := ""
		for _, value := range v {
			headerValue += value
		}
		data += string(k) + ":" + headerValue + "\n"
	}
	return data, nil
}

func main() {

	// Parse and validate command line arguments
	if len(os.Args) == 2 && os.Args[1] == "-h" {
		fmt.Println(usageString)
		os.Exit(0)
	}
	if len(os.Args) != 4 {
		log.Fatal(usageString)
	}
	config, err := validateInput(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Run the command
	var data string
	switch config.Protocol {
	case "http":
		data, err = invokeHttpCmd(config)
	default:
		log.Fatalf("Unrecognized protocol: %s", config.Protocol)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}
