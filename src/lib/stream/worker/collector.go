package worker

import "fa18cs425mp/src/pb"

// TODO: implement emitter
type Collector struct {
	// state to record connection
	stream pb.StreamProcServices_StreamTuplesServer
	err    error
}

func NewCollector(server pb.StreamProcServices_StreamTuplesServer) *Collector {
	return &Collector{stream: server}
}

func (s *Collector) Emit(arr []byte) {
	// TODO: error handle
	err := s.stream.Send(&pb.BytesTuple{
		BytesOneof: &pb.BytesTuple_Tuple{Tuple: arr}})
	if err != nil && s.err != nil {
		s.err = err
	}
}

func (s *Collector) IssueStop() {
	// Send control signal and remove task from taskManager
	_ = s.stream.Send(&pb.BytesTuple{
		BytesOneof: &pb.BytesTuple_ControlSignal{ControlSignal: 1}})
	return
}

func (s *Collector) IssueCheckPoint() {
	_ = s.stream.Send(&pb.BytesTuple{
		BytesOneof: &pb.BytesTuple_ControlSignal{ControlSignal: 0}})
	return
}

func (s *Collector) IsLive() error {
	return s.err
}
