package rvdl

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

var mimeVideo = "video/mp4"
var mimeAudio = "audio/mp4"

type Mpd struct {
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

func VideoAudioFromMpd(url string) (string, string) {
	res, err := http.Get(url)
	if err != nil {
		return "", ""
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "application/dash+xml" {
		return "", ""
	}

	var mpd Mpd
	err = xml.NewDecoder(res.Body).Decode(&mpd)
	if err != nil {
		return "", ""
	}

	videoUrl, videoBandwidth, audioUrl, audioBandwidth := "", 0, "", 0

	for _, adapt := range mpd.Period.AdaptationSet {
		for _, repr := range adapt.Representation {
			if (adapt.MimeType == mimeAudio || repr.MimeType == mimeAudio) && repr.Bandwidth > audioBandwidth {
				audioUrl = repr.BaseURL
				audioBandwidth = repr.Bandwidth
			} else if (adapt.MimeType == mimeVideo || repr.MimeType == mimeVideo) && repr.Bandwidth > videoBandwidth {
				videoUrl = repr.BaseURL
				videoBandwidth = repr.Bandwidth
			}
		}
	}

	if videoUrl != "" {
		videoUrl = UrlResolveRelative(url, videoUrl)
	}
	if audioUrl != "" {
		audioUrl = UrlResolveRelative(url, audioUrl)
	}

	return videoUrl, audioUrl
}

func UrlResolveRelative(base string, ref string) string {
	baseUrl, _ := url.Parse(base)
	refUrl, _ := url.Parse(ref)
	return baseUrl.ResolveReference(refUrl).String()
}
