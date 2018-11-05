package sdfs

import (
	"container/list"
	"fmt"
)

var fileMap map[string]int

type nameVersionPair struct {
	fileName     string
	lastVersions int
}

type FileVersionPair struct {
	name  string
	index int
}

func MemTableIntial() {
	fileMap = make(map[string]int)
}

func GetFileVersion(fileName string) int {
	version, present := fileMap[fileName]
	if !present {
		return 0
	}
	return version
}

func InsertFileVersion(fileName string, version int) int {
	fileMap[fileName] = version
	return version
}

func ListAllFile() list.List {
	ret := list.New()
	for key := range fileMap {
		str := fmt.Sprintf(key+" %d", fileMap[key])
		ret.PushBack(str)
	}

	return *ret
}

func DeleteFileFromTable(fileName string) {
	delete(fileMap, fileName)
}

func getFileFoeKey(key int) []FileVersionPair {
	ret := []FileVersionPair{}
	for k := range fileMap {
		if HashToKey(k) == key {
			pair := FileVersionPair{k, fileMap[k]}
			ret = append(ret, pair)
		}
	}
	return ret
}
