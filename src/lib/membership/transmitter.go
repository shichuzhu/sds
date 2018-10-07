package membership

import (
	"net"
	"time"
)

type NetworkStatsType struct {
	bytesCount int
	startTime  time.Time
	endTime    time.Time
}

var networkStats NetworkStatsType

func (s *NetworkStatsType) initNetworkStats() {
	s.startTime = time.Now()
	s.bytesCount = 0
}

func (s *NetworkStatsType) getBandwidthUsage() float32 {
	s.endTime = time.Now()
	durationInSecond := s.endTime.Sub(s.startTime).Seconds()
	return float32(s.bytesCount) / float32(durationInSecond)
}

func UdpSend(IP string, buf []byte, rep int) {
	addr, _ := net.ResolveUDPAddr("udp", IP)
	conn, _ := net.DialUDP("udp", nil, addr)
	for i := 0; i < rep; i++ {
		_, err := conn.WriteToUDP(buf, addr)
		ErrHandler(err)
		networkStats.bytesCount += len(buf)
	}
	conn.Close()
}

func UdpSendSingle(IP string, buf []byte) {
	addr, _ := net.ResolveUDPAddr("udp", IP)
	conn, _ := net.DialUDP("udp", nil, addr)
	_, err := conn.WriteToUDP(buf, addr)
	ErrHandler(err)
	networkStats.bytesCount += len(buf)
	conn.Close()
}

// TODO: implement this when the message is done.
//func UdpRecvSingle() pb.Message, error  {
//	return nil, nil
//}
