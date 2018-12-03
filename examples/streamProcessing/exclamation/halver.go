package main

import (
	"fa18cs425mp/src/lib/stream/shared"
	"fa18cs425mp/src/lib/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"strings"
)

type Halver struct {
	states [][]string
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

	var another []string
	log.Println("New tuple reaches sink: ", obj.Words)
	for i, word := range obj.Words {
		if i%2 == 0 {
			another = append(another, word)
		}
	}
	//anotherObj := &Words{Words: another}
	s.states = append(s.states, another)
	//anotherArr, _ := proto.Marshal(anotherObj)
	//abc.Emit(anotherArr)
	return
}

func (s *Halver) CheckPoint() {
	log.Println("sink checkPointing...")
	for _, words := range s.states {
		_, _ = fmt.Fprintln(s.file, strings.Join(words, " "))
	}
	s.states = nil
	//s.file.Close()
	_ = s.file.Sync()
	_ = utils.RunShellString(fmt.Sprintf("sds sdfs put %s %s", s.fn, s.fn))
	//s.file, _ = os.OpenFile(s.fn, os.O_APPEND|os.O_WRONLY, 0600)
}
