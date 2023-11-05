package net

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

func MakeRequest(method, url string, body any) (*http.Response, error) {
	client := retryablehttp.NewClient()
	client.RetryMax = 3

	req, err := retryablehttp.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
