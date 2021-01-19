package test

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func TestUseGopacket(t *testing.T) {
	var (
		snapshot_len int32 = 65535
		promiscuous  bool  = false
		err          error
		timeout      time.Duration = 30 * time.Second
		handle       *pcap.Handle
		buffer       gopacket.SerializeBuffer
		options      gopacket.SerializeOptions
	)

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range devices {
		if value.Description == "Realtek Gaming GbE Family Controller" || value.Name == "enp10s0" {
			//Open device
			handle, err = pcap.OpenLive(value.Name, snapshot_len, promiscuous, timeout)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println(value.Description, value.Name)
	}
	defer handle.Close()
	// Send raw bytes over wire
	data := "this is raw data in packets"
	data = ""
	var bf bytes.Buffer
	bf.WriteString(data)
	rawBytes := bf.String()

	// This time lets fill out some information
	ipLayer := &layers.IPv4{
		Protocol: layers.IPProtocolTCP,
		Flags:    0x0000,
		IHL:      0x45, //version + header length
		TTL:      0x80,
		Id:       0x1234,
		Length:   0x014e,
		SrcIP:    net.IP{12, 13, 11, 10},
		DstIP:    net.IP{44, 33, 55, 88},
	}
	ethernetLayer := &layers.Ethernet{
		EthernetType: layers.EthernetTypeIPv4,
		SrcMAC:       net.HardwareAddr{0xFB, 0xBA, 0xFA, 0xAA, 0xF6, 0xAA},
		DstMAC:       net.HardwareAddr{0x2F, 0xFF, 0x4F, 0xF6, 0x3F, 0xF0},
	}
	//udpLayer := &layers.UDP{
	//	SrcPort: layers.UDPPort(68),
	//	DstPort: layers.UDPPort(67),
	//	Length:  0x013a,
	//}
	tcpLayer := &layers.TCP{
		Seq:        0x1234,
		Ack:        0x1235,
		SrcPort:    layers.TCPPort(68),
		DstPort:    layers.TCPPort(67),
		SYN:        true,
		DataOffset: 0x5,
	}
	// And create the packet with the layers
	buffer = gopacket.NewSerializeBuffer()

	count := 0
	lastTime := time.Now()
	for {
		gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer,
			tcpLayer,
			gopacket.Payload(rawBytes),
		)
		outgoingPacket := buffer.Bytes()

		err = handle.WritePacketData(outgoingPacket)
		if err != nil {
			log.Fatal(err)
		}
		count++
		now := time.Now()
		if now.Sub(lastTime) > time.Second*1 {
			fmt.Println("pps:", count)
			lastTime = now
			count = 0
		}
	}

}

func TestCountOx(t *testing.T) {
	fmt.Println(0x45, 0x014e)
}

func TestRand(t *testing.T) {
	fmt.Println(rand.Int())
	fmt.Println(rand.Int31())

}
