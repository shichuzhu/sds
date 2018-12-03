package main

import (
	"fa18cs425mp/src/lib/stream/shared"
	"fa18cs425mp/src/lib/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
)

type Halver struct {
	states []string
	file   *os.File
	fn     string
}

func NewHalver() shared.SinkABC {
	return &Halver{fn: "stream.out"}
}

func (s *Halver) Init() (err error) {
	s.file, err = os.Create(s.fn)
	if err != nil {
		log.Println("Spout error: ", err)
	}
	return err
}

func (s *Halver) Execute(arr []byte, abc shared.CollectorABC) {
	obj := new(Words)
	_ = proto.Unmarshal(arr, obj)

	inStr := obj.Words
	log.Println("New tuple reaches sink: ", inStr)
	another := inStr + "!!!"
	s.states = append(s.states, another)
	return
}

func (s *Halver) CheckPoint() {
	log.Println("sink checkPointing...")
	for _, words := range s.states {
		_, _ = fmt.Fprintln(s.file, words)
	}
	s.states = nil
	//s.file.Close()
	_ = s.file.Sync()
	_ = utils.RunShellString(fmt.Sprintf("sds sdfs put %s %s", s.fn, s.fn))
	//s.file, _ = os.OpenFile(s.fn, os.O_APPEND|os.O_WRONLY, 0600)
}
