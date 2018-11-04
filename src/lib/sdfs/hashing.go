package sdfs

import (
	ms "fa18cs425mp/src/lib/membership"
	"hash/fnv"
)

func HashToKey(str string) int {
	h := fnv.New32()
	h.Write([]byte(str))
	return int(h.Sum32() % ms.RingSize)
}

func FindNodeId(key int, successor int) int {
	return ms.NextNofId(successor, key)
}
