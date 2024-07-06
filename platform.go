package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func getOSVersion() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "ver")
	case "linux":
		cmd = exec.Command("lsb_release", "-d")
	case "darwin":
		cmd = exec.Command("sw_vers", "-productVersion")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	versionInfo := strings.TrimSpace(out.String())

	if runtime.GOOS == "linux" {
		// Extract description from `lsb_release -d` output
		parts := strings.SplitN(versionInfo, ":", 2)
		if len(parts) == 2 {
			versionInfo = strings.TrimSpace(parts[1])
		}
	}

	return versionInfo, nil
}

func getPlatform() (string, error) {
	version, err := getOSVersion()
	if err != nil {
		return runtime.GOOS, err
	}

	name := runtime.GOOS

	if name == "darwin" {
		name = "macOS"
	}

	return name + " " + version, nil
}
