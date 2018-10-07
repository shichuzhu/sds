package membership

import (
	"fmt"
	"net"
)

type Message struct {
	Header  string
	IP      string
	Session int
}

func Parse_Message(m Message) {
	switch m.Header {
	case "Ping":
		Handle_ping_message(m)
	case "Ack":
		Handle_ack_message(m)
	case "Update":
		Handle_update_message(m)
	default:
		fmt.Println("Wrong input message")
	}
}

func Handle_ping_message(m Message) {
	//create sending message from protool buffer
	// marshal the message to byte  string
	var mess []byte

	// We design to send UDP message
	UDP_send(m.IP, mess, 2)
}

func Handle_ack_message(m Message) {
	// mark the receiving status to true
}

func Handle_update_message(m Message) {
	// delete the member if exist
	// prepare the message to send (add to sequence if need)
	// UDP send message to randomly chosen target
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
