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
var ThreadsOn bool

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

func (ml *MembershipListType) insert(index int, memberType MemberType) {
	s := &ml.members
	*s = append(*s, MemberType{})
	copy((*s)[index+1:], (*s)[index:])
	(*s)[index] = memberType
	ml.updateMyIndex()
}

func (ml *MembershipListType) delete(index int) {
	s := &ml.members
	copy((*s)[index:], (*s)[index+1:])
	*s = (*s)[:len(*s)-1]
	ml.updateMyIndex()
}

func (ml *MembershipListType) insertNewID(id string, sessionID int) {
	for i, member := range MembershipList.members {
		if id == member.addr {
			if sessionID > member.sessionCounter {
				ml.members[i].sessionCounter = sessionID
			}
			return
		} else if id < member.addr {
			ml.insert(i, MemberType{addr: id, sessionCounter: sessionID})
			return
		}
	}
	ml.insert(len(MembershipList.members), MemberType{addr: id, sessionCounter: sessionID})
}

func (ml *MembershipListType) lookupID(id string) (MemberType, bool) {
	for _, member := range MembershipList.members {
		if id == member.addr {
			return member, true
		}
	}
	return MemberType{}, false
}
func (ml *MembershipListType) deleteID(id string, sessionID int) {
	for i := range MembershipList.members {
		member := &MembershipList.members[i]
		if member.addr == id {
			if member.sessionCounter <= sessionID {
				ml.delete(i)
			}
			return
		} else if member.addr > id {
			return
		}
	}
}

func (ml *MembershipListType) getRandomTargets(num int) []string {
	num = min(num, len(MembershipList.members))
	targets := make([]string, num)
	for i, j := range rand.Perm(len(MembershipList.members))[:num] {
		targets[i] = ml.members[j].addr
	}

	return targets
}

func (ml *MembershipListType) updateMyIndex() {
	for i, member := range MembershipList.members {
		if member.addr == MyAddr {
			ml.myIndex = i
			return
		}
	}
}

func (ml *MembershipListType) getPingTargets(num int) []string {
	num = max(0, min(NodeNumberToPing, len(MembershipList.members)-1))
	targets := make([]string, num)
	for i := range targets {
		targets[i] = ml.members[(ml.myIndex+i+1)%len(ml.members)].addr
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
		AddSelfToList(0)
	}
}

func AddSelfToList(sessionCounter int) {
	MembershipList.insertNewID(MyAddr, sessionCounter)
}

func StartFailureDetector() {
	InitXmtr()
	if !ThreadsOn {
		go receiverService()
		go senderService()
		ThreadsOn = true
	}
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
		fmt.Printf("Index %d is Process: %s\n", i, t)
	}
}

func FormDumpTable() string {
	var response string
	table, _ := GetListElement()
	for i, t := range table {
		response += fmt.Sprintf("Index %d is Process: %s\n", i, t)
	}
	return response
}
