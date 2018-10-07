package membership

import (
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"strconv"
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
		HandleNewjoinMessage(m)
	default:
		fmt.Println("Wrong input message")
	}
}

func HandlePingMessage(m pb.DetectorMessage) {
	//create sending message from protool buffer
	// marshal the message to byte  string
	addr := m.GetAddr()
	ackMess := pb.DetectorMessage{Header: "Ack", Addr: MyAddr, SessNUm: 0, TTL: 0}
	fm := pb.FullMembershipList{}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &ackMess, Fm: &fm}
	mess, _ := proto.Marshal(&UdpMess)
	// We design to send UDP message
	UDPSend(addr, mess, 2)
}

func HandleAckMessage(m pb.DetectorMessage) {
	// mark the receiving status to true

}

func HandleJoinMessage(m pb.DetectorMessage) {
	addr := m.GetAddr()
	session := m.GetSessNUm()
	MembershipList.insertNewID(addr, int(session))

	TTL := m.GetTTL()
	forwardMess := pb.DetectorMessage{Header: "Join", Addr: addr, SessNUm: session, TTL: TTL + 1}
	fm := pb.FullMembershipList{}
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &forwardMess, Fm: &fm}

	mess, _ := proto.Marshal(&UdpMess)
	if TTL < 4 {
		targets := MembershipList.getRandomTargets(3)

		for _, target := range targets {
			UDPSend(target, mess, 2)
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
	UdpMess := pb.UDPMessage{MessageType: "DetectorMessage", Dm: &forwardMess, Fm: &pb.FullMembershipList{}}
	mess, _ := proto.Marshal(&UdpMess)
	if TTL < 4 {
		targets := MembershipList.getRandomTargets(3)

		for _, target := range targets {
			UDPSend(target, mess, 2)
		}
	}
}

func HandleNewjoinMessage(m pb.DetectorMessage) {
	addr := m.GetAddr()
	fm := GetMemberListMessage()
	ackMess := pb.DetectorMessage{}

	UdpMess := pb.UDPMessage{MessageType: "FullMembershipList", Dm: &ackMess, Fm: &fm}
	mess, _ := proto.Marshal(&UdpMess)

	UDPSend(addr, mess, 4)

}

func UDPSend(IP string, buf []byte, rep int) {
	addr, _ := net.ResolveUDPAddr("udp", IP)
	conn, _ := net.DialUDP("udp", nil, addr)
	for i := 0; i < rep; i++ {
		_, err := conn.WriteToUDP(buf, addr)
		ErrHandler(err)
	}
	conn.Close()
}

func ErrHandler(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func receiverService() {
	// for
	fd, err := net.ListenPacket("udp", ":"+strconv.Itoa(DefaultUdpPort))
	if err != nil {
		log.Panicln("Cannot listen UDP packet")
	}
	buf := make([]byte, 4096)
	defer fd.Close()
	for {
		n, _, err := fd.ReadFrom(buf)
		if err != nil {
			fmt.Printf("ReadError from UDP: %s", err.Error())
		}
		UdpMess := pb.UDPMessage{}
		err = proto.Unmarshal(buf[0:n], &UdpMess)

		if err != nil {
			fmt.Printf("Error in Unmarshal proto message: %s", err.Error())
			return
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
		List[i].Addr = MembershipList.members[i].addr
		List[i].SessNum = int32(MembershipList.members[i].sessionCounter)
	}
	fm := pb.FullMembershipList{List: List}

	return fm
}
