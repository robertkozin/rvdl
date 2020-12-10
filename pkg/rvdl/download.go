package rvdl

import (
	"bytes"
	"fmt"
	"os/exec"
	"rvdl/pkg/util"
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

func Download(info *VideoInfo) (string, error) {
	switch info.VideoType {
	case VideoTypeMp4:
		return DownloadMp4(info)
	case VideoTypeDash:
		return DownloadDash(info)
	case VideoTypeGif:
		return DownloadGif(info)
	default:
		return "", nil
	}
}

func DownloadDash(info *VideoInfo) (string, error) {
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

func DownloadMp4(info *VideoInfo) (string, error) {
	filepath := info.Filepath()

	err := Ffmpeg(
		"-i", info.VideoUrl,
		"-codec", "copy",
		filepath,
	)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func DownloadGif(info *VideoInfo) (string, error) {
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
