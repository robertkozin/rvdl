package rvdl

import "errors"

var ErrIdNotFound = errors.New("can't match an ID to that url")

var FfmpegError = errors.New("ffmpeg error")
