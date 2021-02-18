package reddit

type RssClient struct {
	Client
	token string
	username string
}

func NewRssClient(token string, username string, userAgent string) Client {
	return &RssClient{token: token, username: username, Client: NewPublicClient(userAgent)}
}

func (c *RssClient) Get(url string, data V, response interface{}) error {
	data["feed"] = c.token
	data["user"] = c.username
	return c.Client.Get(url, data, response)
}

func (c *RssClient) GetListing(url string, data V, response interface{}) error {
	data["feed"] = c.token
	data["user"] = c.username
	return c.Client.Get(url, data, response)
}

func (c *RssClient) Post(url string, data V, response interface{}) error {
	return nil
}

