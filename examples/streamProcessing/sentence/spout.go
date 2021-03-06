package main

import (
	"bufio"
	"fa18cs425mp/src/lib/stream/shared"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"time"
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
		s.lineNumber++
		instr := s.scanner.Text()
		if len(instr) >= 28 {
			obj := Words{Words: instr}
			bts, _ := proto.Marshal(&obj)
			collector.Emit(bts)
		}
		if s.lineNumber%15 == 0 {
			collector.IssueCheckPoint()
		}
		//fmt.Println("Spout Sending: ", s.scanner.Text())
		time.Sleep(time.Second)
	} else {
		collector.IssueStop()
	}
}
