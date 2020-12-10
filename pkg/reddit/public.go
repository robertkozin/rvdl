package reddit

import (
	"net/http"
)

var publicEndpoint = "https://www.reddit.com"

type PublicClient struct {
	UserAgent string
}

func NewPublicClient(userAgent string) *PublicClient {
	client := &PublicClient{UserAgent: userAgent}
	return client
}

func (c *PublicClient) Get(url string, data V, response interface{}) error {
	data["raw_json"] = "1"

	headers := http.Header{
		"User-Agent": {c.UserAgent},
	}

	return do(http.MethodGet, publicEndpoint+url+".json"+data.Query(), headers, nil, response)
}

func (c *PublicClient) Post(url string, data V, response interface{}) error {
	data["raw_json"] = "1"

	headers := http.Header{
		"User-Agent":   {c.UserAgent},
		"Content-Type": {"application/x-www-form-urlencoded"},
	}

	return do(http.MethodPost, privateEndpoint+url, headers, data.EncodeReader(), response)
}
