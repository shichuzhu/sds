package main

import (
	"fa18cs425mp/src/lib/stream/shared"
	"github.com/golang/protobuf/proto"
	"strings"
)

type Halver struct {
	states [][]string
}

func NewHalver() shared.SinkABC {
	return &Halver{}
}

func (s *Halver) Init() error {
	return nil
}

func (s *Halver) Execute(arr []byte, abc shared.CollectorABC) {
	obj := new(Words)
	_ = proto.Unmarshal(arr, obj)

	var another []string
	//log.Println("sink got: ", obj.Words)
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
	println("sink checkPointing:")
	for _, words := range s.states {
		println(strings.Join(words, " "))
	}
	s.states = nil
}
