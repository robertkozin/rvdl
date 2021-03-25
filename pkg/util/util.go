package util

import (
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
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

func UrlRawString(u *url.URL) string {
	var buf strings.Builder

	buf.WriteString(u.Scheme)
	buf.WriteString("://")
	buf.WriteString(u.Host)
	buf.WriteString(u.Path)
	if u.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.RawQuery)
	}

	return buf.String()
}

func IifString(cond bool, a string, b string) string {
	if cond {
		return a
	} else {
		return b
	}
}
