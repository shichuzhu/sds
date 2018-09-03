package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

var dir = "mp1/src/toys/"
var python3 = "/usr/bin/python3"

func main() {
	cmd := exec.Command(python3, dir+"gen.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	fmt.Println("Stream created!")
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		//time.Sleep(time.Second)
	}
}
