package gfycat

import (
	"fmt"
	"github.com/robertkozin/rvdl/internal/util"
	"regexp"
)

// TODO: Replace with yt-dlp gfycat pattern
var gifNamePattern = regexp.MustCompile(`https://\S+/([a-zA-Z]+)`)

func MatchGifName(gifUrl string) string {
	return util.FirstSubmatch(gifNamePattern, gifUrl)
}

func GetMp4Url(gifName string) string {
	return fmt.Sprintf(`https://giant.gfycat.com/%s.mp4`, gifName)
}
