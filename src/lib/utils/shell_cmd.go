package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func RunShellArgs(cmdStrs []string) error {
	cmd := exec.Command(cmdStrs[0], cmdStrs[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func RunShellString(cmd string) error {
	fmt.Println(cmd)
	return RunShellArgs(strings.Split(cmd, " "))
}