package membership

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
)

var MyAddr string
var MembershipList MembershipListType
var ackWaitEntries []AckWaitEntry

// Use getter to get list for thread safety issue.
type MemberType struct {
	addr           string
	sessionCounter int
}

type MembershipListType struct {
	members []MemberType
	// Potential global config about MembershipList
	myIP    net.IP
	myIndex int
	MyPort  int
}

func (t *MembershipListType) insert(index int, memberType MemberType) {
	fmt.Println("insert")
	s := &t.members
	*s = append(*s, MemberType{})
	copy((*s)[index+1:], (*s)[index:])
	(*s)[index] = memberType
	t.updateMyIndex()
}

func (t *MembershipListType) delete(index int) {
	fmt.Println("delete")
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

func (s *MembershipListType) lookupID(id string) (MemberType, bool) {
	for _, member := range MembershipList.members {
		if id == member.addr {
			return member, true
		}
	}
	return MemberType{}, false
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
	num = min(num, len(MembershipList.members))
	targets := make([]string, num)
	for i, j := range rand.Perm(len(MembershipList.members))[:num] {
		targets[i] = s.members[j].addr
	}

	return targets
}

func (s *MembershipListType) updateMyIndex() {
	for i, member := range MembershipList.members {
		if member.addr == MyAddr {
			s.myIndex = i
			return
		}
	}
}

func (s *MembershipListType) getPingTargets(num int) []string {
	targets := make([]string, num)
	for i, _ := range targets {
		targets[i] = s.members[(s.myIndex+i)%len(s.members)].addr
	}
	return targets
}

func InitInstance() {
	if ackWaitEntries == nil {
		if MembershipList.MyPort == 0 {
			MembershipList.MyPort = DefaultUdpPort
		}
		ackWaitEntries = make([]AckWaitEntry, NodeNumberToPing)
		MembershipList.myIP = GetOutboundIP()
		MyAddr = MembershipList.myIP.String() + ":" + strconv.Itoa(MembershipList.MyPort)
	}
}

func SetPortNumber(port int) {
	MembershipList.MyPort = port
}

func AddSelfToList(sessionCounter int) {
	MembershipList.insertNewID(MyAddr, sessionCounter)
}

func StartFailureDetector() {
	InitXmtr()
	go receiverService()
	go senderService()
	NetworkStats.InitNetworkStats()
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

func GetListElement() ([]string, []int) {
	ret1 := make([]string, len(MembershipList.members))
	ret2 := make([]int, len(MembershipList.members))
	for i := 0; i < len(MembershipList.members); i++ {
		ret1[i] = MembershipList.members[i].addr
		ret2[i] = MembershipList.members[i].sessionCounter
	}

	return ret1, ret2
}

func DumpTable() {
	table, _ := GetListElement()
	for i, t := range table {
		fmt.Printf("Entry %d has ID: %s\n", i+1, t)
	}
}
