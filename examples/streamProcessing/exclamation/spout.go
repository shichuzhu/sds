package main

import (
	"bufio"
	"fa18cs425mp/src/lib/stream/shared"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"strings"
)

type Spout struct {
	// states
	lineNumber int
	fileName   string
	scanner    *bufio.Scanner
}

func NewSpout() shared.SpoutABC {
	return &Spout{}
}

func (s *Spout) Init() error {
	log.Println("Plugin invoked successfully!")
	s.fileName = "test/mp4/data/input.txt"
	file, err := os.Open(s.fileName)
	if err != nil {
		log.Println(err)
		return err
	}
	s.scanner = bufio.NewScanner(file)
	return nil
}

func (s *Spout) NextTuple(collector shared.CollectorABC) {
	if s.scanner.Scan() {
		obj := Words{Words: strings.Split(s.scanner.Text(), " ")}
		bts, _ := proto.Marshal(&obj)
		fmt.Println("Spout Sending: ", s.scanner.Text())
		collector.Emit(bts)
	} else {
		collector.IssueStop()
	}
}
