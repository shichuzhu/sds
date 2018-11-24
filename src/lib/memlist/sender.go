package memlist

import (
	"fa18cs425mp/src/pb"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
	"time"
)

type AckWaitEntry struct {
	addr  string
	acked bool
	lock  sync.Mutex
}

func (s *AckWaitEntry) ack() {
	s.lock.Lock()
	s.acked = true
	s.lock.Unlock()
}

func (s *AckWaitEntry) reset() {
	s.lock.Lock()
	s.acked = false
	s.lock.Unlock()
}

func ContactIntroducer(introAddr string) {
	message, err := proto.Marshal(
		&pb.UDPMessage{MessageType: "DetectorMessage",
			Dm: &pb.DetectorMessage{
				Header: "NewJoin", Ttl: 0, Mem: &pb.Member{
					Addr:     MyAddr,
					GrpcAddr: MemList.members[MemList.myIndex].grpcAddr,
					SessNum:  0,
					NodeId:   int32(MemList.MyNodeId)}}})
	ErrHandler(err)

	UdpSend(introAddr, message, 1)
	StartFailureDetector()
}

func ReportFailure(addr string) {
	log.Println("Failure DETECTED: ", addr)
	if member, exist := MemList.lookupID(addr); exist {
		message, _ := proto.Marshal(
			&pb.UDPMessage{MessageType: "DetectorMessage",
				Dm: &pb.DetectorMessage{Header: "Delete", Mem: mem2pb(&member), Ttl: 0}})
		for _, addr := range MemList.getRandomTargets(len(MemList.members)) {
			UdpSend(addr, message, 1)
		}
		MemList.deleteID(addr, int(^uint(0)>>1))
	}
	if addr == MyAddr {
		log.Fatalln("False positive detected, auto terminating.")
	}
}

func senderService() {
	for {
		memsToPing := MemList.getPingTargets(NodeNumberToPing)
		//log.Printf("memsToPing %s\n", memsToPing)
		for i, addr := range memsToPing {
			ackWaitEntries[i] = AckWaitEntry{addr: addr}
		}
		for i := 0; i < MultiSendNumber; i++ { // Send multiple times.
			for _, addr := range memsToPing {
				message, _ := proto.Marshal(
					&pb.UDPMessage{MessageType: "DetectorMessage",
						Dm: &pb.DetectorMessage{Header: "Ping", Mem: &pb.Member{Addr: MyAddr}}})
				UdpSendSingle(addr, message)
			}
		}
		time.Sleep(time.Duration(FailureTimeout) * time.Millisecond)
		for i := range memsToPing {
			if ackWaitEntries[i].acked == true {
				continue
			} else {
				ReportFailure(ackWaitEntries[i].addr)
			}
		}
	}
}
