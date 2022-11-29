package server

import (
	"github.com/robertkozin/rvdl/pkg/rvdl"
	"net/http"
	"net/url"
	"strings"
)

func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		serveHome(res, req)
	case "/favicon.ico", "/apple-touch-icon.png":
		serveFavicon(res, req)
	case "/robots.txt":
		serveRobots(res, req)
	default:
		serveRvdl(res, req)
	}
}

func serveHome(res http.ResponseWriter, req *http.Request) {

}

func serveFavicon(res http.ResponseWriter, req *http.Request) {

}

func serveRobots(res http.ResponseWriter, req *http.Request) {

}

func serveRvdl(res http.ResponseWriter, req *http.Request) {
	rvdlUrl := req.URL
	redditUrl := GetRedditUrlFromRvdlUrl(rvdlUrl)
	redditId := rvdl.FindRedditIdFromRedditUrl(redditUrl)

}

func GetRedditUrlFromRvdlUrl(rvdlUrl *url.URL) (redditUrl string) {
	if strings.HasPrefix(rvdlUrl.Path, "https://") {
		// https://rvdl.com/https://reddit.com/abc123
		return rvdlUrl.Path
	} else {

	}


	if strings.Contains(u.Path, "reddit.com/") || strings.Contains(u.Path, "redd.it/") {
		if u.RawQuery != "" {
		return u.Path + "?" + u.RawQuery
	} else {
		return u.Path
	}
	} else {
		u.Host = strings.Replace(u.Host, Domain, "reddit.com", 1)
		u.Host = strings.Replace(u.Host, ShortDomain, "redd.it", 1)
		return u.String()
	}
	}
}
