package main

import (
	"bufio"
	ms "fa18cs425mp/src/lib/membership"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	introFlag := flag.Bool("i", false, "Indicator of the node is the introducer of port 11000")
	port := flag.Int("p", 11000, "port number to use for the failure detector")
	drop := flag.Float64("d", 0.0, "Simulated packet drop rate.")
	flag.Parse()

	ms.PacketDrop.SetDropRate(float32(*drop))
	if !(*introFlag) {
		ms.MembershipList.MyPort = *port
	} else {
		ms.MembershipList.MyPort = *port
		ms.InitInstance()
		ms.StartFailureDetector()
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "ls":
			ms.ListInfo()
		case "leave":
			ms.LeaveGroup()
		default:
			if len(text) >= 4 && text[:4] == "join" {
				introAddr := text[5:]
				ms.JoinByIntroducer(introAddr)
			} else {
				fmt.Println("Invalid input.")
			}
		}
	}
}
