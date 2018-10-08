package membership

import "fmt"

const FailureTimeout int = 1800 // in millisecond
const DefaultUdpPort int = 11000
const NodeNumberToPing = 3
const MultiSendNumber = 1

func ListInfo() {
	DumpTable() // Will print to stdout
	fmt.Println("Current process ID is ", MyAddr)
}

func JoinByIntroducer(introAddr string) {
	ContactIntroducer(introAddr)
}

func LeaveGroup() {
	ReportFailure(MyAddr)
}
