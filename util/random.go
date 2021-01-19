package util

import (
	"log"
	"math/rand"
	"net"
)

func RandomInt(start, end int) int {
	if end <= start {
		log.Panic("RandomInt end is less than start")
	}
	return rand.Intn(end-start) + start
}

func RandomMac() net.HardwareAddr {
	return RandomMacN(net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, net.HardwareAddr{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
}

func RandomPort() int {
	return RandomPortN(0, 65535)
}

func RandomIPv4() net.IP {
	return RandomIPv4N(net.IP{0, 0, 0, 0}, net.IP{255, 255, 255, 255})
}

func RandomSequence() uint32 {
	return rand.Uint32()
}

func RandomBoolean() bool {
	const N = 1 << 10
	return rand.Intn(N) >= N/2
}

func RandomString(n int) string {
	return string(RandomBytes(n))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomBytes(length int) []byte {
	if length < 0 {
		log.Panic("RandomString length is less than 0")
	}
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

//RandomBoolean1 return random bool lists which only has 1 true element.
func RandomBoolean1(n int) (blist []bool) {
	for i := 0; i < n; i++ {
		blist = append(blist, false)
	}
	blist[rand.Intn(n)] = true
	return
}

//RandomMacN return random mac addr between start and end,
// notice that only support every position start[i]<end[i]
func RandomMacN(start, end net.HardwareAddr) (rd net.HardwareAddr) {
	for i := 0; i < 6; i++ {
		if end[i] <= start[i] || end[i]-start[i] >= 0xFF {
			rd = append(rd, byte(rand.Intn(0xFF)))
		} else {
			rd = append(rd, byte(rand.Intn(int(end[i]-start[i]))+int(start[i])))
		}
	}
	return rd
}

func RandomPortN(start, end int) int {
	if end <= start {
		log.Panic("RandomPortN end is less than start")
	}
	if end-start >= 65535 {
		return rand.Intn(65535)
	}
	return rand.Intn(end-start) + start
}

func RandomIPv4N(start, end net.IP) (rd net.IP) {
	for i := 0; i < 4; i++ {
		if end[i] <= start[i] || end[i]-start[i] >= 255 {
			rd = append(rd, byte(rand.Intn(255)))
			continue
		}
		rd = append(rd, byte(rand.Intn(int(end[i]-start[i]))+int(start[i])))
	}
	return rd
}
