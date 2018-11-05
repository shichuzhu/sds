package main

import (
	"container/list"
	"fmt"
)

var fileMap map[string]int

func main() {
	MemTableIntial()
	version := GetFileVersion("hello.txt")
	fmt.Println(version)
	InsertFileVersion("hello.txt", 1)
	version = GetFileVersion("hello.txt")
	deleteFileFromTable("hello.txt")
	version = GetFileVersion("hello.txt")
	fmt.Println(version)
}

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

func listAllFile() list.List {
	ret := list.New()
	for key := range fileMap {
		str := fmt.Sprintf(key+" %d", fileMap[key])
		ret.PushBack(str)
	}

	return *ret
}

func deleteFileFromTable(fileName string) {
	delete(fileMap, fileName)
}

/*func getFileFoeKey(key int) []FileVersionPair {
	ret := []FileVersionPair{}
	for k := range fileMap {
		if HashToKey(k) == key {
			pair := FileVersionPair{k, fileMap[k].Front().Value.(int)}
			ret = append(ret, pair)
		}
	}
	return ret
}*/
