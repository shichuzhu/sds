/**
Package sds stands for simple distributed system. It provides all client operations as a client using the distributed system.
*/
package main

import (
	"fmt"
	"os"
)

var ArgsCopy []string

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No action specified by 'sds' command.")
		return
	}
	ArgsCopy = make([]string, len(os.Args))
	copy(ArgsCopy, os.Args)
	action := ArgsCopy[1]
	ArgsCopy = ArgsCopy[1:]

	switch action {
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
	default:
		fmt.Println("Invalid action ", action)
	}
	return
}
