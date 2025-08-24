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

// Extended network checking with multiple metrics
type NetworkInfo struct {
	IsConnected   bool
	Latency       string
	PublicIP      string
	DNSResolution bool
}

// CheckNetworkExtended performs comprehensive network connectivity checks
func CheckNetworkExtended(client *ssh.Client) (*NetworkInfo, error) {
	info := &NetworkInfo{
		IsConnected:   false,
		Latency:       "N/A",
		PublicIP:      "N/A",
		DNSResolution: false,
	}

	// Check 1: Basic connectivity with ping
	pingCmd := "ping -c 1 8.8.8.8 2>&1 | grep -q '0% packet loss' && echo 'connected' || echo 'disconnected'"
	pingResult, _ := RunCommand(client, pingCmd)
	if strings.Contains(pingResult, "connected") {
		info.IsConnected = true

		// Measure latency
		latencyCmd := "ping -c 3 8.8.8.8 | tail -1 | awk '{print $4}' | cut -d '/' -f 2"
		latency, _ := RunCommand(client, latencyCmd)
		if latency != "" {
			info.Latency = strings.TrimSpace(latency) + " ms"
		}
	}

	// Check 2: DNS resolution
	dnsCmd := "nslookup google.com 2>&1 | grep -q 'Address' && echo 'success' || echo 'failed'"
	dnsResult, _ := RunCommand(client, dnsCmd)
	info.DNSResolution = strings.Contains(dnsResult, "success")

	// Check 3: Get public IP
	ipCmd := "curl -s https://api.ipify.org"
	publicIP, _ := RunCommand(client, ipCmd)
	if publicIP != "" {
		info.PublicIP = strings.TrimSpace(publicIP)
	}

	return info, nil
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
