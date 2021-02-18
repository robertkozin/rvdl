package reddit

import (
	"net/http"
)

var publicEndpoint = "https://www.reddit.com"

type PublicClient struct {
	UserAgent string
	httpClient *http.Client
}

func NewPublicClient(userAgent string) Client {
	client := &PublicClient{
		UserAgent: userAgent,
		httpClient: &http.Client{},
	}
	return client
}

func (c *PublicClient) Get(url string, data V, response interface{}) error {
	data["raw_json"] = "1"

	headers := http.Header{
		"User-Agent": {c.UserAgent},
	}

	return do(c.httpClient, http.MethodGet, publicEndpoint+url+".json"+data.Query(), headers, nil, response)
}

func (c *PublicClient) GetListing(url string, data V, response interface{}) error {
	return getListing(c, url, data, response)
}

func (c *PublicClient) Post(url string, data V, response interface{}) error {
	return nil
}
