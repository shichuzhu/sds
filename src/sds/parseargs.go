package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ArrayFlags []int

func (i *ArrayFlags) String() string {
	return fmt.Sprint(*i)
}

func (i *ArrayFlags) Set(value string) error {
	strArr := strings.Split(value, ",")
	for _, strItem := range strArr {
		if strItem != "" {
			if intItem, err := strconv.Atoi(strItem); err != nil {
				log.Panicln("Can't parse command-line arguments.")
			} else {
				*i = append(*i, intItem)
			}
		}
	}
	return nil
}

func RegisterNodeArgs(flagSet *flag.FlagSet, serverIds *ArrayFlags) *flag.FlagSet {
	if flagSet == nil {
		flagSet = flag.NewFlagSet("dclient", flag.ErrorHandling(flag.PanicOnError))
	}
	flagSet.Var(serverIds, "n", "Node list, comma separated string")
	return flagSet
}

func ParseArgs(flagSet *flag.FlagSet, argvp *[]string) *flag.FlagSet {
	var serverIds ArrayFlags

	argv := *argvp
	argv = argv[1:]
	flagSet = RegisterNodeArgs(flagSet, &serverIds)
	if err := flagSet.Parse(argv); err != nil {
		return nil
	}
	*argvp = argv[len(argv)-flagSet.NArg():]
	return flagSet
}
