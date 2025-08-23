package goscout

import (
	"strings"

	"golang.org/x/crypto/ssh"
)

type OSInfo struct {
	Name    string
	Version string
	IDLike  []string
}

// DetectOS phát hiện hệ điều hành từ /etc/os-release
func DetectOS(client *ssh.Client) (*OSInfo, error) {
	out, err := RunCommand(client, `cat /etc/os-release`)
	if err != nil {
		return nil, err
	}

	info := &OSInfo{}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "NAME=") {
			info.Name = strings.Trim(line[5:], `"`)
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			info.Version = strings.Trim(line[11:], `"`)
		} else if strings.HasPrefix(line, "ID_LIKE=") {
			val := strings.Trim(line[8:], `"`)
			info.IDLike = strings.Split(val, " ")
		}
	}

	return info, nil
}
