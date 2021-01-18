package producer

import (
	"flow-generator/log"
	"testing"
)

func TestRandPktGenerator_Run(t *testing.T) {
	// var nic string
	// flag.StringVar(&nic, "i", "", "nic which sending packets")

	nic := "enp10s0"
	// flag.Parse()
	if len(nic) == 0 {
		log.Panic("Var port name([-i]) is empty.")
	}
	config := make(map[string]string, 0)
	config["i"] = nic

	gen := RandPktGenerator{}
	gen.Init(config)
	gen.Run()
}

func BenchmarkRandPktConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randPktConfig()
	}
}
func BenchmarkRandomChangePacket(b *testing.B) {
	lyrs := randPktConfig()
	for i := 0; i < b.N; i++ {
		randomChangePacket(lyrs)
	}
}
