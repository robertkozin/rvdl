package rvdl

import (
	"errors"
	"github.com/robertkozin/rvdl/pkg/reddit"
)

// TODO: Probably just move info info.go
func GetPostFromId(id *VideoId) (*RedditPost, error) {
	data := reddit.V{"id": "t3_" + id.Id}
	var res RedditListing

	err := redditClient.Get("/api/info", data, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Data.Children) <= 0 {
		return nil, errors.New("post not found")
	}

	return &(res.Data.Children[0].Data), nil
}
