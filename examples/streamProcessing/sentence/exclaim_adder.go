package main

import (
	"fa18cs425mp/src/lib/stream/shared"
	"github.com/golang/protobuf/proto"
	"strings"
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

	//log.Println("bolt got: ", obj.Words)
	inStr := obj.Words
	newArr, _ := proto.Marshal(&Words{Words: strings.ToUpper(inStr)})
	abc.Emit(newArr)
	return
}
