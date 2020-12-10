package reddit

import (
	"encoding/base64"
	"net/http"
	"time"
)

var privateEndpoint = "https://oauth.reddit.com"
var refreshEndpoint = "https://www.reddit.com/api/v1/access_token"

type PrivateClient struct {
	clientId     string
	clientSecret string
	Username     string
	password     string

	UserAgent string

	refreshError error
	accessToken  string
}

func NewPrivateClient(clientId string, clientSecret string, username string, password string, userAgent string) (Client, error) {
	client := &PrivateClient{clientId: clientId, clientSecret: clientSecret, Username: username, password: password, UserAgent: userAgent}

	client.refresh()
	if client.refreshError != nil {
		return nil, client.refreshError
	}

	go func() {
		for range time.Tick((1 * time.Hour) - (5 * time.Second)) {
			client.refresh()
		}
	}()

	return client, nil
}

func (c *PrivateClient) refresh() {
	data := V{
		"grant_type": "password",
		"username":   c.Username,
		"password":   c.password,
	}

	headers := http.Header{
		"Authorization": {basicAuth(c.clientId, c.clientSecret)},
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"User-Agent":    {c.UserAgent},
	}

	var res accessTokenResponse

	err := do(http.MethodPost, refreshEndpoint, headers, data.EncodeReader(), &res)
	if err != nil {
		c.refreshError = err
		return
	}

	c.accessToken = res.AccessToken
	c.refreshError = nil
	return
}

func (c *PrivateClient) Get(url string, data V, response interface{}) error {
	if c.refreshError != nil {
		return c.refreshError
	}

	data["raw_json"] = "1"

	headers := http.Header{
		"User-Agent":    {c.UserAgent},
		"Authorization": {"bearer " + c.accessToken},
	}

	return do(http.MethodGet, privateEndpoint+url+data.Query(), headers, nil, response)
}

func (c *PrivateClient) Post(url string, data V, response interface{}) error {
	if c.refreshError != nil {
		return c.refreshError
	}

	data["raw_json"] = "1"

	headers := http.Header{
		"User-Agent":    {c.UserAgent},
		"Authorization": {"bearer " + c.accessToken},
		"Content-Type":  {"application/x-www-form-urlencoded"},
	}

	return do(http.MethodPost, privateEndpoint+url, headers, data.EncodeReader(), response)
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func basicAuth(username string, password string) string {
	auth := username + ":" + password
	auth = base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + auth
}
