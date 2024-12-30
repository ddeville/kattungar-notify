package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func MakeRequest(method, path string, body []byte, authToken *string) *http.Response {
	root := GetServerUrl()
	if root == "" {
		log.Fatalln("Missing 'server_url' in config (or 'KATTUNGAR_NOTIFY_SERVER_URL' environment variable)")
	}

	reqUrl, err := url.JoinPath(root, path)
	if err != nil {
		log.Fatalf("URL is malformed root = %s, path = %s\n", root, path)
	}

	var buf io.Reader
	if body != nil {
		buf = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, reqUrl, buf)
	if err != nil {
		log.Fatalln(err)
	}

	if authToken != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *authToken))
	}
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
