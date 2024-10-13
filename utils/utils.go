package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
)

func ExtractThumbnail(videoPath string, outputPath string, time string) error {
	cmd := exec.Command("ffmpeg", "-ss", time, "-i", videoPath, "-vframes", "1", outputPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error extracting thumbnail: %w", err)
	}
	return nil
}

func AsPrettyJson(input interface{}) string {
	jsonB, _ := json.MarshalIndent(input, "", "  ")
	return fmt.Sprintf("```%s```", string(jsonB))
}

func GetPrivateIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error: " + err.Error())
	}

	for _, i := range interfaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}

func GetFileHash(path string) string {
	cmd := exec.Command("sha256sum", path)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return string(output)[:64]
}
