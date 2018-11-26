package main

import (
	"bufio"
	"fa18cs425mp/src/lib/stream/shared"
	"fmt"
	"log"
	"os"
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
		text := s.scanner.Text()
		fmt.Println("Sending: ", text)
		collector.Emit([]byte(text))
	}
}
