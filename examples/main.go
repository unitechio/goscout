package main

import (
	"fmt"
	"log"

	"github.com/unitechio/goscout"
)

func main() {
	server := goscout.Server{
		Addr:     "127.0.0.1:2222", // hoặc IP server thật
		User:     "testuser",
		Password: "testpass",
	}

	client, err := goscout.NewSSHClient(server)
	if err != nil {
		log.Fatal("SSH connect failed:", err)
	}
	defer client.Close()

	info, err := goscout.GetResourceInfo(client)
	if err != nil {
		log.Fatal("GetResourceInfo failed:", err)
	}

	fmt.Println("CPU:", info.CPU)
	fmt.Println("RAM:", info.RAM)
	fmt.Println("Disk:", info.Disk)
}
