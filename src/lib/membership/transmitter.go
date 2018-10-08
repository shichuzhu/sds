package membership

import (
	pb "fa18cs425mp/src/protobuf"
	"github.com/golang/protobuf/proto"
	"log"
	"math/rand"
	"net"
	"time"
)

type NetworkStatsType struct {
	bytesCount int
	startTime  time.Time
	endTime    time.Time
}

type PacketDropType struct {
	dropRate float32
}

var NetworkStats NetworkStatsType
var xmtr *net.UDPConn
var buffer []byte
var PacketDrop PacketDropType

func (s *PacketDropType) SetDropRate(rate float32) {
	s.dropRate = rate
	rand.Seed(time.Now().Unix())
}

func (s *PacketDropType) rollDiceToDrop() bool {
	if rand.Float32() >= s.dropRate {
		return false
	} else {
		return true
	}
}

func (s *NetworkStatsType) InitNetworkStats() {
	s.startTime = time.Now()
	s.bytesCount = 0
}

func (s *NetworkStatsType) GetBandwidthUsage() float32 {
	s.endTime = time.Now()
	durationInSecond := s.endTime.Sub(s.startTime).Seconds()
	count := s.bytesCount
	s.startTime = time.Now()
	s.bytesCount = 0
	return float32(count) / float32(durationInSecond)
}

func UdpSend(IP string, buf []byte, rep int) {
	//for i := 0; i < rep; i++ {
	for i := 0; i < rep/rep; i++ { // TODO: delete this one, debug use only
		UdpSendSingle(IP, buf)
	}
}

func InitXmtr() {
	if xmtr != nil {
		return
	}
	local := AddrStrToBin(MyAddr)
	var err error
	xmtr, err = net.ListenUDP("udp", local)
	ErrHandler(err)
	buffer = make([]byte, 4096)
}

func UdpSendSingle(IP string, buf []byte) {
	remote, err := net.ResolveUDPAddr("udp", IP)
	ErrHandler(err)
	_, err = xmtr.WriteToUDP(buf, remote)
	ErrHandler(err)
	NetworkStats.bytesCount += len(buf)
}

func AddrStrToBin(addr string) *net.UDPAddr {
	bin, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Panicln("Can't convert address")
	}
	return bin
}

func UdpRecvSingle() (*pb.UDPMessage, error) {
	n, err := xmtr.Read(buffer)
	ErrHandler(err)
	if PacketDrop.rollDiceToDrop() {
		log.Println("Packet dropped.")
		return nil, nil
	}
	UdpMess := pb.UDPMessage{}
	err = proto.Unmarshal(buffer[0:n], &UdpMess)
	ErrHandler(err)
	//if mesgType := UdpMess.GetMessageType(); mesgType == "DetectorMessage" {
	//	log.Printf("Received byte %d TYPE %s\n", n, UdpMess.GetDm().GetHeader())
	//} else {
	//	log.Printf("Received byte %d TYPE %s\n", n, mesgType)
	//}
	return &UdpMess, nil
}
