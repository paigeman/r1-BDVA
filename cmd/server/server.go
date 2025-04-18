package main

import (
	"fmt"
	"net"
)
import "github.com/paigeman/r1-BDVA/cmd/config"

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":"+config.ServerPort)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}(conn)
	fmt.Printf("UDP server started on port %s\n", config.ServerPort)
	buffer := make([]byte, config.BufferSize)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}
		fmt.Printf("Received from %s: %s\n", clientAddr, string(buffer[:n]))
	}
}
