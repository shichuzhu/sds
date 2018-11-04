package sdfs

import "container/list"

var fileMap map[string]*list.List

func memTableIntial() {
	fileMap = make(map[string]*list.List)
}

func getFileVersion(fileName string) int {
	list, ok := fileMap[fileName]

}
