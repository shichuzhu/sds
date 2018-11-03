package main

import (
	"os"
	"regexp"
	"strconv"
)

func main() {

	if hName, err := os.Hostname(); err != nil {
		return
	} else {
		re := regexp.MustCompile("fa18-cs425-g44-(\\d{2})\\.cs\\.illinois\\.edu")
		strId := re.FindStringSubmatch(hName)[1]
		if id, err := strconv.Atoi(strId); err != nil {
			return
		} else {
			println(id)
		}
	}
}
