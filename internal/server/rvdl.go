package server

import (
	"net/http"
	"net/url"
	"os"
	"rvdl/pkg/rvdl"
	"strings"
	"time"
)

var lastModified = "Tue, 01 Dec 2020 00:00:00 GMT"
var lastModifiedTime, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", lastModified)

func ProcessUrl(u url.URL) string {
	if strings.Contains(u.Path, "reddit.com/") || strings.Contains(u.Path, "redd.it/") {
		if u.RawQuery != "" {
			return u.Path + "?" + u.RawQuery
		} else {
			return u.Path
		}
	} else {
		u.Host = strings.Replace(u.Host, "rvdl.com", "reddit.com", 1)
		u.Host = strings.Replace(u.Host, "rvdl.it", "redd.it", 1)
		return u.String()
	}
}

func handleRvdl(res http.ResponseWriter, req *http.Request) {

	u := ProcessUrl(*req.URL)
	// TODO: Temp fix. Think about the proper solution a little more.
	// I should probably store the encoded permalink instead.
	reqUrl := "https://" + req.URL.Host + req.URL.Path
	if req.URL.RawQuery != "" {
		reqUrl += "?" + req.URL.RawQuery
	}

	id := rvdl.FindIdCache(u)
	if id.IdType == rvdl.VideoIdNone {
		http.NotFound(res, req)
		return
	}

	info, err := rvdl.InfoFromIdCache(id)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	if info.VideoType == rvdl.VideoTypeNone {
		http.NotFound(res, req)
		return
	}

	if info.Permalink != reqUrl {
		redirect(res, req, info.Permalink, http.StatusFound)
		// TODO: maybe start download
		return
	}

	filename, err := rvdl.DownloadCache(info)
	if err != nil {
		// TODO: figure out error types
		http.Error(res, err.Error(), 500)
		return
	}

	res.Header().Set("Cache-Control", "public, max-age=31536000")
	res.Header().Set("Content-Type", "video/mp4")

	_ = os.Chtimes(filename, lastModifiedTime, lastModifiedTime)

	http.ServeFile(res, req, filename)

	return
}
