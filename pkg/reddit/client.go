package reddit

import (
	"github.com/segmentio/encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

type Client interface {
	Get(url string, data V, response interface{}) error
	GetListing(url string, data V, response interface{}) error
	Post(url string, data V, response interface{}) error
}

func do(httpClient *http.Client, method string, url string, headers http.Header, body io.Reader, response interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header = headers

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		var redditErr *RedditError
		json.NewDecoder(res.Body).Decode(&redditErr) // TODO
		return redditErr
	}

	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return err
	}

	return nil
}

type listing struct {
	//Kind string `json:"kind"`
	Data struct {
		//Modhash  string `json:"modhash"`
		Dist     int    `json:"dist"`
		Children []struct {
			//Kind string `json:"kind"`
			Data json.RawMessage `json:"data"`
		} `json:"children"`
		//After  interface{} `json:"after"`
		//Before interface{} `json:"before"`
	} `json:"data"`
}

func getListing(client Client, url string, data V, response interface{}) error {
	var listing listing
	err := client.Get(url, data, &listing)
	if err != nil {
		return err
	}

	x := reflect.ValueOf(response).Elem()

	dist := listing.Data.Dist
	x.Set(reflect.MakeSlice(x.Type(), dist, dist))
	for i, child := range listing.Data.Children {
		_ = json.Unmarshal(child.Data, x.Index(i).Addr().Interface())
	}

	return nil
}
// {"reason": "private", "message": "Forbidden", "error": 403}
type RedditError struct {
	Code int `json:"error"`
	Message string `json:"message"`
	Reason string `json:"reason"`
}

func (e *RedditError) Error() string {
	return strconv.Itoa(e.Code) + e.Message + e.Reason
}
