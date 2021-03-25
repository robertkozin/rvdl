package server

import (
	"fmt"
	"github.com/robertkozin/rvdl/pkg/rvdl"
	"github.com/robertkozin/rvdl/pkg/util"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var IsDev = util.EnvBool("RVDL_IS_DEV", true)
var DefaultDomain = util.IifString(IsDev, "rvdl-dev.com", "rvdl.com")
var DefaultShortDomain = util.IifString(IsDev, "rvdl-dev.it", "rvdl.it")

var Domain = util.EnvString("RVDL_DOMAIN", DefaultDomain)
var ShortDomain = util.EnvString("RVDL_SHORT_DOMAIN", DefaultShortDomain)

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
		u.Host = strings.Replace(u.Host, Domain, "reddit.com", 1)
		u.Host = strings.Replace(u.Host, ShortDomain, "redd.it", 1)
		return u.String()
	}
}

func handleRvdl(res http.ResponseWriter, req *http.Request) {

	u := ProcessUrl(*req.URL)
	// TODO: Temp fix. Think about the proper solution a little more.
	// I should probably store the encoded permalink instead.
	reqUrl := util.UrlRawString(req.URL)

	id := rvdl.FindIdCache(u)
	if id.IdType == rvdl.VideoIdNone {
		fmt.Println("404")
		//res.Header().Set("Cache-Control", "public, max-age=2678400")
		res.Header().Set("Cache-Control", "public, max-age=604800")
		http.NotFound(res, req)
		return
	}

	info, err := rvdl.InfoFromIdCache(id)
	if err != nil {
		fmt.Println(err)
		res.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(res, req, "./web/static/500_server_error.mp4")
		return
	}

	if info.VideoType == rvdl.VideoTypeNone {
		fmt.Println("404")
		res.Header().Set("Cache-Control", "public, max-age=604800")
		http.ServeFile(res, req, "./web/static/404_video_not_found.mp4")
		return
	}

	if info.Permalink != reqUrl {
		res.Header().Set("Cache-Control", "public, max-age=604800")
		http.Redirect(res, req, info.Permalink, http.StatusFound)
		// TODO: Preserve download query
		// TODO: maybe start download
		return
	}

	filename, err := rvdl.DownloadCache(info)
	if err != nil {
		fmt.Println(err)
		res.Header().Set("Cache-Control", "public, max-age=86400")
		// TODO: figure out error types
		http.ServeFile(res, req, "./web/static/500_server_error.mp4")
		return
	}

	res.Header().Set("Video-Found", "?1")
	res.Header().Set("Cache-Control", "public, max-age=31536000")
	res.Header().Set("Content-Type", "video/mp4")

	_ = os.Chtimes(filename, lastModifiedTime, lastModifiedTime)

	http.ServeFile(res, req, filename)

	return
}
