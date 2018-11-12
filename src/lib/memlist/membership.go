package memlist

import (
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/shared/sdfs2fd"
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

func (s *MemberType) SessionCounter() int {
	return s.sessionCounter
}

func (s *MemberType) NodeId() int {
	return s.nodeId
}

func (s *MemberType) Addr() string {
	return s.grpcAddr
}

func (s *MemberType) SetTcpAddrWithPort(tcpPort int) {
	str := s.addr
	re := regexp.MustCompile("(.*):(\\d*)")
	matches := re.FindStringSubmatch(str)
	if len(matches) >= 3 {
		s.grpcAddr = matches[1] + ":" + strconv.Itoa(tcpPort)
	} else {
		log.Fatalln("Cannot set Tcp address")
	}
}

func (ml *MembershipListType) insert(index int, memberType *MemberType) {
	// TODO: Sdfs re-replicate
	log.Println("Member Added: ", memberType.addr)
	s := &ml.members
	*s = append(*s, MemberType{})
	copy((*s)[index+1:], (*s)[index:])
	(*s)[index] = *memberType
	ml.updateMyIndex()
	ml.sort()
}

func (ml *MembershipListType) delete(index int) {
	// Sdfs re-replicate
	failId := ml.members[index].nodeId
	//log.Println("channel to send: ", failId)
	sdfs2fd.Fd2Sdfs <- failId

	<-sdfs2fd.Sdfs2Fd // Barrier to avoid read-write conflict
	log.Println("Member Rmved: ", ml.members[index].addr)
	s := &ml.members
	copy((*s)[index:], (*s)[index+1:])
	*s = (*s)[:len(*s)-1]
	ml.updateMyIndex()
	ml.sort()
	sdfs2fd.Fd2Sdfs <- 1
}

func (ml *MembershipListType) insertNewID(m *MemberType) {
	lk.Lock()
	defer lk.Unlock()
	id := m.addr
	sessionID := m.sessionCounter
	for i, member := range MemList.members {
		if id == member.addr {
			if sessionID > member.sessionCounter {
				ml.members[i].sessionCounter = sessionID
			}
			return
		}
	}
	ml.insert(len(MemList.members), m)
}

func (ml *MembershipListType) lookupID(id string) (MemberType, bool) {
	for _, member := range MemList.members {
		if id == member.addr {
			return member, true
		}
	}
	return MemberType{}, false
}

func (ml *MembershipListType) sort() {
	sort.Slice(MemList.members, func(i, j int) bool {
		return MemList.members[i].nodeId < MemList.members[j].nodeId
	})
}

func (ml *MembershipListType) searchIndexById(key int) int {
	index := 0
	for i, member := range MemList.members {
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
func NextNofId(n, key int) *MemberType {
	index := MemList.searchIndexById(key)
	index = (index + n) % len(MemList.members)
	return &MemList.members[index]
}

func GetKeysOfId(nodeId int) []int {
	ml := &MemList
	index := ml.searchIndexById(nodeId)
	retArr := make([]int, 0)
	prevId := MemList.members[utils.PosMod(index-1, len(MemList.members))].nodeId
	for key := index; key != prevId; {
		retArr = append(retArr, key)
		key = utils.PosMod(key-1, RingSize)
	}
	return retArr
}

func GetDistByKey(from, to int) int {
	fromInd := MemList.searchIndexById(from)
	toInd := MemList.searchIndexById(to)
	if toInd < fromInd {
		toInd += len(MemList.members)
	}
	return toInd - fromInd
}

func PrevKOfKey(k, key int) int {
	if len(MemList.members) <= k {
		return key
	}
	index := MemList.searchIndexById(key)
	index = utils.PosMod(index-k, len(MemList.members))
	return MemList.members[index].nodeId
}

func (ml *MembershipListType) deleteID(id string, sessionID int) {
	lk.Lock()
	defer lk.Unlock()
	for i := range MemList.members {
		member := &MemList.members[i]
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
	num = utils.Min(num, len(MemList.members))
	targets := make([]string, num)
	for i, j := range rand.Perm(len(MemList.members))[:num] {
		targets[i] = ml.members[j].addr
	}
	return targets
}

func (ml *MembershipListType) updateMyIndex() {
	for i, member := range MemList.members {
		if member.addr == MyAddr {
			ml.myIndex = i
			return
		}
	}
	log.Fatalln("Self not existing in membership list. Aborting...")
}

func (ml *MembershipListType) getPingTargets(num int) []string {
	num = utils.Max(0, utils.Min(NodeNumberToPing, len(MemList.members)-1))
	targets := make([]string, num)
	for i := range targets {
		targets[i] = ml.members[(ml.myIndex+i+1)%len(ml.members)].addr
	}
	return targets
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
	fmt.Println("Cannot derive nodeID from host, random generating one")
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int()%RingSize + 1
}

func GetListElement() ([]string, []int, []int) {
	ret1 := make([]string, len(MemList.members))
	ret2 := make([]int, len(MemList.members))
	ret3 := make([]int, len(MemList.members))
	for i := 0; i < len(MemList.members); i++ {
		ret1[i] = MemList.members[i].addr
		ret2[i] = MemList.members[i].sessionCounter
		ret3[i] = MemList.members[i].nodeId
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
