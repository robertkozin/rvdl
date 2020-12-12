package rvdl

import (
	"github.com/robertkozin/rvdl/pkg/cache"
	"github.com/robertkozin/rvdl/pkg/reddit"
	"github.com/robertkozin/rvdl/pkg/util"
)

var FfmpegPath = util.EnvString("RVDL_FFMPEG_PATH", "/bin/ffmpeg")
var VideosDir = util.EnvString("RVDL_VIDEOS_DIR", "./web/static/videos/")

var Domain = util.EnvString("RVDL_DOMAIN", "rvdl.com")
var ShortDomain = util.EnvString("RVDL_SHORT_DOMAIN", "rvdl.it")

var CacheDir = util.EnvString("RVDL_CACHE_DIR", "./cache/")

var RedditUsername = util.EnvString("RVDL_REDDIT_USERNAME", "")
var RedditPassword = util.EnvString("RVDL_REDDIT_PASSWORD", "")
var RedditClientId = util.EnvString("RVDL_REDDIT_CLIENT_ID", "")
var RedditClientSecret = util.EnvString("RVDL_REDDIT_CLIENT_SECRET", "")

var UserAgent = util.EnvString("RVDL_USER_AGENT", "rvdl")

var urlToInfoCache *cache.LRU
var idToInfoCache *cache.Cache
var redditClient reddit.Client

func Init() error {
	var err error

	if RedditClientId != "" && RedditClientSecret != "" && RedditUsername != "" && RedditPassword != "" {
		redditClient, err = reddit.NewPrivateClient(RedditClientId, RedditClientSecret, RedditUsername, RedditPassword, UserAgent)
		if err != nil {
			return err
		}
	} else {
		redditClient = reddit.NewPublicClient(UserAgent)
	}

	urlToInfoCache = cache.NewLru(100)
	idToInfoCache = cache.NewCache(100, CacheDir+"id_to_info")

	return nil
}

func Close() {
	idToInfoCache.Close()
}
