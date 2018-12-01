package worker

import "fa18cs425mp/src/pb"

// TODO: implement emitter
type Collector struct {
	// state to record connection
	stream pb.StreamProcServices_StreamTuplesServer
}

func NewCollector(cfg *pb.TaskCfg) *Collector {
	return new(Collector)
}

func (s *Collector) Emit(arr []byte) {
	// TODO: error handle
	_ = s.stream.Send(&pb.BytesTuple{
		BytesOneof: &pb.BytesTuple_Tuple{Tuple: arr}})
}
