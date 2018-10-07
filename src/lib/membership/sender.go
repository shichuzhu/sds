package membership

import (
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
				UdpSendSingle(addr, []byte{}) // TODO: add message here.
			}
		}
		time.Sleep(1500 * time.Millisecond)
		for i := range ackWaitEntries {
			if ackWaitEntries[i].acked == true {
				continue
			} else {
				// TODO: Form fail message for ackW[i].addr
				message := []byte{}
				for _, addrs := range MembershipList.getRandomTargets(len(MembershipList.members)) {
					UdpSend(addrs, message, 2)
				}
				MembershipList.deleteID(ackWaitEntries[i].addr, int(^uint(0)>>1))
			}
		}

	}
	return nil
}
