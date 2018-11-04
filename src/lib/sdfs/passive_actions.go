package sdfs

import "fa18cs425mp/src/lib/membership"

func ReReplicateUponFailure(failId int) {
	myId := membership.MembershipList.MyNodeId
	lostKeyList := membership.GetKeysOfId(failId)
	if dist := membership.GetDistByKey(failId, myId); dist <= REPLICA {
		pullId := membership.PrevKOfKey(REPLICA, myId)
		PullKeyFromNode(membership.GetKeysOfId(pullId), pullId)
	}
}
