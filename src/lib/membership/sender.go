package membership

import (
	pb "fa18cs425mp/src/protobuf"
	"github.com/golang/protobuf/proto"
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

func senderService() error {
	ackWaitEntries := make([]AckWaitEntry, NodeNumberToPing)
	for {
		memsToPing := MembershipList.getPingTargets(NodeNumberToPing)
		for i, addr := range memsToPing {
			ackWaitEntries[i] = AckWaitEntry{addr: addr}
		}
		for i := 0; i < 3; i++ { // Send 3 times.
			for j, addr := range memsToPing {
				message, _ := proto.Marshal(
					&pb.UDPMessage{MessageType: "DetectorMessage",
						Dm: &pb.DetectorMessage{Header: "Ping", Addr: MyAddr, SessNUm: 0, TTL: 0},
						Fm: &pb.FullMembershipList{}}) // TODO: see if any can be omitted.
				UdpSendSingle(addr, message)
			}
		}
		time.Sleep(1500 * time.Millisecond)
		for i := range ackWaitEntries {
			if ackWaitEntries[i].acked == true {
				continue
			} else {
				addr := ackWaitEntries[i].addr
				if member, exist := MembershipList.lookupID(addr); exist {
					message, _ := proto.Marshal(
						&pb.UDPMessage{MessageType: "DetectorMessage",
							Dm: &pb.DetectorMessage{Header: "Delete", Addr: addr, SessNUm: int32(member.sessionCounter), TTL: 0},
							Fm: &pb.FullMembershipList{}}) // TODO: see if any can be omitted.
					for _, addr := range MembershipList.getRandomTargets(len(MembershipList.members)) {
						UdpSend(addr, message, 2)
					}
					MembershipList.deleteID(ackWaitEntries[i].addr, int(^uint(0)>>1))
				}
			}
		}

	}
	return nil
}
