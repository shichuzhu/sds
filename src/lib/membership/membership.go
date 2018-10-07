package membership

import (
	"log"
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
}

func (s *MembershipListType) insertNewID(id string, sessionID int) {
	for i, member := range MembershipList.members {
		if id == member.addr {
			if sessionID > member.sessionCounter {
				member.sessionCounter = sessionID
			}
			return
		} else if id > member.addr {
			continue
		} else {
			s.insert(i, MemberType{addr: id, sessionCounter: sessionID})
			return
		}
	}

}
func (s *MembershipListType) deleteID(id string, sessionID int) {

}
func (s *MembershipListType) getRandomTargets(num int) []string {
	return nil
}
func (s *MembershipListType) updateMyIndex() {
	return
}

func (s *MembershipListType) getPingTargets(num int) []string {
	return nil
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
