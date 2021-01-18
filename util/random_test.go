package util

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
)

func TestRandomMacN(t *testing.T) {
	N := 10000
	s := net.HardwareAddr{0xFB, 0xBA, 0xFA, 0xAA, 0xF6, 0xAA}
	e := net.HardwareAddr{0x2F, 0xFF, 0x4F, 0x00, 0x3F, 0xF0}
	for i := 0; i < N; i++ {
		s, e = RandomMacN(s, e), RandomMacN(s, e)
		fmt.Println(s.String(), e.String())
	}
}

func BenchmarkRandomMacN(b *testing.B) {
	s := net.HardwareAddr{0xFB, 0xBA, 0xFA, 0xAA, 0xF6, 0xAA}
	e := net.HardwareAddr{0x2F, 0xFF, 0x4F, 0x00, 0x3F, 0xF0}
	for i := 0; i < b.N; i++ {
		RandomMacN(s, e)
	}
}

func TestRandomPortN(t *testing.T) {
	N := 1000
	var r int
	for i := 0; i < N; i++ {
		r = RandomPortN(10, 65535)
		fmt.Println(r)
	}
}

func BenchmarkRandomPortN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomPortN(10, 65535)
	}
}

func TestRandomIPv4N(t *testing.T) {
	N := 10000
	s := net.IP{12, 13, 11, 10}
	e := net.IP{44, 33, 55, 88}
	for i := 0; i < N; i++ {
		a, b := RandomIPv4N(s, e), RandomIPv4N(s, e)
		fmt.Println(a.String(), b.String())
	}
	for i := 0; i < N; i++ {
		s, e = RandomIPv4N(s, e), RandomIPv4N(s, e)
		fmt.Println(s.String(), e.String())
	}
}

func BenchmarkRandomIPv4N(b *testing.B) {
	s := net.IP{12, 13, 11, 10}
	e := net.IP{144, 233, 155, 188}
	for i := 0; i < b.N; i++ {
		RandomIPv4N(s, e)
	}
}

func TestRandomBoolean(t *testing.T) {
	N := 1000
	for i := 0; i < N; i++ {
		fmt.Println(RandomBoolean())
	}
}

func BenchmarkRandomBoolean(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomBoolean()
	}
}

func TestRandomString(t *testing.T) {
	N := 1000
	for i := 0; i < N; i++ {
		fmt.Println(RandomString(rand.Intn(20)))
	}
}

func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomString(16)
	}
}
