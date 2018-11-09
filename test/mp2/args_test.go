package mp2_test

import (
	"fmt"
	"os"
)

func main_test_mp2() {
	if len(os.Args) <= 1 {
		fmt.Println("No action specified by 'sds' command.")
		return
	}
	ArgsCopy := make([]string, len(os.Args))
	copy(ArgsCopy, os.Args)
	action := ArgsCopy[1]
	ArgsCopy[1] = ArgsCopy[0]
	ArgsCopy = ArgsCopy[1:]

	switch action {
	case "sdfs":
		fmt.Println("sdfs invoked")
	case "grep":
		fmt.Println("grep invoked")
	case "close":
		fmt.Println("close invoked")
	case "config":
		fmt.Println("config invoked")
	default:
		fmt.Println("Invalid action ", action)
	}
	fmt.Println(ArgsCopy)
	return
}
