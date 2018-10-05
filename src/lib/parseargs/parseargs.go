package parseargs

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type arrayFlags []int

func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

func (intArr *arrayFlags) Set(value string) error {
	strArr := strings.Split(value, ",")
	for _, strItem := range strArr {
		if strItem != "" {
			if intItem, err := strconv.Atoi(strItem); err != nil {
				log.Panicln("Can't parse command-line arguments.")
			} else {
				*intArr = append(*intArr, intItem)
			}
		}
	}
	return nil
}

var ServerIds arrayFlags
var alreadyParsed bool

func registerArgs() {
	flag.Var(&ServerIds, "n", "Node list, comma separated string")
}

func ParseArgs() {
	if !alreadyParsed {
		registerArgs()
		flag.Parse()
		alreadyParsed = true
	}
}
