package memlist

import (
	"fa18cs425mp/src/shared/cfg"
	"fmt"
)

var (
	FailureTimeout   = cfg.Cfg.Memlist.FailureTimeout // in millisecond
	NodeNumberToPing = cfg.Cfg.Memlist.NodeNumberToPing
	MultiSendNumber  = cfg.Cfg.Memlist.MultiSendNumber
	RingSize         = cfg.Cfg.Memlist.RingSize
)

func FormListInfo() string {
	var response string
	response += FormDumpTable()
	response += fmt.Sprintf("Current process NodeId %d is %s\n", MemList.MyNodeId, MyAddr)
	return response
}

func JoinByIntroducer(introAddr string) {
	ContactIntroducer(introAddr)
}

func LeaveGroup() {
	ReportFailure(MyAddr)
}
