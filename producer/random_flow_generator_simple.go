package producer

import (
	"flag"
	"flow-generator/util"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type SimpleFlowGenerator struct {
	config  map[string]string
	handler *pcap.Handle
}

func (s *SimpleFlowGenerator) getOptions() {
	var nicName string
	flag.StringVar(&nicName, "i", "enp10s0", "the nic name used for sending packets.")

	flag.Parse()
	s.config["i"] = nicName
}

func (s *SimpleFlowGenerator) Init() {
	//init struct
	s.config = make(map[string]string, 0)
	//get options
	s.getOptions()

	//init other
	s.initHandler()
}

func (s *SimpleFlowGenerator) initHandler() {
	var (
		snapshotLen int32 = 65535
		promiscuous bool  = false
		err         error
		timeout     time.Duration = 30 * time.Second
		nicName     string
	)
	if v, ok := s.config["i"]; ok {
		nicName = v
	} else {
		log.Panic("nic name needed!")
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range devices {
		//这里为了兼容win与linux，同时检查2个
		if value.Description == nicName || value.Name == nicName {
			//Open device
			s.handler, err = pcap.OpenLive(value.Name, snapshotLen, promiscuous, timeout)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println(value.Description, value.Name)
	}

	if s.handler == nil {
		log.Panicln("can not find nic:", nicName)
	}
}

func (s *SimpleFlowGenerator) Run() {
	defer s.handler.Close()

	var (
		buffer      gopacket.SerializeBuffer
		options     gopacket.SerializeOptions
		err         error
		flowPackets int
		lys         *randomPacketLayer
	)

	// And create the packet with the layers
	buffer = gopacket.NewSerializeBuffer()

	count := 0
	lastTime := time.Now()
	//loop
	for {
		flowPackets = util.RandomInt(6, 32)
		lys = randPktConfig()
		//gen a flow
		for flowPackets > 0 {
			//change content
			gopacket.SerializeLayers(buffer, options,
				lys.ether,
				lys.ipv4,
				lys.tcp,
				gopacket.Payload(util.RandomBytes(20)),
			)
			err = s.handler.WritePacketData(buffer.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			//change flag
			if flowPackets == 1 {
				//the last of a flow
				lys.tcp.SYN = false
				lys.tcp.PSH = false
				lys.tcp.FIN = util.RandomBoolean() //have ending or not
			} else {
				lys.tcp.SYN = false
				lys.tcp.FIN = false
				lys.tcp.PSH = true
			}
			//count
			count++
			now := time.Now()
			if now.Sub(lastTime) > time.Second*1 {
				fmt.Println("pps:", count)
				lastTime = now
				count = 0
			}
			//
			flowPackets--

		}

	}
}

func simpleRandomChangePacket(lys *randomPacketLayer) {
	lys.ipv4.SrcIP = util.RandomIPv4()
	lys.ipv4.DstIP = util.RandomIPv4()
	lys.tcp.SrcPort = layers.TCPPort(util.RandomPort())
	lys.tcp.DstPort = layers.TCPPort(util.RandomPort())
}
