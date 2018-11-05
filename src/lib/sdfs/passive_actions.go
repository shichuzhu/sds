package sdfs

import "fa18cs425mp/src/lib/membership"

// TODO: This needs to be async, or to collect actions and return immediately to update ml
func ReReplicateUponFailure(failId int) {
	myId := membership.MembershipList.MyNodeId
	//lostKeyList := membership.GetKeysOfId(failId)
	if dist := membership.GetDistByKey(failId, myId); dist <= REPLICA {
		pullId := membership.PrevKOfKey(REPLICA, myId)
		for _, key := range membership.GetKeysOfId(pullId) {
			PullKeyFromNode(key, pullId)
		}
	}
}
