package main

import (
	"fmt"
	"github.com/paigeman/r1-BDVA/cmd/config"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", config.ServerIP+":"+config.ServerPort)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}(conn)
	message := []byte("Hello, UDP Server!")
	if _, err = conn.Write(message); err != nil {
		fmt.Println("Error sending:", err)
		return
	}
	fmt.Println("Message sent successfully")
}
