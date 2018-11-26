package main

import "fa18cs425mp/src/lib/stream/shared"

// TODO: implement me and chain the bytes in spawn_task_test.go
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
	abc.Emit(arr)
	return
}
