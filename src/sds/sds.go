/**
Package sds stands for simple distributed system. It provides all client operations as a client using the distributed system.
*/
package main

import (
	"errors"
	"fa18cs425mp/src/shared/cfg"
	"fmt"
	"os"
)

var (
	ArgsCopy    []string
	TargetNodes []int
	Action      string
)

func parseOverallParas() error {
	if len(os.Args) <= 1 {
		return errors.New("no action specified by 'sds' command")
	}
	ArgsCopy = make([]string, len(os.Args))
	copy(ArgsCopy, os.Args)

	flagSet := ParseArgs(nil, &ArgsCopy)
	if flagSet == nil {
		return errors.New("fail to parse Overall parameters")
	}

	Action = ArgsCopy[0]
	TargetNodes = *flagSet.Lookup("n").Value.(*ArrayFlags)
	if TargetNodes == nil {
		for i := 0; i < len(cfg.Cfg.Addrs); i++ {
			TargetNodes = append(TargetNodes, i)
		}
	}
	return nil
}

func main() {
	if err := parseOverallParas(); err != nil {
		panic(err)
	}
	switch Action {
	case "sdfs":
		fmt.Println("sdfs invoked")
		dsdfs()
	case "grep":
		fmt.Println("grep invoked")
		dgrep()
	case "close":
		fmt.Println("close invoked")
		dclose()
	case "config":
		fmt.Println("config invoked")
		dconfig()
	case "swim":
		fmt.Println("swim membership invoked")
		dswim()
	default:
		fmt.Println("Invalid action ", Action)
	}
	return
}
