package model

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
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

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	// 地球半径（米）
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	// 单位：米
	return R * c
}

func (b *BusStatus) Move(r *rand.Rand) error {
	oldLatitude := b.Latitude
	oldLongitude := b.Longitude
	oldTime := time.Unix(0, b.Timestamp)
	now := time.Now()
	if oldTime.After(now) {
		return fmt.Errorf("%s: Illegal timestamp argument", b.BusID)
	}
	// 纬度或经度变化约对应10米左右
	maxDelta := 0.0001
	deltaLatitude := (r.Float64()*2 - 1) * maxDelta
	deltaLongitude := (r.Float64()*2 - 1) * maxDelta
	newLatitude := oldLatitude + deltaLatitude
	newLongitude := oldLongitude + deltaLongitude
	// 根据haversine公式计算出来的距离
	distance := haversine(oldLatitude, oldLongitude, newLatitude, newLongitude)
	deltaTimestamp := now.Sub(oldTime).Seconds()
	// m/s
	speed := distance / deltaTimestamp
	// 更新数据
	b.Latitude = newLatitude
	b.Longitude = newLongitude
	b.Speed = speed
	b.Timestamp = now.UnixNano()
	return nil
}
