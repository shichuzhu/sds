package sdfs

import (
	ms "fa18cs425mp/src/lib/membership"
	"hash/fnv"
	"log"
	"strconv"
	"strings"
)

var RUNES = []rune("01/V")

func HashToKey(str string) int {
	h := fnv.New32()
	h.Write([]byte(str))
	return int(h.Sum32() % ms.RingSize)
}

func FindNodeId(key int, successor int) ms.MemberType {
	return ms.NextNofId(successor, key)
}

func SdfsToLfs(s string, v int) string {
	n := len(s)
	lfn := strings.Builder{}
	lfn.Grow(2*n + 2)
	for _, c := range s {
		if c != '/' {
			lfn.WriteRune(RUNES[0])
			lfn.WriteRune(c)
		} else {
			lfn.WriteRune(RUNES[1])
			lfn.WriteRune(RUNES[1])
		}
	}
	lfn.WriteString("V" + strconv.Itoa(v))
	return lfn.String()
}

func LfsToSdfs(localFilename string) (string, int) {
	s := []rune(localFilename)
	sfn := strings.Builder{}
	sfn.Grow(len(s) / 2)
	for i := 0; i < len(s); {
		if s[i] == RUNES[0] {
			sfn.WriteRune(s[i+1])
			i += 2
		} else if s[i] == RUNES[1] {
			sfn.WriteRune(RUNES[2])
			i += 2
		} else { // s[i] is 'V'
			if v, err := strconv.Atoi(string(s[i+1:])); err != nil {
				log.Panic(err)
			} else {
				return sfn.String(), v
			}
		}
	}
	log.Panic()
	return "", 0
}
