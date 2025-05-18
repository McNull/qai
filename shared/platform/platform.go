// filepath: /home/null/projects/go/github.com/mcnull/qai/cmd/platform.go
package platform

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// getPlatform returns platform information as a string
func GetInfo() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxInfo()
	case "windows":
		return getWindowsInfo()
	case "darwin":
		return getDarwinInfo()
	default:
		return runtime.GOOS, nil
	}
}

func getCmdOutput(cmd string, args ...string) (string, error) {
	// run the command
	cmdOut, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s", err)
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
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "ID":
			id = strings.Trim(kv[1], "\"")
		case "VERSION_ID":
			version = strings.Trim(kv[1], "\"")
		}
	}

	if id == "" || version == "" {
		return "linux", nil
	}

	return fmt.Sprintf("linux %s %s", id, version), nil
}

func getWindowsInfo() (string, error) {
	cmdOut, err := getCmdOutput("cmd", "/c", "ver")
	if err != nil {
		return "windows", nil
	}

	version := strings.TrimSpace(cmdOut)

	return version, nil
}

func getDarwinInfo() (string, error) {
	// grab the version from sw_vers
	cmdOut, err := getCmdOutput("sw_vers", "-productVersion")
	if err != nil {
		return "macos", nil
	}

	version := strings.TrimSpace(cmdOut)

	if version == "" {
		return "macos", nil
	}

	return fmt.Sprintf("macos %s", version), nil
}
