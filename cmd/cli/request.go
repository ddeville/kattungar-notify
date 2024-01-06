package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func makeRequest(method, url string, body []byte, authToken *string) *http.Response {
	if authToken == nil {
		apiKey, has := os.LookupEnv("KATTUNGAR_NOTIFY_API_KEY")
		if !has {
			log.Fatalln("Missing KATTUNGAR_NOTIFY_API_KEY environment variable")
		}
		authToken = &apiKey
	}

	var buf io.Reader
	if body != nil {
		buf = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *authToken))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode == 401 {
		log.Fatalln("API key is incorrect")
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Fatalf("There was an unknown server error: %v", res.Status)
	}

	return res
}
