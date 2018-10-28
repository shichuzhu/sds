/**
Package sds stands for simple distributed system. It provides all client operations as a client using the distributed system.
*/
package main

import (
	"errors"
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
		return errors.New("No action specified by 'sds' command.")
	}
	ArgsCopy = make([]string, len(os.Args))
	copy(ArgsCopy, os.Args)

	flagSet := ParseArgs(nil, &ArgsCopy)
	if flagSet == nil {
		return errors.New("Fail to parse Overall parameters")
	}

	Action = ArgsCopy[0]
	TargetNodes = *flagSet.Lookup("n").Value.(*ArrayFlags)
	loadConfig()
	if TargetNodes == nil {
		for i := 0; i < len(config.Addrs); i++ {
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
