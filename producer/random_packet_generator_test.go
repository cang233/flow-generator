package producer

import (
	"flag"
	"flow-generator/log"
	"testing"
)

func TestRandPktGenerator_Run(t *testing.T) {
	var nic string
	flag.StringVar(&nic, "i", "", "nic which sending packets")

	flag.Parse()
	if len(nic) == 0 {
		log.Panic("Var port name([-i]) is empty.")
	}
	config := make(map[string]string, 0)
	config["i"] = nic

	gen := RandPktGenerator{}
	gen.Init(config)
	gen.Run()
}
