package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
)
import "github.com/paigeman/r1-BDVA/cmd/config"

func main() {
	kafkaConn, err := kafka.DialLeader(context.Background(), "tcp", config.KafkaAddress, config.KafkaTopic, config.KafkaPartition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer func(kafkaConn *kafka.Conn) {
		err := kafkaConn.Close()
		if err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}(kafkaConn)
	addr, err := net.ResolveUDPAddr("udp", ":"+config.ServerPort)
	if err != nil {
		log.Println("Error resolving address:", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
		}
	}(conn)
	log.Printf("UDP server started on port %s\n", config.ServerPort)
	buffer := make([]byte, config.BufferSize)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error reading:", err)
			continue
		}
		log.Printf("Received from %s: %s\n", clientAddr, string(buffer[:n]))
		_, err = kafkaConn.WriteMessages(
			kafka.Message{Value: buffer[:n]},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}
}
