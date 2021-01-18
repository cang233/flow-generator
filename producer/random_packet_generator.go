package producer

import (
	"flow-generator/util"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type WorkStatus int

const (
	WorkStatusStop = iota
	WorkStatusRunning
)

var (
	wg   sync.WaitGroup
	lock sync.Mutex
)

type RandPktGenerator struct {
	workers       []*worker
	sender        *sender
	runningStatus WorkStatus
	defaultPacket *randomPacketLayer
	speed         uint64
	config        map[string]string
}

type randomPacketLayer struct {
	ether *layers.Ethernet
	ipv4  *layers.IPv4
	tcp   *layers.TCP
}

func (r *RandPktGenerator) Init(config map[string]string) {
	r.config = config
	r.init()
}

func (r *RandPktGenerator) Run() {
	r.sender.Run()
	for i := range r.workers {
		r.workers[i].Run()
	}

	wg.Wait()
}

func defaultConfig() map[string]int {
	return map[string]int{
		"default_worker_count": 10,
	}
}

func defaultPacketConfig() *randomPacketLayer {
	return &randomPacketLayer{
		ether: &layers.Ethernet{
			EthernetType: layers.EthernetTypeIPv4,
			SrcMAC:       net.HardwareAddr{0xFB, 0xBA, 0xFA, 0xAA, 0xF6, 0xAA},
			DstMAC:       net.HardwareAddr{0x2F, 0xFF, 0x4F, 0xF6, 0x3F, 0xF0},
		},
		ipv4: &layers.IPv4{
			Protocol: layers.IPProtocolTCP,
			Flags:    0x0000,
			IHL:      0x45, //version + header length
			TTL:      0x80,
			Id:       0x1234,
			Length:   0x014e,
			SrcIP:    net.IP{12, 13, 11, 10},
			DstIP:    net.IP{44, 33, 55, 88},
		},
		tcp: &layers.TCP{
			Seq:        0x1234,
			Ack:        0x1235,
			SrcPort:    layers.TCPPort(68),
			DstPort:    layers.TCPPort(67),
			SYN:        true,
			DataOffset: 0x5,
		},
	}
}

//初始化流信息和worker信息
func (r *RandPktGenerator) init() {
	config := defaultConfig()
	var (
		workerCount int
	)
	//set worker count
	if v, ok := config["default_worker_count"]; ok {
		workerCount = v
	} else {
		workerCount = 10
	}
	//init default packet config
	r.defaultPacket = defaultPacketConfig()
	//init sender
	r.sender = new(sender)
	r.sender.Init(r.config)
	r.sender.father = r //set father
	//init workers
	r.workers = nil
	for i := 0; i < workerCount; i++ {
		w := new(worker)
		w.Init(r.sender.dsChan, len(r.workers))
		r.workers = append(r.workers, w)
	}
}

func (r *RandPktGenerator) addWorker() {
	lock.Lock()
	defer lock.Unlock()
	w := new(worker)
	w.Init(r.sender.dsChan, len(r.workers))
	r.workers = append(r.workers, w)
	w.Run() //run immediately
}

func (r *RandPktGenerator) removeWorker() {
	lock.Lock()
	defer lock.Unlock()
	r.workers[len(r.workers)-1].Stop()
	r.workers = r.workers[:len(r.workers)]
}

//
func (r *RandPktGenerator) Stop() {
	for i := range r.workers {
		r.workers[i].Stop()
	}
	r.sender.Stop()
}

func randPktConfig() *randomPacketLayer {
	return &randomPacketLayer{
		ether: &layers.Ethernet{
			EthernetType: layers.EthernetTypeIPv4,
			SrcMAC:       util.RandomMac(),
			DstMAC:       util.RandomMac(),
		},
		ipv4: &layers.IPv4{
			Protocol: layers.IPProtocolTCP,
			Flags:    0x0000,
			IHL:      0x45, //version + header length
			TTL:      0x80,
			Id:       0x1234,
			//Length:   0x014e,
			SrcIP: util.RandomIPv4(),
			DstIP: util.RandomIPv4(),
		},
		tcp: &layers.TCP{
			Seq:        util.RandomSequence(),
			Ack:        util.RandomSequence(),
			SrcPort:    layers.TCPPort(util.RandomPort()),
			DstPort:    layers.TCPPort(util.RandomPort()),
			SYN:        true,
			DataOffset: 0x5,
		},
	}
}

type worker struct {
	id            int
	runningStatus WorkStatus
	exportChan    chan []byte
	sts           statics
	tmplate       *randomPacketLayer
}

func (w *worker) Init(exportChan chan []byte, id int) {
	w.runningStatus = WorkStatusStop
	w.exportChan = exportChan
	w.id = id
	w.tmplate = randPktConfig()
}

func randomChangePacket(lys *randomPacketLayer) {
	lys.ether.SrcMAC = util.RandomMac()
	lys.ether.DstMAC = util.RandomMac()
	lys.ipv4.SrcIP = util.RandomIPv4()
	lys.ipv4.DstIP = util.RandomIPv4()
	lys.tcp.Seq = util.RandomSequence()
	lys.tcp.Ack = util.RandomSequence()
	lys.tcp.SrcPort = layers.TCPPort(util.RandomPort())
	lys.tcp.DstPort = layers.TCPPort(util.RandomPort())
}

func (w *worker) Run() {
	w.runningStatus = WorkStatusRunning
	wg.Add(1)
	go func() {
		defer wg.Done()
		stopTicker := time.NewTicker(time.Second * 5)
		logTicker := time.NewTicker(time.Second * 1)
		defer stopTicker.Stop()
		defer logTicker.Stop()
		for w.runningStatus == WorkStatusRunning {
			select {
			case <-stopTicker.C:
			case <-logTicker.C:
				log.Printf("[worker-%d]pps=%d", w.id, w.sts.CountLoad())
				w.sts.CountClear()
			default: //do rand update
				randomChangePacket(w.tmplate)
				buffer := gopacket.NewSerializeBuffer()
				gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{},
					w.tmplate.ether,
					w.tmplate.ipv4,
					w.tmplate.tcp,
					gopacket.Payload(util.RandomBytes(util.RandomInt(0, 1000))),
				)
				w.exportChan <- buffer.Bytes()
				w.sts.CountAdd()
			}
		}
		w.Clear()
	}()
}
func (w *worker) Stop() {
	w.runningStatus = WorkStatusStop
}

func (w *worker) Clear() {
}

type sender struct {
	dsChan       chan []byte
	running      WorkStatus
	handler      *pcap.Handle
	sts          statics
	father       *RandPktGenerator
	reachTopTime time.Time
}

func (s *sender) Init(config map[string]string) {
	s.dsChan = make(chan []byte, 100*1000)
	s.running = WorkStatusStop
	//get handler
	var (
		snapshot_len int32 = 65535
		promiscuous  bool  = false
		err          error
		timeout      time.Duration = 30 * time.Second
		nicName      string        = "Realtek Gaming GbE Family Controller"
	)
	//init config
	if v, ok := config["i"]; ok {
		nicName = v
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range devices {
		if value.Name == nicName || value.Description == nicName { //为了兼容windows和linux
			//Open device
			s.handler, err = pcap.OpenLive(value.Name, snapshot_len, promiscuous, timeout)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println(value.Name, ":", value.Description)
	}
	if s.handler == nil {
		log.Panic("Init handle in sender is nil,can not find the nic named " + nicName)
	}
}

func (s *sender) Run() {
	s.running = WorkStatusRunning
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s.handler.Close()

		stopTicker := time.NewTicker(time.Second * 5)
		logTicker := time.NewTicker(time.Second * 1)
		upperTicker := time.NewTicker(time.Second * 5)
		defer stopTicker.Stop()
		for s.running == WorkStatusRunning {
			select {
			case <-stopTicker.C: //do nothing
			case <-logTicker.C:
				log.Printf("[sender]chan.cap()=%d,chan.len()=%d,pps=%d", cap(s.dsChan), len(s.dsChan), s.sts.CountLoad())
				s.sts.CountClear()
			case <-upperTicker.C:
				now := time.Now()
				if now.Sub(s.reachTopTime) < time.Second*60 {
					continue
				}
				if cap(s.dsChan) > 3*len(s.dsChan) {
					s.father.addWorker()
					log.Println("[sender]chan is not full,call father adding a worker")
				} else {
					s.father.removeWorker()
					s.reachTopTime = now
					log.Println("[sender]chan is full,call father removing a worker")
				}

			case v, ok := <-s.dsChan:
				if !ok {
					continue
				}
				err := s.handler.WritePacketData(v)
				s.sts.CountAdd()
				if err != nil {
					log.Panic(err)
				}
			}
		}
		s.running = WorkStatusStop

	}()
}
func (s *sender) Stop() {
	s.running = WorkStatusStop
}

type statics struct {
	countPerS uint64
}

func (s *statics) CountAdd() {
	atomic.AddUint64(&s.countPerS, 1)
}
func (s *statics) CountLoad() uint64 {
	return atomic.LoadUint64(&s.countPerS)
}
func (s *statics) CountClear() {
	atomic.StoreUint64(&s.countPerS, 0)
}
