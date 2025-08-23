package goscout

import (
	"strings"

	"golang.org/x/crypto/ssh"
)

type ResourceInfo struct {
	CPU  string
	RAM  string
	Disk string
}

func GetResourceInfo(client *ssh.Client) (*ResourceInfo, error) {
	osinfo, _ := DetectOS(client)

	var cpuCmd, ramCmd, diskCmd string

	switch {
	case strings.Contains(strings.ToLower(osinfo.Name), "ubuntu"),
		strings.Contains(strings.ToLower(osinfo.Name), "debian"):
		cpuCmd = `top -bn1 | grep "Cpu(s)"`
		ramCmd = `free -h | grep Mem`
		diskCmd = `df -h --total | grep total`

	case strings.Contains(strings.ToLower(osinfo.Name), "centos"),
		strings.Contains(strings.ToLower(osinfo.Name), "fedora"),
		strings.Contains(strings.ToLower(osinfo.Name), "red hat"):
		cpuCmd = `top -bn1 | grep Cpu`
		ramCmd = `free -m | grep Mem`
		diskCmd = `df -h --total | grep total`

	default: // fallback
		cpuCmd = `uptime`
		ramCmd = `free -h`
		diskCmd = `df -h`
	}

	cpu, _ := RunCommand(client, cpuCmd)
	ram, _ := RunCommand(client, ramCmd)
	disk, _ := RunCommand(client, diskCmd)

	return &ResourceInfo{
		CPU:  cpu,
		RAM:  ram,
		Disk: disk,
	}, nil
}
