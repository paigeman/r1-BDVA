package model

import (
	"fmt"
	"hash/fnv"
	"time"
)

type BusStatus struct {
	BusID     string  `json:"bus_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
	Timestamp int64   `json:"timestamp"`
}

func (b *BusStatus) Now() {
	b.Timestamp = time.Now().UnixNano()
}

func GenerateBusIDs(n int) []string {
	ids := make([]string, n)
	width := len(fmt.Sprintf("%d", n))      // 自动计算需要几位数
	format := fmt.Sprintf("B%%0%dd", width) // 如 "B%03d" 或 "B%05d"
	for i := 0; i < n; i++ {
		ids[i] = fmt.Sprintf(format, i+1)
	}
	return ids
}

func GenerateBusRandomSeed(id string) int64 {
	h := fnv.New64()
	_, err := h.Write([]byte(id))
	if err != nil {
		return time.Now().UnixNano()
	}
	// 做哈希扰动
	return time.Now().UnixNano() ^ int64(h.Sum64())
}
