package sdfs

import (
	"fa18cs425mp/src/lib/membership"
	"log"
)

func ReReplicateUponFailure(failId int) {
	fetchedKeys := GetFetchKeys(failId)
	if len(fetchedKeys) > 0 {
		go FetchKeys(fetchedKeys)
	}
}

func GetFetchKeys(failId int) []int {
	myId := membership.MembershipList.MyNodeId
	//lostKeyList := membership.GetKeysOfId(failId)
	if dist := membership.GetDistByKey(failId, myId); dist <= REPLICA {
		pullId := membership.PrevKOfKey(REPLICA, myId)
		return membership.GetKeysOfId(pullId)
	}
	return nil
}

func FetchKeys(keys []int) {
	for _, key := range keys {
		for suc := 0; suc < 4; suc++ {
			pullId := FindNodeId(key, suc)
			if err := PullKeyFromNode(key, pullId); err != nil {
				log.Println(err)
				log.Printf("Fail to fetch key %d from the {%d}th successor, trying next", key, suc)
			} else {
				log.Printf("Fetched key %d from the {%d}th successor", key, suc)
				break
			}
		}
	}
	log.Printf("Rereplication process finished.")
}
