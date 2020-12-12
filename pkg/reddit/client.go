package reddit

import (
	"errors"
	"github.com/robertkozin/rvdl/pkg/util"
	"github.com/segmentio/encoding/json"
	"io"
	"net/http"
)

type Client interface {
	Get(url string, data V, response interface{}) error
	Post(url string, data V, response interface{}) error
}

func do(method string, url string, headers http.Header, body io.Reader, response interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header = headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New(util.ReadAllString(res.Body))
	}

	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return err
	}

	return nil
}
