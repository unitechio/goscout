package goscout

import (
	"bytes"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Server struct {
	Addr     string // ví dụ: "192.168.1.10:22"
	User     string
	Password string
	// TODO: sau này hỗ trợ PrivateKey
}

func NewSSHClient(s Server) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", s.Addr, config)
}

func RunCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var out bytes.Buffer
	session.Stdout = &out
	if err := session.Run(cmd); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
