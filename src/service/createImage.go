package service

import (
	"os/exec"
	"strings"
)

func FfmpegCreateCover(videoFileName string) string {
	coverName := videoFileName[:strings.Index(videoFileName, ".")] + ".jpg"

	cmdArguments := []string{"-i", videoFileName, "-vf", "select=eq(n\\,1)",
		"-vframes", "1", coverName}
	cmd := exec.Command("ffmpeg", cmdArguments...)
	cmd.Run()
	return coverName
}
