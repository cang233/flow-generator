package producer

import "testing"

func TestRandPktGenerator_Run(t *testing.T) {
	gen := RandPktGenerator{}
	gen.Run()
}
