package rvdl

import (
	"errors"
	"net/http"
	"strings"
	"sync"
)

var idInfoLock sync.Mutex

func InfoFromIdCache(id *VideoId) (*VideoInfo, error) {
	var info *VideoInfo
	idString := id.IdString()

	if idToInfoCache.Get(idString, &info); info != nil {
		return info, nil
	}

	idInfoLock.Lock()
	defer idInfoLock.Unlock()

	if idToInfoCache.Get(idString, &info); info != nil { //TODO: Fast get only
		return info, nil
	}

	var err error
	info, err = InfoFromId(id)
	if err != nil {
		return nil, err
	}

	idToInfoCache.Put(idString, info)
	urlToInfoCache.Put(info.Permalink, info)

	return info, nil
}

func InfoFromId(id *VideoId) (*VideoInfo, error) {
	switch id.IdType {
	case VideoIdRedditPost:
		return FromPost(id)
	case VideoIdRedditVideo:
		return FromVideo(id)
	case VideoIdRedditGif:
		return FromGif(id)
	default:
		return &VideoInfo{
			VideoId: id,
		}, nil
	}
}

func FromPost(id *VideoId) (*VideoInfo, error) {
	p, err := GetPostFromId(id)
	if err != nil {
		return nil, err
	}

	info := &VideoInfo{
		VideoId:   id,
		Permalink: "https://www." + Domain + strings.TrimSuffix(p.Permalink, "/") + ".mp4",
		VideoType: VideoTypeNone,
	}

	if p.Media.RedditVideo.DashUrl != "" {
		info.VideoType = VideoTypeDash
		info.VideoUrl, info.AudioUrl = VideoAudioFromMpd(p.Media.RedditVideo.DashUrl)
	} else if len(p.CrossPostParentList) > 0 && p.CrossPostParentList[0].Media.RedditVideo.DashUrl != "" {
		info.VideoType = VideoTypeDash
		info.VideoUrl, info.AudioUrl = VideoAudioFromMpd(p.CrossPostParentList[0].Media.RedditVideo.DashUrl)
	} else if p.Preview.RedditVideoPreview.DashUrl != "" {
		info.VideoType = VideoTypeDash
		info.VideoUrl, info.AudioUrl = VideoAudioFromMpd(p.Preview.RedditVideoPreview.DashUrl)
	} else if len(p.Preview.Images) > 0 && p.Preview.Images[0].Variants.Mp4.Source.URL != "" {
		info.VideoType = VideoTypeMp4
		info.VideoUrl = p.Preview.Images[0].Variants.Mp4.Source.URL
	} //else if id = FindId(p.URL); id.IdType != VideoIdRedditPost && id.IdType != VideoIdNone  {
	//	return InfoFromId(id)
	//}

	if info.VideoType == VideoTypeDash && info.VideoUrl == "" {
		return nil, errors.New("") // TODO
	}

	return info, nil
}

func FromVideo(id *VideoId) (*VideoInfo, error) {
	info := &VideoInfo{
		VideoId:   id,
		Permalink: "https://v." + ShortDomain + "/" + id.Id + ".mp4",
		VideoType: VideoTypeDash,
	}

	info.VideoUrl, info.AudioUrl = VideoAudioFromMpd("https://v.redd.it/" + id.Id + "/DASHPlaylist.mpd")
	if info.VideoUrl == "" {
		return nil, errors.New("") //TODO
	}

	return info, nil
}

func FromGif(id *VideoId) (*VideoInfo, error) {
	info := &VideoInfo{
		VideoId:   id,
		Permalink: "https://i." + ShortDomain + "/" + id.Id + ".mp4",
		VideoType: VideoTypeGif,
		VideoUrl:  "https://i.redd.it/" + id.Id + ".gif",
	}

	if res, err := http.Head(info.VideoUrl); err != nil || res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "image/gif" {
		return nil, errors.New("") //TODO
	}

	return info, nil
}
