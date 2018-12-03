package worker

import (
	"errors"
	"fa18cs425mp/src/pb"
	"fmt"
	"log"
)

type Collector struct {
	// state to record connection
	stream pb.StreamProcServices_StreamTuplesServer
	err    error
	cpFlag bool
}

func NewCollector(server pb.StreamProcServices_StreamTuplesServer) *Collector {
	return &Collector{stream: server}
}

func (s *Collector) Emit(arr []byte) {
	err := s.stream.Send(&pb.BytesTuple{
		BytesOneof: &pb.BytesTuple_Tuple{Tuple: arr}})
	if err != nil {
		s.err = err
		log.Println(err)
		fmt.Println("error while EMITTING", err)
	}
	s.cpFlag = false
}

func (s *Collector) IssueStop() {
	if !s.cpFlag {
		s.IssueCheckPoint()
	}
	// Send control signal and remove task from taskManager
	s.err = s.stream.Send(&pb.BytesTuple{
		BytesOneof: &pb.BytesTuple_ControlSignal{ControlSignal: 1}})
	if s.err == nil {
		s.err = errors.New("stream stopped by user")
		log.Println(s.err)
	}
	return
}

func (s *Collector) IssueCheckPoint() {
	s.err = s.stream.Send(&pb.BytesTuple{
		BytesOneof: &pb.BytesTuple_ControlSignal{ControlSignal: 0}})
	s.cpFlag = true
	return
}

func (s *Collector) IsLive() error {
	return s.err
}
