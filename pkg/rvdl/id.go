package rvdl

import (
	"regexp"
	"sync"
)

var matches = []match{
	makeMatch(`v.redd\.it/([0-9a-z]+)`, VideoIdRedditVideo),
	makeMatch(`(preview|i).redd\.it/([0-9a-z]+)`, VideoIdRedditGif),
	makeMatch(`reddit\.com/(?:r|u|user)/\S{2,}/comments/([0-9a-z]+)`, VideoIdRedditPost),
	makeMatch(`reddit\.com/comments/([0-9a-z]+)`, VideoIdRedditPost),
	makeMatch(`reddit\.com/([0-9a-z]+)`, VideoIdRedditPost),
	makeMatch(`redd\.it/([0-9a-z]+)`, VideoIdRedditPost),
}

type match struct {
	pattern *regexp.Regexp
	idType  string
}

func makeMatch(pattern string, idType string) match {
	return match{regexp.MustCompile(pattern), idType}
}

var idLock sync.Mutex //TODO: Use multi locker

func FindIdCache(url string) *VideoId {
	var id *VideoId

	if urlToInfoCache.Get(url, &id); id != nil {
		return id
	}

	idLock.Lock()
	defer idLock.Unlock()

	if urlToInfoCache.Get(url, &id); id != nil {
		return id
	}

	id = FindId(url)

	urlToInfoCache.Put(url, id)

	return id
}

func FindId(url string) *VideoId {
	for _, match := range matches {
		if submatch := match.pattern.FindStringSubmatch(url); submatch != nil {
			return &VideoId{Id: submatch[1], IdType: match.idType}
		}
	}
	return &VideoId{IdType: VideoIdNone}
}
