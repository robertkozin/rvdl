package util

import (
	"io"
	"io/ioutil"
	"net/url"
	"os"
)

func ReadAllString(reader io.Reader) string {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return ""
	}

	return string(bytes)
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func FullPath(u url.URL) string {
	if u.RawQuery != "" {
		return u.Path[1:] + "?" + u.RawPath
	} else {
		return u.Path[1:]
	}
}
