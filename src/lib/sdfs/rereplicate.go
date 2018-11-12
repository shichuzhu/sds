package sdfs

import (
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/shared/sdfs2fd"
	"log"
)

func ReReplicateHandler() {
	for {
		failId, more := <-sdfs2fd.Fd2Sdfs
		if more {
			ReReplicateUponFailure(failId)
		} else {
			return
		}
	}
}

func ReReplicateUponFailure(failId int) {
	fetchedKeys := GetFetchKeys(failId)
	sdfs2fd.Sdfs2Fd <- 1 // Sync the FD to remove node
	<-sdfs2fd.Fd2Sdfs    // Wait until deletion completed
	if len(fetchedKeys) > 0 {
		FetchKeys(fetchedKeys)
	}
}

func GetFetchKeys(failId int) []int {
	myId := memlist.MemList.MyNodeId
	//lostKeyList := membership.GetKeysOfId(failId)
	if dist := memlist.GetDistByKey(failId, myId); dist <= REPLICA {
		pullId := memlist.PrevKOfKey(REPLICA, myId)
		return memlist.GetKeysOfId(pullId)
	}
	return nil
}

func FetchKeys(keys []int) {
	for _, key := range keys {
		for suc := 0; suc < 4; suc++ {
			pullId := FindNodeId(key, suc)
			if pullId == memlist.MemList.MyNodeId {
				log.Printf("Key %d exists locally already, skip fetching", key)
				break
			}
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
