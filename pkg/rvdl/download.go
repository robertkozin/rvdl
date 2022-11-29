package rvdl

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"
)

var downloadLock sync.Mutex

func DownloadCache(info *VideoInfo) (string, error) {
	filepath := info.Filepath()
	if util.FileExists(filepath) {
		return filepath, nil
	}

	downloadLock.Lock()
	defer downloadLock.Unlock()

	if util.FileExists(filepath) {
		return filepath, nil
	}

	_, err := Download(info)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func Download(vid *RedditVideo) (string, error) {
	switch vid.VideoType {
	case RedditVideoMp4:
		return DownloadMp4(vid)
	case RedditVideoGif:
		return DownloadGif(vid)
	default:
		return "", nil // TODO: error
	}
}

func DownloadDash(vid *RedditVideo) (string, error) {
	filepath := info.Filepath()
	args := []string{
		"-i", info.VideoUrl,
		"-codec", "copy",
		filepath,
	}
	if info.AudioUrl != "" {
		args = append([]string{"-i", info.AudioUrl}, args...)
	}

	err := Ffmpeg(args...)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func DownloadGif(vid *RedditVideo) (string, error) {
	filepath := info.Filepath()

	err := Ffmpeg(
		"-i", info.VideoUrl,
		"-pix_fmt", "yuv420p",
		"-c:v", "libx264",
		"-filter:v", "crop='floor(in_w/2)*2:floor(in_h/2)*2'",
		filepath,
	)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func Ffmpeg(arg ...string) error {
	arg = append(
		arg,
		"-hide_banner",
		"-user_agent", "rvdl",
		"-loglevel", "error",
		"-f", "mp4",
		"-movflags", "faststart",
		"-y",
	)
	cmd := exec.Command(FfmpegPath, arg...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%w: %s > %s", FfmpegError, cmd.String(), buf.String())
	}

	return nil
}
