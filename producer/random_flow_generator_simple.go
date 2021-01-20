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
	config  *simpleConfig
	handler *pcap.Handle
}

type simpleConfig struct {
	nicName       string
	flowPktMax    int
	pktPayloadMax int
}

func (s *SimpleFlowGenerator) getOptions() {
	s.config = &simpleConfig{}
	flag.StringVar(&s.config.nicName, "i", "enp10s0", "the nic name used for sending packets.")
	flag.IntVar(&s.config.flowPktMax, "pmn", 32, "max number of a flow's packet,not less than 6")
	flag.IntVar(&s.config.pktPayloadMax, "ps", 10, "max packet payload size")

	flag.Parse()
}

func (s *SimpleFlowGenerator) Init() {
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
	)

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range devices {
		//这里为了兼容win与linux，同时检查2个
		if value.Description == s.config.nicName || value.Name == s.config.nicName {
			//Open device
			s.handler, err = pcap.OpenLive(value.Name, snapshotLen, promiscuous, timeout)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println(value.Description, value.Name)
	}

	if s.handler == nil {
		log.Panicln("can not find nic:", s.config.nicName)
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
		flowPackets = util.RandomInt(6, s.config.flowPktMax)
		lys = randPktConfig()
		//gen a flow
		for flowPackets > 0 {
			//change content
			gopacket.SerializeLayers(buffer, options,
				lys.ether,
				lys.ipv4,
				lys.tcp,
				gopacket.Payload(util.RandomBytes(s.config.pktPayloadMax)),
			)
			err = s.handler.WritePacketData(buffer.Bytes())
			if err != nil {
				for i := 0; i < 10; i++ {
					err = s.handler.WritePacketData(buffer.Bytes())
					if err != nil {
						log.Panicln(err)
					}else{
						break
					}
					time.Sleep(time.Millisecond * 10)
				}
				if err != nil {
					log.Fatal(err)
				}
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
