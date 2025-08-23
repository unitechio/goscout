package goscout

import (
	"testing"

	"golang.org/x/crypto/ssh"
)

// mockSSHClient tạo client SSH fake (ở đây chưa có server test thật)
// Thực tế bạn có thể spin lên container SSH để test.
func mockSSHClient(t *testing.T) *ssh.Client {
	// ⚠️ Trong CI/CD bạn có thể skip test này nếu không có server
	t.Skip("Integration test: cần SSH server để chạy")
	return nil
}

func TestNewSSHClient_Invalid(t *testing.T) {
	server := Server{
		Addr:     "127.0.0.1:12345", // cổng không mở
		User:     "fakeuser",
		Password: "fakepass",
	}
	_, err := NewSSHClient(server)
	if err == nil {
		t.Error("expected error when connecting to invalid SSH server, got nil")
	}
}

func TestRunCommand_Error(t *testing.T) {
	client := mockSSHClient(t)
	if client == nil {
		return
	}
	_, err := RunCommand(client, "exit 1")
	if err == nil {
		t.Error("expected error when running bad command, got nil")
	}
}

func TestDetectOS_Skip(t *testing.T) {
	client := mockSSHClient(t)
	if client == nil {
		return
	}
	osinfo, err := DetectOS(client)
	if err != nil {
		t.Fatalf("DetectOS failed: %v", err)
	}
	if osinfo.Name == "" {
		t.Error("OS name should not be empty")
	}
}

func TestGetResourceInfo_Skip(t *testing.T) {
	client := mockSSHClient(t)
	if client == nil {
		return
	}
	info, err := GetResourceInfo(client)
	if err != nil {
		t.Fatalf("GetResourceInfo failed: %v", err)
	}
	if info.CPU == "" {
		t.Error("CPU info should not be empty")
	}
}

func TestIntegrationWithSSHServer(t *testing.T) {
	server := Server{
		Addr:     "127.0.0.1:2222",
		User:     "testuser",
		Password: "testpass",
	}
	client, err := NewSSHClient(server)
	if err != nil {
		t.Skipf("Cannot connect to SSH server: %v", err)
	}
	defer client.Close()

	osinfo, err := DetectOS(client)
	if err != nil {
		t.Fatalf("DetectOS failed: %v", err)
	}
	t.Logf("Detected OS: %s %s", osinfo.Name, osinfo.Version)

	info, err := GetResourceInfo(client)
	if err != nil {
		t.Fatalf("GetResourceInfo failed: %v", err)
	}
	if info.CPU == "" || info.RAM == "" || info.Disk == "" {
		t.Error("Resource info should not be empty")
	}
}
