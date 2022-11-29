package rvdl

import "regexp"

var PRIMARY_DOMAIN = "rvdl.com"
var SHORT_DOMAIN = "rvdl.it"

func FirstSubmatch(s string, patterns ...*regexp.Regexp) string {
	for _, pattern := range patterns {
		if submatch := pattern.FindStringSubmatch(s); len(submatch) >= 2 {
			return submatch[1]
		}
	}
	return ""
}

var a = []func(string) (*RedditVideo, error){
	GetRedditVideoFromRedditPostId,
	GetRedditVideoFromRedditGif,
	GetRedditVideoFromRedditVideoUrl,
}

func GetRedditVideoFromRedditUrl(redditUrl string) (vid *RedditVideo, err error) {
	for _, fn := range a {
		if vid, err = fn(redditUrl); vid != nil {
			return vid, err
		}
	}
	return nil, nil
}
