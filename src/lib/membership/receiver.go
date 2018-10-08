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
		log.Println("Wrong message input message")
	}
}

func HandlePingMessage(m pb.DetectorMessage) {
	//create sending message from proto buffer
	// marshal the message to byte string
	addr := m.GetAddr()
	ackMess := pb.DetectorMessage{Header: "Ack", Addr: MyAddr, SessNUm: 0, TTL: 0}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &ackMess}
	mess, _ := proto.Marshal(&UdpMess)
	// We design to send UDP message
	UdpSend(addr, mess, 2)
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
	session := m.GetSessNUm()
	MembershipList.insertNewID(addr, int(session))

	TTL := m.GetTTL()
	forwardMess := pb.DetectorMessage{Header: "Join", Addr: addr, SessNUm: session, TTL: TTL + 1}
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
	session := m.GetSessNUm()
	MembershipList.deleteID(addr, int(session))

	TTL := m.GetTTL()
	forwardMess := pb.DetectorMessage{Header: "Delete", Addr: addr, SessNUm: session, TTL: TTL + 1}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &forwardMess}
	mess, _ := proto.Marshal(&UdpMess)
	if TTL < 4 {
		targets := MembershipList.getRandomTargets(3)

		for _, target := range targets {
			UdpSend(target, mess, 1)
		}
	}
}

func HandleNewJoinMessage(m pb.DetectorMessage) {
	addr := m.GetAddr()
	MembershipList.insertNewID(addr, 0) // TODO: change sessionID
	fm := GetMemberListMessage()
	ackMess := pb.DetectorMessage{}

	UdpMess := pb.UDPMessage{MessageType: "FullMembershipList", Dm: &ackMess, Fm: &fm}
	mess, _ := proto.Marshal(&UdpMess)

	UdpSend(addr, mess, 4)

	ackMess = pb.DetectorMessage{Header: "Join", Addr: addr, SessNUm: 0, TTL: 0}
	UdpMess = pb.UDPMessage{MessageType: "DetectorMessage", Dm: &ackMess}
	mess, _ = proto.Marshal(&UdpMess)

	targets := MembershipList.getRandomTargets(3)
	for _, target := range targets {
		UdpSend(target, mess, 2)
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
		MembershipList.insertNewID(member.GetAddr(), int(member.GetSessNum()))
	}
}

func GetMemberListMessage() pb.FullMembershipList {
	List := make([]*pb.Member, len(MembershipList.members))
	for i := 0; i < len(MembershipList.members); i++ {
		List[i] = &pb.Member{Addr: MembershipList.members[i].addr,
			SessNum: int32(MembershipList.members[i].sessionCounter)}
	}
	fm := pb.FullMembershipList{List: List}

	return fm
}
