package sdfs_test

import (
	"os"
	"regexp"
	"strconv"
	"testing"
)

func TestMemTableIntial(t *testing.T) {

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
