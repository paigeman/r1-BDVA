package main

import (
	"encoding/json"
	"github.com/paigeman/r1-BDVA/cmd/config"
	"github.com/paigeman/r1-BDVA/cmd/model"
	"log"
	"math/rand"
	"net"
	"sync"
)

func simulateBus(id string, wg *sync.WaitGroup) {
	defer wg.Done()
	addr, err := net.ResolveUDPAddr("udp", config.ServerIP+":"+config.ServerPort)
	if err != nil {
		log.Println(id, ": Error resolving address:", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Println(id, ": Error connecting:", err)
		return
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			log.Println(id, ": Error closing connection:", err)
		}
	}(conn)
	// 确保每辆车生成的随机数不一致
	r := rand.New(rand.NewSource(model.GenerateBusRandomSeed(id)))
	// 初始经纬度可以调整
	latitude := 23.12
	longitude := 113.28
	busStatus := model.BusStatus{
		BusID:     id,
		Latitude:  latitude,
		Longitude: longitude,
		Speed:     0,
	}
	busStatus.Now()
	for {
		data, err := json.Marshal(&busStatus)
		if err != nil {
			log.Println(id, ": Error marshalling busStatus:", err)
		}
		_, err = conn.Write(data)
		if err != nil {
			log.Println(id, ": Error writing busStatus:", err)
		}
		//time.Sleep()
		// 更新状态
		err = busStatus.Move(r)
		if err != nil {
			log.Println(id, ": Error moving busStatus:", err)
		}
	}
}

func main() {
	numOfBuses := 10
	idSet := model.GenerateBusIDs(numOfBuses)
	var wg sync.WaitGroup
	for _, id := range idSet {
		wg.Add(1)
		go simulateBus(id, &wg)
	}
	wg.Wait()
}
