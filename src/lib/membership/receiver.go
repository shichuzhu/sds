package membership

import (
	pb "fa18cs425mp/src/protobuf"
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
	addr := m.GetAddr()
	ackMess := pb.DetectorMessage{Header: "Ack", Addr: MyAddr}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &ackMess}
	mess, _ := proto.Marshal(&UdpMess)
	// We design to send UDP message
	UdpSend(addr, mess, 1)
}

func HandleAckMessage(m pb.DetectorMessage) {
	// mark the receiving status to true
	addr := m.GetAddr()

	for i := 0; i < len(ackWaitEntries); i++ {
		if addr == ackWaitEntries[i].addr {
			ackWaitEntries[i].ack()
		}
	}
}

func HandleJoinMessage(m pb.DetectorMessage) {
	addr := m.GetAddr()
	session := m.GetSessNum()
	nodeId := m.GetNodeId()
	MembershipList.insertNewID(addr, int(session), int(nodeId))

	TTL := m.GetTtl()
	forwardMess := pb.DetectorMessage{
		Header: "Join", Addr: addr, SessNum: session, Ttl: TTL + 1,
		NodeId: nodeId,
	}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &forwardMess}

	mess, _ := proto.Marshal(&UdpMess)
	if TTL < 4 {
		targets := MembershipList.getRandomTargets(3)

		for _, target := range targets {
			UdpSend(target, mess, 1)
		}
	}
}

func HandleDeleteMessage(m pb.DetectorMessage) {
	// delete the member if exist
	// prepare the message to send (add to sequence if need)
	// UDP send message to randomly chosen target
	addr := m.GetAddr()
	session := m.GetSessNum()
	MembershipList.deleteID(addr, int(session))

	TTL := m.GetTtl()
	forwardMess := pb.DetectorMessage{Header: "Delete", Addr: addr, SessNum: session, Ttl: TTL + 1}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &forwardMess}
	mess, _ := proto.Marshal(&UdpMess)
	log.Println("Failure PROPAGATED: ", addr)
	if TTL < 4 {
		targets := MembershipList.getRandomTargets(3)

		for _, target := range targets {
			UdpSend(target, mess, 1)
		}
	}
}

/*
Introducer initialize the join message and gossip out
*/
func HandleNewJoinMessage(m pb.DetectorMessage) {
	addr := m.GetAddr()
	nodeId := m.GetNodeId()
	MembershipList.insertNewID(addr, 0, int(nodeId)) // TODO: change sessionID
	fm := GetMemberListMessage()

	UdpMess := pb.UDPMessage{MessageType: "FullMembershipList", Fm: &fm}
	mess, _ := proto.Marshal(&UdpMess)

	UdpSend(addr, mess, 1)

	ackMess := pb.DetectorMessage{Header: "Join", Addr: addr, SessNum: 0, Ttl: 0, NodeId: nodeId}
	UdpMess = pb.UDPMessage{MessageType: "DetectorMessage", Dm: &ackMess}
	mess, _ = proto.Marshal(&UdpMess)

	targets := MembershipList.getRandomTargets(3)
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
		MembershipList.insertNewID(member.GetAddr(), int(member.GetSessNum()), int(member.GetNodeId()))
	}
}

func GetMemberListMessage() pb.FullMembershipList {
	List := make([]*pb.Member, len(MembershipList.members))
	for i := 0; i < len(MembershipList.members); i++ {
		List[i] = &pb.Member{Addr: MembershipList.members[i].addr,
			SessNum: int32(MembershipList.members[i].sessionCounter),
			NodeId:  int32(MembershipList.members[i].nodeId)}
	}
	fm := pb.FullMembershipList{List: List}
	return fm
}
