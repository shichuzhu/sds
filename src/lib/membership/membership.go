package membership

import "time"

type Member struct {
	addr           string
	sessionCounter int
	deadline       time.Time
	ddlPending     bool
}

type MembershipList struct {
	members []Member
	// Potential global config about MembershipList
}
