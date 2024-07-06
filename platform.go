package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func getCmdOutput(cmd string, args ...string) (string, error) {
	// run the command
	cmdOut, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "", fmt.Errorf("Failed to run command: %s", err)
	}

	return string(cmdOut), nil
}

func getLinuxInfo() (string, error) {
	// grab the ID and VERSION_ID from /etc/os-release
	cmdOut, err := getCmdOutput("cat", "/etc/os-release")
	if err != nil {
		return "", err
	}

	// parse the output
	var id, version string
	for _, line := range strings.Split(cmdOut, "\n") {
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		kv := strings.Split(line, "=")
		switch kv[0] {
		case "ID":
			id = kv[1]
		case "VERSION_ID":
			version = kv[1]
		}
	}

	if id == "" || version == "" {
		return "", fmt.Errorf("Failed to parse /etc/os-release")
	}

	return fmt.Sprintf("linux %s %s", id, version), nil
}

func getWindowsInfo() (string, error) {
	cmdOut, err := getCmdOutput("cmd", "/c", "ver")
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(cmdOut)

	return version, nil
}

func getDarwinInfo() (string, error) {
	// grab the version from sw_vers
	cmdOut, err := getCmdOutput("sw_vers", "-productVersion")
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(cmdOut)

	if version == "" {
		return "", fmt.Errorf("Failed to parse sw_vers")
	}

	return fmt.Sprintf("macos %s", version), nil
}

func getPlatform() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxInfo()
	case "windows":
		return getWindowsInfo()
	case "darwin":
		return getDarwinInfo()
	default:
		return "", fmt.Errorf("Unsupported platform")
	}
}
