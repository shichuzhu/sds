package main

import (
	"fa18cs425mp/src/lib/stream/shared"
	"log"
)

type Spout struct {
	// states
}

func NewSpout() shared.SpoutABC {
	return &Spout{}
}

func (s *Spout) Init() {
	log.Println("wow")
	return
}

func (s *Spout) NextTuple() {
	return
}
