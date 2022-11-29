package rvdl

// TODO: Probably better to put these into their respective files

type VideoId struct {
	Id     string `json:"id"`
	IdType string `json:"id_type"`
}

type VideoInfo struct {
	*VideoId

	Permalink string `json:"permalink"`

	VideoType string `json:"video_type"`
	VideoUrl  string `json:"video_url"`
	AudioUrl  string `json:"audio_url"`
}

func (id *VideoId) IdString() string {
	return id.IdType + ":" + id.Id
}

func (id *VideoId) Filename() string {
	return id.IdType + "_" + id.Id + ".mp4"
}

func (id *VideoId) Filepath() string {
	return VideosDir + id.Filename()
}

const (
	VideoIdRedditPost  = "post"
	VideoIdRedditVideo = "video"
	VideoIdRedditGif   = "gif"
	VideoIdNone        = "none"
)

const (
	VideoTypeMp4  = "mp4"
	VideoTypeDash = "dash"
	VideoTypeGif  = "gif"
	VideoTypeNone = "none"
)

type RedditListing struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditPost struct {
	URL       string `json:"url"`
	Permalink string `json:"permalink"`
	Media     struct {
		RedditVideo RedditVideo `json:"reddit_video"`
		OEmbed      struct {
			ProviderName string `json:"provider_name"` // Gfycat, RedGIFs
			ThumbnailUrl string `json:"thumbnail_url"` // https://giant.gfycat.com/{{name}}.mp4, // https://\S+/([a-zA-Z]+)
		} `json:"oembed"`
	} `json:"secure_media"`
	Preview struct {
		RedditVideoPreview RedditVideo `json:"reddit_video_preview"`
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
			RedditVideo RedditVideo `json:"reddit_video"`
		} `json:"secure_media"`
	} `json:"crosspost_parent_list"`
}

type RedditVideo struct {
	DashUrl string `json:"dash_url"`
}
