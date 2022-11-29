package rvdl

import (
	"net/http"
	"regexp"
)

var gifPattern = regexp.MustCompile(`(?:preview|i).redd\.it/([0-9a-z]+)\.gif`)

func GetRedditVideoFromRedditGif(gif string) (*RedditVideo, error) {
	id := FirstSubmatch(gif, gifPattern)
	if id == "" {
		return nil, nil
	}
	vid, hit := LookupVideo()
	if hit {
		return vid, nil
	}

	videoUrl := "https://i.redd.it/" + id + ".gif"

	// TODO: Have my own Http client
	if res, err := http.Head(videoUrl); err != nil || res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "image/gif" {
		return vid, nil
	}

	vid.Permalink = "https://i." + SHORT_DOMAIN + "/" + id + ".gif"
	vid.VideoType = RedditVideoGif
	vid.VideoUrl = videoUrl

	return vid, nil
}
