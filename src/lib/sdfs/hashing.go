package sdfs

import (
	ms "fa18cs425mp/src/lib/membership"
	"hash/fnv"
	"log"
	"strconv"
)

var RUNES = []rune("01/V")

type Builder struct {
	arr []rune
	loc int
}

func (b *Builder) Grow(n int) {
	b.arr = make([]rune, n)
	b.loc = 0
}

func (b *Builder) WriteRune(r rune) {
	b.arr[b.loc] = r
	b.loc++
}

func (b *Builder) WriteString(s string) {
	for _, c := range s {
		b.WriteRune(c)
	}
}

func (b *Builder) String() string {
	return string(b.arr)
}

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
	lfn := &Builder{}
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
	sfn := &Builder{}
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
