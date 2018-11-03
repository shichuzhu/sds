package membership

import "fmt"

const FailureTimeout int = 1800 // in millisecond
const DefaultUdpPort int = 11000
const NodeNumberToPing = 3
const MultiSendNumber = 1
const RingSize = 10

func ListInfo() {
	DumpTable() // Will print to stdout
	fmt.Println("Current process ID is ", MyAddr)
}

func FormListInfo() string {
	var response string
	response += FormDumpTable()
	response += fmt.Sprintf("Current process NodeId %d is %s\n", MembershipList.MyNodeId, MyAddr)
	return response
}

func JoinByIntroducer(introAddr string) {
	ContactIntroducer(introAddr)
}

func LeaveGroup() {
	ReportFailure(MyAddr)
}
