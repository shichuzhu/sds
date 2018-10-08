package membership

import (
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"time"
)

type NetworkStatsType struct {
	bytesCount int
	startTime  time.Time
	endTime    time.Time
}

var NetworkStats NetworkStatsType
var xmtr *net.UDPConn
var buffer []byte

func (s *NetworkStatsType) InitNetworkStats() {
	s.startTime = time.Now()
	s.bytesCount = 0
}

func (s *NetworkStatsType) GetBandwidthUsage() float32 {
	fmt.Println(s.startTime)
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
	buffer = make([]byte, 4096) // TODO: move this to init
}

func UdpSendSingle(IP string, buf []byte) {
	remote, err := net.ResolveUDPAddr("udp", IP)
	ErrHandler(err)
	_, err = xmtr.WriteToUDP(buf, remote)
	ErrHandler(err)
	NetworkStats.bytesCount += len(buf)
}

// TODO: return pointer type
func AddrStrToBin(addr string) *net.UDPAddr {
	bin, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Panicln("Can't convert address")
	}
	return bin
}

// TODO: add drop packet functionality
func UdpRecvSingle() (*pb.UDPMessage, error) {
	n, err := xmtr.Read(buffer)
	ErrHandler(err)
	UdpMess := pb.UDPMessage{}
	err = proto.Unmarshal(buffer[0:n], &UdpMess)
	ErrHandler(err)
	if mesgType := UdpMess.GetMessageType(); mesgType == "DetectorMessage" {
		fmt.Printf("Received byte %d TYPE %s\n", n, UdpMess.GetDm().GetHeader())
	} else {
		fmt.Printf("Received byte %d TYPE %s\n", n, mesgType)
	}
	return &UdpMess, nil
}
