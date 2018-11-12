package memlist

import (
	"fa18cs425mp/src/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

type Message struct {
	Header  string
	IP      string
	Session int
}

func ParseMessage(m pb.DetectorMessage) {
	switch m.GetHeader() {
	case "Ping":
		HandlePingMessage(m)
	case "Ack":
		HandleAckMessage(m)
	case "Join":
		HandleJoinMessage(m)
	case "Delete":
		HandleDeleteMessage(m)
	case "NewJoin":
		HandleNewJoinMessage(m)
	default:
		log.Println("Wrong input message")
	}
}

func HandlePingMessage(m pb.DetectorMessage) {
	//create sending message from proto buffer
	// marshal the message to byte string
	addr := m.GetMem().GetAddr()
	ackMess := pb.DetectorMessage{Header: "Ack", Mem: &pb.Member{Addr: MyAddr}}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &ackMess}
	mess, _ := proto.Marshal(&UdpMess)
	// We design to send UDP message
	UdpSend(addr, mess, 1)
}

func HandleAckMessage(m pb.DetectorMessage) {
	// mark the receiving status to true
	addr := m.GetMem().GetAddr()

	for i := 0; i < len(ackWaitEntries); i++ {
		if addr == ackWaitEntries[i].addr {
			ackWaitEntries[i].ack()
		}
	}
}

func mem2pb(m *MemberType) *pb.Member {
	mem := new(pb.Member)
	mem.NodeId = int32(m.nodeId)
	mem.Addr = m.addr
	mem.SessNum = int32(m.sessionCounter)
	mem.GrpcAddr = m.grpcAddr
	return mem
}

func pb2mem(m *pb.Member) *MemberType {
	mem := new(MemberType)
	mem.nodeId = int(m.NodeId)
	mem.addr = m.Addr
	mem.sessionCounter = int(m.SessNum)
	mem.grpcAddr = m.GrpcAddr
	return mem
}

func HandleJoinMessage(m pb.DetectorMessage) {
	MemList.insertNewID(pb2mem(m.GetMem()))

	TTL := m.GetTtl()
	forwardMess := pb.DetectorMessage{
		Header: "Join", Ttl: TTL + 1, Mem: m.GetMem()}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &forwardMess}

	mess, _ := proto.Marshal(&UdpMess)
	if TTL < 4 {
		targets := MemList.getRandomTargets(3)

		for _, target := range targets {
			UdpSend(target, mess, 1)
		}
	}
}

func HandleDeleteMessage(m pb.DetectorMessage) {
	// delete the member if exist
	// prepare the message to send (add to sequence if need)
	// UDP send message to randomly chosen target
	addr := m.GetMem().GetAddr()
	MemList.deleteID(addr, int(m.GetMem().GetSessNum()))

	TTL := m.GetTtl()
	forwardMess := pb.DetectorMessage{Header: "Delete", Mem: m.GetMem(), Ttl: TTL + 1}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &forwardMess}
	mess, _ := proto.Marshal(&UdpMess)
	log.Println("Failure PROPAGATED: ", addr)
	if TTL < 4 {
		targets := MemList.getRandomTargets(3)

		for _, target := range targets {
			UdpSend(target, mess, 1)
		}
	}
}

/*
Introducer initialize the join message and gossip out
*/
func HandleNewJoinMessage(m pb.DetectorMessage) {
	addr := m.GetMem().GetAddr()
	// TODO: introducer assign nodeID possible to implement here
	MemList.insertNewID(pb2mem(m.GetMem())) // TODO: change sessionID
	fm := GetMemberListMessage()

	UdpMess := pb.UDPMessage{MessageType: "FullMembershipList", Fm: &fm}
	mess, _ := proto.Marshal(&UdpMess)

	UdpSend(addr, mess, 1)

	ackMess := pb.DetectorMessage{Header: "Join", Mem: m.GetMem(), Ttl: 0}
	UdpMess = pb.UDPMessage{MessageType: "DetectorMessage", Dm: &ackMess}
	mess, _ = proto.Marshal(&UdpMess)

	targets := MemList.getRandomTargets(3)
	for _, target := range targets {
		UdpSend(target, mess, 1)
	}
}

func ErrHandler(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}

func receiverService() {
	for {
		UdpMess, _ := UdpRecvSingle()
		if UdpMess == nil {
			continue
		}
		switch UdpMess.GetMessageType() {
		case "DetectorMessage":
			ParseMessage(*UdpMess.GetDm())
		case "FullMembershipList":
			InitialNewList(*UdpMess.GetFm())
		}
	}
	return
}

/*
	Initial new membership list at first join
*/
func InitialNewList(L pb.FullMembershipList) {
	ML := L.GetList()

	for _, member := range ML {
		MemList.insertNewID(pb2mem(member))
	}
}

func GetMemberListMessage() pb.FullMembershipList {
	List := make([]*pb.Member, len(MemList.members))
	for i := 0; i < len(MemList.members); i++ {
		List[i] = &pb.Member{Addr: MemList.members[i].addr,
			SessNum:  int32(MemList.members[i].sessionCounter),
			GrpcAddr: MemList.members[i].grpcAddr,
			NodeId:   int32(MemList.members[i].nodeId)}
	}
	fm := pb.FullMembershipList{List: List}
	return fm
}
