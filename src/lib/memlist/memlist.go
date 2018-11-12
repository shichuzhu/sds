package memlist

import (
	"fa18cs425mp/src/shared/cfg"
	"flag"
	"net"
	"strconv"
	"sync"
)

var MyAddr string
var MemList MembershipListType
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
	// Potential global config about MemList
	myIP     net.IP
	myIndex  int
	MyPort   int
	MyNodeId int
}

var (
	port    = flag.Int("pfd", cfg.Cfg.DefaultUDPPort, "port number to use for the failure detector")
	TcpPort = flag.Int("port", cfg.Cfg.DefaultTCPPort, "The server port")
	drop    = flag.Float64("drop", 0.0, "Simulated packet drop rate.")
	nodeId  = flag.Int("nodeid", -1, "The nodeid if not randomized")
	lk      = sync.Mutex{}
)

func SetupLocalMemList() {
	MemList.MyPort = *port
	MemList.myIP = GetOutboundIP()
	MyAddr = MemList.myIP.String() + ":" + strconv.Itoa(MemList.MyPort)
	ackWaitEntries = make([]AckWaitEntry, NodeNumberToPing)
	if *nodeId != -1 {
		MemList.MyNodeId = *nodeId
	} else {
		MemList.MyNodeId = GetNodeIdFromHostname()
	}

	tmp := &MemberType{
		addr:           MyAddr,
		sessionCounter: 0,
		nodeId:         MemList.MyNodeId,
		grpcAddr:       MyAddr}
	tmp.SetTcpAddrWithPort(*TcpPort)
	MemList.insertNewID(tmp)

	PacketDrop.SetDropRate(float32(*drop))
}
