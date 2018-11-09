package membership

import (
	"fa18cs425mp/src/lib/sdfs/sdfs2fd"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var MyAddr string
var MembershipList MembershipListType
var ackWaitEntries []AckWaitEntry
var ThreadsOn bool

// Use getter to get list for thread safety issue.
type MemberType struct {
	addr           string
	sessionCounter int
	nodeId         int
	grpcAddr       string
}

type MembershipListType struct {
	members []MemberType
	// Potential global config about MembershipList
	myIP     net.IP
	myIndex  int
	MyPort   int
	MyNodeId int
}

func (s *MemberType) SessionCounter() int {
	return s.sessionCounter
}

func (s *MemberType) NodeId() int {
	return s.nodeId
}

func (s *MemberType) Addr() string {
	// TODO: naive version of TCP port = UDP port - 1000
	if s.grpcAddr == "" {
		str := s.addr
		re := regexp.MustCompile("(.*):(\\d*)")
		matches := re.FindStringSubmatch(str)
		if len(matches) >= 3 {
			if udpPort, err := strconv.Atoi(matches[2]); err != nil {
				log.Panicln("Invalid port number")
			} else {
				return matches[1] + ":" + strconv.Itoa(udpPort-1000)
			}
		}
	}
	return s.grpcAddr
}

func PosMod(a, b int) int {
	return (a%b + b) % b
}

func (ml *MembershipListType) insert(index int, memberType MemberType) {
	// TODO: Sdfs re-replicate
	log.Println("Member Added: ", memberType.addr)
	s := &ml.members
	*s = append(*s, MemberType{})
	copy((*s)[index+1:], (*s)[index:])
	(*s)[index] = memberType
	ml.updateMyIndex()
	ml.sort()
}

func (ml *MembershipListType) delete(index int) {
	// Sdfs re-replicate
	failId := ml.members[index].nodeId
	log.Println("channel to send: ", failId)
	sdfs2fd.Communicate <- failId

	log.Println("Member Rmved: ", ml.members[index].addr)
	s := &ml.members
	copy((*s)[index:], (*s)[index+1:])
	*s = (*s)[:len(*s)-1]
	ml.updateMyIndex()
	ml.sort()
}

func (ml *MembershipListType) insertNewID(id string, sessionID int, nodeId int) {
	for i, member := range MembershipList.members {
		if id == member.addr {
			if sessionID > member.sessionCounter {
				ml.members[i].sessionCounter = sessionID
			}
			return
		}
	}
	ml.insert(len(MembershipList.members),
		MemberType{addr: id, sessionCounter: sessionID, nodeId: nodeId})
}

func (ml *MembershipListType) lookupID(id string) (MemberType, bool) {
	for _, member := range MembershipList.members {
		if id == member.addr {
			return member, true
		}
	}
	return MemberType{}, false
}

func (ml *MembershipListType) sort() {
	sort.Slice(MembershipList.members, func(i, j int) bool {
		return MembershipList.members[i].nodeId < MembershipList.members[j].nodeId
	})
}

func (ml *MembershipListType) searchIndexById(key int) int {
	index := 0
	for i, member := range MembershipList.members {
		if member.nodeId >= key {
			index = i
			break
		}
	}
	return index
}

/*
API to SDFS
*/
func NextNofId(n, key int) MemberType {
	index := MembershipList.searchIndexById(key)
	index = (index + n) % len(MembershipList.members)
	return MembershipList.members[index]
}

func GetKeysOfId(nodeId int) []int {
	ml := &MembershipList
	index := ml.searchIndexById(nodeId)
	retArr := make([]int, 0)
	prevId := MembershipList.members[PosMod(index-1, len(MembershipList.members))].nodeId
	for key := index; key != prevId; {
		retArr = append(retArr, key)
		key = PosMod(key-1, RingSize)
	}
	return retArr
}

func GetDistByKey(from, to int) int {
	fromInd := MembershipList.searchIndexById(from)
	toInd := MembershipList.searchIndexById(to)
	if toInd < fromInd {
		toInd += len(MembershipList.members)
	}
	return toInd - fromInd
}

func PrevKOfKey(k, key int) int {
	if len(MembershipList.members) <= k {
		return key
	}
	index := MembershipList.searchIndexById(key)
	index = PosMod(index-k, len(MembershipList.members))
	return MembershipList.members[index].nodeId
}

func (ml *MembershipListType) deleteID(id string, sessionID int) {
	for i := range MembershipList.members {
		member := &MembershipList.members[i]
		if member.addr == id {
			if member.sessionCounter <= sessionID {
				ml.delete(i)
			}
			return
			//} else if member.addr > id {
			//	return
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
		if MembershipList.MyNodeId == -1 {
			MembershipList.MyNodeId = GetNodeIdFromHostname()
		}
		MyAddr = MembershipList.myIP.String() + ":" + strconv.Itoa(MembershipList.MyPort)
		AddSelfToList(0, MembershipList.MyNodeId)
	}
}

func AddSelfToList(sessionCounter int, nodeId int) {
	MembershipList.insertNewID(MyAddr, sessionCounter, nodeId)
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

func GetNodeIdFromHostname() int {
	if hName, err := os.Hostname(); err != nil {
		ErrHandler(err)
		return 0
	} else {
		re := regexp.MustCompile("fa18-cs425-g44-(\\d{2})\\.cs\\.illinois\\.edu")
		matches := re.FindStringSubmatch(hName)
		if len(matches) >= 2 {
			strId := matches[1]
			if id, err := strconv.Atoi(strId); err != nil {
				//ErrHandler(err)
			} else {
				return id
			}
		}
	}
	fmt.Println("Local debug purpose only")
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int()%RingSize + 1
}

func GetListElement() ([]string, []int, []int) {
	ret1 := make([]string, len(MembershipList.members))
	ret2 := make([]int, len(MembershipList.members))
	ret3 := make([]int, len(MembershipList.members))
	for i := 0; i < len(MembershipList.members); i++ {
		ret1[i] = MembershipList.members[i].addr
		ret2[i] = MembershipList.members[i].sessionCounter
		ret3[i] = MembershipList.members[i].nodeId
	}

	return ret1, ret2, ret3
}

func DumpTable() {
	table, _, _ := GetListElement()
	for i, t := range table {
		fmt.Printf("Index %d is Process: %s\n", i, t)
	}
}

func FormDumpTable() string {
	var response string
	table, _, tid := GetListElement()
	for i, t := range table {
		response += fmt.Sprintf("NodeId %d is Process: %s\n", tid[i], t)
	}
	return response
}
