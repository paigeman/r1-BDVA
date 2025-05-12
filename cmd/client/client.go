package main

import (
	"fmt"
	"github.com/paigeman/r1-BDVA/cmd/config"
	"github.com/paigeman/r1-BDVA/cmd/model"
	"math/rand"
	"net"
	"sync"
)

func simulateBus(id string, wg *sync.WaitGroup) {
	defer wg.Done()
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
	// 确保每辆车生成的随机数不一致
	rand.New(rand.NewSource(model.GenerateBusRandomSeed(id)))
}

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
