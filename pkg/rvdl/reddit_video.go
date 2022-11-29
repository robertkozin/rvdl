package rvdl

import "net/http"

type RedditVideo struct {
	Id     string `json:"id"`
	IdType string `json:"id_type"`

	Permalink string `json:"permalink"`

	VideoType string `json:"video_type"`
	VideoUrl  string `json:"video_url"`
	AudioUrl  string `json:"audio_url"`
}

const (
	RedditVideoMp4  = "mp4"
	RedditVideoGif  = "gif"
	RedditVideoNone = "none"
)

func GetRedditVideoFromRedditVideoUrl(videoUrl string) (*RedditVideo, error) {
	info := &RedditVideo{
		RedditId:  videoId,
		VideoType: RedditVideoNone,
	}

	videoUrl, audioUrl := BestMp4VideoAndAudioFromRedditMpd(http.DefaultClient, "https://v.redd.it/"+videoId.Id+"/DASHPlaylist.mpd")
	if videoUrl == "" {
		return info, nil
	}

	info.Permalink = "https://v." + SHORT_DOMAIN + "/" + videoId.Id + ".mp4"
	info.VideoType = RedditVideoMp4
	info.VideoUrl = videoUrl
	info.AudioUrl = audioUrl

	return info, nil
}
