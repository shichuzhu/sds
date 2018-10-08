package parseargs

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type arrayFlags []int

var ArgsCopy []string

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

var alreadyParsed bool
var ServerIds arrayFlags

func RegisterNodeArgs(flagSet *flag.FlagSet) *flag.FlagSet {
	if flagSet == nil {
		flagSet = flag.NewFlagSet("dclient", flag.ErrorHandling(flag.PanicOnError))
	}
	flagSet.Var(&ServerIds, "n", "Node list, comma separated string")
	return flagSet
}

func ParseArgs(flagSet *flag.FlagSet, name string) bool {
	if !alreadyParsed {
		ArgsCopy = make([]string, len(os.Args))
		copy(ArgsCopy, os.Args)
		alreadyParsed = true
	} else if name == ArgsCopy[0] {
		return false
	}
	flagSet.Parse(ArgsCopy[1:])
	ArgsCopy = ArgsCopy[len(ArgsCopy)-flagSet.NArg():]
	return true
}
