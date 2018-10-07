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

func Parse_Message(m pb.DetectorMessage) {
	switch m.GetHeader() {
	case "Ping":
		Handle_ping_message(m)
	case "Ack":
		Handle_ack_message(m)
	case "Join":
		Handle_join_message(m)
	case "Delete":
		Handle_delete_message(m)
	default:
		fmt.Println("Wrong input message")
	}
}

func Handle_ping_message(m pb.DetectorMessage) {
	//create sending message from protool buffer
	// marshal the message to byte  string
	addr := m.GetAddr()
	ackMess := pb.DetectorMessage{"Ack", MyAddr, 0, 0}
	mess, _ := proto.Marshal(&ackMess)
	// We design to send UDP message
	UDP_send(addr, mess, 2)
}

func Handle_ack_message(m pb.DetectorMessage) {
	// mark the receiving status to true

}

func Handle_join_message(m pb.DetectorMessage) {
	addr := m.GetAddr()
	session := m.GetSessNUm()
	MembershipList.insertNewID(addr, int(session))

	TTL := m.GetTTL()
	forwardMess := pb.DetectorMessage{"Join", addr, session, TTL + 1}
	mess, _ := proto.Marshal(&forwardMess)
	if TTL < 4 {
		targets := MembershipList.getRandomTargets(3)

		for _, target := range targets {
			UDP_send(target, mess, 2)
		}
	}
}

func Handle_delete_message(m pb.DetectorMessage) {
	// delete the member if exist
	// prepare the message to send (add to sequence if need)
	// UDP send message to randomly chosen target
	addr := m.GetAddr()
	session := m.GetSessNUm()
	MembershipList.deleteID(addr, int(session))

	TTL := m.GetTTL()
	forwardMess := pb.DetectorMessage{"Delete", addr, session, TTL + 1}
	mess, _ := proto.Marshal(&forwardMess)
	if TTL < 4 {
		targets := MembershipList.getRandomTargets(3)

		for _, target := range targets {
			UDP_send(target, mess, 2)
		}
	}
}

func UDP_send(IP string, buf []byte, rep int) {
	addr, _ := net.ResolveUDPAddr("udp", IP)
	conn, _ := net.DialUDP("udp", nil, addr)
	for i := 0; i < rep; i++ {
		_, err := conn.WriteToUDP(buf, addr)
		Err_handler(err)
	}
	conn.Close()
}

func Err_handler(err error) {
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
	buf := make([]byte, 1024)
	defer fd.Close()
	for {
		n, _, err := fd.ReadFrom(buf)
		if err != nil {
			fmt.Printf("ReadError from UDP: %s", err.Error())
		}
		DetMess := pb.DetectorMessage{}
		err = proto.Unmarshal(buf[0:n], &DetMess)

		if err != nil {
			fmt.Printf("Error in Unmarshal proto message: %s", err.Error())
			return
		}

		Parse_Message(DetMess)
	}

	return
}
