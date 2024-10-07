package utils

import (
	"fmt"
	"os/exec"
)

func ExtractThumbnail(videoPath string, outputPath string, time string) error {
	cmd := exec.Command("ffmpeg", "-ss", time, "-i", videoPath, "-vframes", "1", outputPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error extracting thumbnail: %w", err)
	}
	return nil
}
