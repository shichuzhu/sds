package sdfs

import "container/list"

var fileMap map[string]*list.List

func memTableIntial() {
	fileMap = make(map[string]*list.List)
}

func getFileVersion(fileName string) int {
	list, present := fileMap[fileName]
	if !present {
		return 0
	}

	return list.Front().Value.(int)
}

func insertFileVersion(fileName string, version int) int {
	_, present := fileMap[fileName]
	if !present {
		fileMap[fileName] = list.New()
		fileMap[fileName].PushFront(version)
		return 1
	}

	list, _ := fileMap[fileName]
	if list.Len() == 5 {
		/*
			TODO: delete the oldest file here and insert new version
		*/
		list.Remove(list.Back())
		list.PushFront(version)
		return 1
	}
	list.PushFront(version)
	return 1
}

func listAllFile(key int) list.List {
	ret := list.New()
	for k := range fileMap {
		if HashToKey(k) == key {
			ret.PushBack(k)
		}
	}

	return *ret
}

func deleteFileFromTable(fileName string) {
	delete(fileMap, fileName)
}
