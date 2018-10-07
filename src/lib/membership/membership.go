package membership

import (
	"log"
	"math/rand"
	"net"
)

// Use getter to get list for thread safety issue.
type MemberType struct {
	addr           string
	sessionCounter int
	acked          bool
}

type MembershipListType struct {
	members []MemberType
	// Potential global config about MembershipList
	myID    string
	myIP    net.IP
	myIndex int
}

func (t *MembershipListType) insert(index int, memberType MemberType) {
	s := &t.members
	*s = append(*s, MemberType{})
	copy((*s)[index+1:], (*s)[index:])
	(*s)[index] = memberType
	t.updateMyIndex()
}

func (t *MembershipListType) delete(index int) {
	s := &t.members
	copy((*s)[index:], (*s)[index+1:])
	*s = (*s)[:len(*s)-1]
	t.updateMyIndex()
}

func (s *MembershipListType) insertNewID(id string, sessionID int) {
	for i, member := range MembershipList.members {
		if id == member.addr {
			if sessionID > member.sessionCounter {
				s.members[i].sessionCounter = sessionID
			}
			return
		} else if id < member.addr {
			s.insert(i, MemberType{addr: id, sessionCounter: sessionID})
			return
		}
	}
	s.insert(len(MembershipList.members), MemberType{addr: id, sessionCounter: sessionID})
}

func (s *MembershipListType) deleteID(id string, sessionID int) {
	for i, _ := range MembershipList.members {
		member := &MembershipList.members[i]
		if member.addr == id {
			if member.sessionCounter <= sessionID {
				s.delete(i)
			}
			return
		} else if member.addr > id {
			return
		}
	}
}

func (s *MembershipListType) getRandomTargets(num int) []string {
	targets := make([]string, num)
	for i, j := range rand.Perm(len(MembershipList.members))[:num] {
		targets[i] = s.members[j].addr
	}
}

func (s *MembershipListType) updateMyIndex() {
	for i, member := range MembershipList.members {
		if member.addr == s.myID {
			s.myIndex = i
			return
		}
	}
}

func (s *MembershipListType) getPingTargets(num int) []string {
	targets := make([]string, NodeNumberToPing)
	for i, _ := range targets {
		targets[i] = s.members[(s.myIndex+i)%len(s.members)].addr
	}
	return targets
}

func StartFailureDetector() {
	MembershipList.myIP = GetOutboundIP()
	go receiverService()
	go senderService()
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

var MembershipList MembershipListType
