package membership

import (
	"net"
	"time"
)

// Use getter to get list for thread safety issue.
type MemberType struct {
	addr           string
	sessionCounter int
	deadline       time.Time
	ddlPending     bool
}

type MembershipListType struct {
	members []MemberType
	// Potential global config about MembershipList
	myID int
}

var MembershipList MembershipListType
var Xmtr *net.UDPConn // Transmitter
