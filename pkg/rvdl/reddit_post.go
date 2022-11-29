package rvdl

import (
	"github.com/robertkozin/rvdl/pkg/gfycat"
	"strings"
)

//makeMatch(`reddit\.com/(?:r|u|user)/\S{2,}/comments/([0-9a-z]+)`, RedditIdPost),
//makeMatch(`reddit\.com/comments/([0-9a-z]+)`, RedditIdPost),
//makeMatch(`reddit\.com/([0-9a-z]+)`, RedditIdPost),
//makeMatch(`redd\.it/([0-9a-z]+)`, RedditIdPost),

func GetRedditVideoFromRedditPostId(postUrl string) (*RedditVideo, error) {
	id := FirstSubmatch(postUrl)
	if id == "" {
		return nil, nil
	}
	vid, hit := LookupVideo()
	if hit {
		return vid, nil
	}

	// GET POST
	var p RedditPost

	vid.Permalink = "https://www." + PRIMARY_DOMAIN + strings.TrimSuffix(p.Permalink, "/") + ".mp4"

	if p.Media.OEmbed.ProviderName == "Gfycat" {
		name := gfycat.MatchGifName(p.Media.OEmbed.ThumbnailUrl)
		vid.VideoType = RedditVideoMp4
		vid.VideoUrl = gfycat.GetMp4Url(name)
		//} else if p.Media.OEmbed.ProviderName == "RedGIFs" {
		//	name := gfycat.MatchGifName(p.Media.OEmbed.ThumbnailUrl)
		// TODO: Implement native RedGIFs info, watch for video expiry and deprecated api
	} else if p.Media.RedditVideo.DashUrl != "" {
		vid.VideoType = RedditVideoMp4
		vid.VideoUrl, vid.AudioUrl = avurl(&p.Media.RedditVideo)
	} else if len(p.CrossPostParentList) > 0 && p.CrossPostParentList[0].Media.RedditVideo.DashUrl != "" {
		vid.VideoType = RedditVideoMp4
		vid.VideoUrl, vid.AudioUrl = avurl(&p.CrossPostParentList[0].Media.RedditVideo)
	} else if p.Preview.RedditVideoPreview.DashUrl != "" {
		vid.VideoType = RedditVideoMp4
		vid.VideoUrl, vid.AudioUrl = avurl(&p.Preview.RedditVideoPreview)
	} else if len(p.Preview.Images) > 0 && p.Preview.Images[0].Variants.Mp4.Source.URL != "" {
		vid.VideoType = RedditVideoMp4
		vid.VideoUrl = p.Preview.Images[0].Variants.Mp4.Source.URL
	}

	PutVideo("", vid)
	PutVideo(postUrl, vid)
	PutVideo(vid.Permalink, vid)

	return vid, nil
}

func avurl(v *RedditVideoo) (string, string) {
	return "", ""
}

type RedditPost struct {
	URL       string `json:"url"`
	Permalink string `json:"permalink"`
	Media     struct {
		RedditVideo RedditVideoo `json:"reddit_video"`
		OEmbed      struct {
			ProviderName string `json:"provider_name"` // Gfycat, RedGIFs
			ThumbnailUrl string `json:"thumbnail_url"` // https://giant.gfycat.com/{{name}}.mp4, // https://\S+/([a-zA-Z]+)
		} `json:"oembed"`
	} `json:"secure_media"`
	Preview struct {
		RedditVideoPreview RedditVideoo `json:"reddit_video_preview"`
		Images             []struct {
			Variants struct {
				Mp4 struct {
					Source struct {
						URL string `json:"url"`
					} `jsonPermalink:"source"`
				} `json:"mp4"`
			} `json:"variants"`
		} `json:"images"`
	} `json:"preview"`
	CrossPostParentList []struct {
		Media struct {
			RedditVideo RedditVideoo `json:"reddit_video"`
		} `json:"secure_media"`
	} `json:"crosspost_parent_list"`
}

type RedditVideoo struct {
	DashUrl     string `json:"dash_url"`
	FallbackUrl string `json:"fallback_url"`
	IsGif       bool   `json:"is_gif"`
}
