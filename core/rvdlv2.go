package rvdl

import (
	"github.com/moby/locker"
	"github.com/robertkozin/rvdl/pkg/cache"
	"github.com/robertkozin/rvdl/pkg/reddit"
)

type Rvdl struct {
	config   RvdlConfig
	locker   *locker.Locker
	urlCache *cache.LRU
	idCache  *cache.Cache
	reddit   *reddit.Client
}

type RvdlConfig struct {
	FfmpegPath string
	VideosDir string
	CacheDir string
}

func NewRvdl(config RvdlConfig, redditClient *reddit.Client, ) (Rvdl, error) {
	r := Rvdl{
		config: config,
		reddit: redditClient,
		urlCache: cache.NewLru(100),
		idCache: cache.NewCache(100, config.CacheDir),
	}

	return r, nil
}

func (rvdl *Rvdl) MatchId() {

}

func (rvdl *Rvdl) MatchInfo() {

}

func (rvdl *Rvdl) Download() {

}

