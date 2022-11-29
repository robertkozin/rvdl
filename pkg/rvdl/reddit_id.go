package rvdl

import "regexp"

type RedditId struct {
	Id     string `json:"id"`
	IdType string `json:"id_type"`
}

const (
	RedditIdPost  = "post"
	RedditIdVideo = "video"
	RedditIdGif   = "gif"
	RedditIdNone  = "none"
)

var matches = []match{
	makeMatch(`v.redd\.it/([0-9a-z]+)`, RedditIdVideo),
	makeMatch(`(?:preview|i).redd\.it/([0-9a-z]+)\.gif`, RedditIdGif),
	makeMatch(`reddit\.com/(?:r|u|user)/\S{2,}/comments/([0-9a-z]+)`, RedditIdPost),
	makeMatch(`reddit\.com/comments/([0-9a-z]+)`, RedditIdPost),
	makeMatch(`reddit\.com/([0-9a-z]+)`, RedditIdPost),
	makeMatch(`redd\.it/([0-9a-z]+)`, RedditIdPost),
}

type match struct {
	pattern      *regexp.Regexp
	redditIdType string
}

func makeMatch(pattern string, redditIdType string) match {
	return match{regexp.MustCompile(pattern), redditIdType}
}

func FindRedditIdFromRedditUrl(redditUrl string) *RedditId {
	for _, match := range matches {
		if submatch := match.pattern.FindStringSubmatch(redditUrl); submatch != nil {
			return &RedditId{Id: submatch[1], IdType: match.redditIdType}
		}
	}
	return &RedditId{IdType: RedditIdNone}
}
