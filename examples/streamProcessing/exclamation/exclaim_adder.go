package main

import (
	"fa18cs425mp/src/lib/stream/shared"
	"github.com/golang/protobuf/proto"
)

type ExclaimAdder struct {
	// states
}

func NewExclaimAdder() shared.BoltABC {
	return &ExclaimAdder{}
}

func (s *ExclaimAdder) Init() error {
	return nil
}

func (s *ExclaimAdder) Execute(arr []byte, abc shared.CollectorABC) {
	obj := new(Words)
	_ = proto.Unmarshal(arr, obj)

	another := []string{}
	for _, word := range obj.Words {
		another = append(another, word+"!!!")
	}
	anotherObj := &Words{Words: another}
	anotherArr, _ := proto.Marshal(anotherObj)
	abc.Emit(anotherArr)
	return
}
