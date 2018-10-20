package membership

import (
	pb "fa18cs425mp/src/protobuf"
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
	InitInstance()
	InitXmtr()
	message, err := proto.Marshal(
		&pb.UDPMessage{MessageType: "DetectorMessage",
			Dm: &pb.DetectorMessage{Header: "NewJoin", Addr: MyAddr, SessNum: 0, Ttl: 0},
			Fm: nil})
	ErrHandler(err)

	UdpSend(introAddr, message, 3)
	StartFailureDetector()
}

func ReportFailure(addr string) {
	log.Println("Failure DETECTED: ", addr)
	if member, exist := MembershipList.lookupID(addr); exist {
		message, _ := proto.Marshal(
			&pb.UDPMessage{MessageType: "DetectorMessage",
				Dm: &pb.DetectorMessage{Header: "Delete", Addr: addr, SessNum: int32(member.sessionCounter), Ttl: 0},
				Fm: nil})
		for _, addr := range MembershipList.getRandomTargets(len(MembershipList.members)) {
			UdpSend(addr, message, 2)
		}
		MembershipList.deleteID(addr, int(^uint(0)>>1))
	}
	if addr == MyAddr {
		log.Fatalln("False positive detected, auto terminating.")
	}
}

func senderService() error {
	for {
		memsToPing := MembershipList.getPingTargets(NodeNumberToPing)
		//log.Printf("memsToPing %s\n", memsToPing)
		for i, addr := range memsToPing {
			ackWaitEntries[i] = AckWaitEntry{addr: addr}
		}
		for i := 0; i < MultiSendNumber; i++ { // Send 3 times.
			for _, addr := range memsToPing {
				message, _ := proto.Marshal(
					&pb.UDPMessage{MessageType: "DetectorMessage",
						Dm: &pb.DetectorMessage{Header: "Ping", Addr: MyAddr, SessNum: 0, Ttl: 0},
						Fm: nil})
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
	return nil
}
