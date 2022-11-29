package rvdl

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

const (
	mimeMp4Video = "video/mp4"
	mimeMp4Audio = "audio/mp4"
)

type RedditMpd struct {
	Period struct {
		AdaptationSet []struct {
			MimeType       string `xml:"mimeType,attr"`
			Representation []struct {
				Bandwidth int    `xml:"bandwidth,attr"`
				MimeType  string `xml:"mimeType,attr"`
				BaseURL   string `xml:"BaseURL"`
			} `xml:"Representation"`
		} `xml:"AdaptationSet"`
	} `xml:"Period"`
}

func BestMp4VideoAndAudioFromRedditMpd(httpClient *http.Client, mpdUrl string) (videoUrl string, audioUrl string) {
	res, err := httpClient.Get(mpdUrl)
	if err != nil {
		return "", ""
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "application/dash+xml" {
		return "", ""
	}

	var mpd RedditMpd
	err = xml.NewDecoder(res.Body).Decode(&mpd)
	if err != nil {
		return "", ""
	}

	videoBandwidth, audioBandwidth := 0, 0

	for _, adapt := range mpd.Period.AdaptationSet {
		for _, repr := range adapt.Representation {
			if (adapt.MimeType == mimeMp4Audio || repr.MimeType == mimeMp4Audio) && repr.Bandwidth > audioBandwidth {
				audioUrl = repr.BaseURL
				audioBandwidth = repr.Bandwidth
			} else if (adapt.MimeType == mimeMp4Video || repr.MimeType == mimeMp4Video) && repr.Bandwidth > videoBandwidth {
				videoUrl = repr.BaseURL
				videoBandwidth = repr.Bandwidth
			}
		}
	}

	urlResolveRelative := func(base string, ref string) string {
		baseUrl, _ := url.Parse(base)
		refUrl, _ := url.Parse(ref)
		return baseUrl.ResolveReference(refUrl).String()
	}

	if videoUrl != "" {
		videoUrl = urlResolveRelative(mpdUrl, videoUrl)
	}
	if audioUrl != "" {
		audioUrl = urlResolveRelative(mpdUrl, audioUrl)
	}

	return videoUrl, audioUrl
}
