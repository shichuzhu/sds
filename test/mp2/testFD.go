package main

import (
	ms "fa18cs425mp/src/lib/membership"
	"flag"
	"fmt"
	"time"
)

func main() {
	introFlag := flag.Bool("i", false, "Indicator of the node is the introducer of port 11000")
	port := flag.Int("p", 11000, "port number to use for the failure detector")
	drop := flag.Float64("d", 0.0, "Simulated packet drop rate.")
	ms.PacketDrop.SetDropRate(float32(*drop))
	flag.Parse()
	if !(*introFlag) {
		ms.MembershipList.MyPort = *port
		ms.ContactIntroducer("10.195.35.3:11000")
	} else {
		ms.MembershipList.MyPort = *port
		ms.InitInstance()
		ms.AddSelfToList(0)
		ms.StartFailureDetector()
	}
	for {
		ms.DumpTable()
		time.Sleep(time.Duration(2000) * time.Millisecond)
		fmt.Println(time.Now())
		fmt.Printf("Average usage is %v \n", ms.NetworkStats.GetBandwidthUsage())
	}
}
