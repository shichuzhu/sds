package worker

import (
	"context"
	"errors"
	"fa18cs425mp/src/lib/stream/shared"
	"fa18cs425mp/src/pb"
	"google.golang.org/grpc"
	"log"
)

type TaskId struct {
	Addr   string
	TaskId int
}

type Task struct {
	CollectorSet chan bool
	Pred         *TaskId
	Self         *TaskId
	Next         *TaskId
	Cfg          *pb.TaskCfg

	Receiver  pb.StreamProcServices_StreamTuplesClient
	Collector shared.CollectorABC
	Executor  shared.BoltABC
}

func NewTask() *Task {
	return &Task{CollectorSet: make(chan bool, 1)}
}

func (s *Task) StreamTuple() error {
	<-s.CollectorSet
	for {
		arr, err := s.GetNextTupleBytes()
		if err != nil {
			return err
		}
		s.Executor.Execute(arr, s.Collector)
	}
}

func (s *Task) Run(cfg *pb.TaskCfg) error {
	s.ConnectUpStream()
	s.ConnectDownStream(cfg)
	return nil
}

func (s *Task) ConnectUpStream() {
	conn, err := grpc.Dial(s.Cfg.PredAddrs[0], nil)
	if err != nil {
		log.Printf("fail to dial: %v", err)
		return
	}
	client := pb.NewStreamProcServicesClient(conn)
	ctx := context.Background()
	s.Receiver, err = client.StreamTuples(ctx, s.Cfg)
	return
}

// To be called by rpc from downstream
func (s *Task) ConnectDownStream(cfg *pb.TaskCfg) {
	s.CollectorSet <- true
	s.Cfg = cfg // TODO: only update what is needed
	s.Collector = NewCollector(s.Cfg)
	close(s.CollectorSet)
	s.CollectorSet = nil
}

func (s *Task) GetNextTupleBytes() ([]byte, error) {
	data, err := s.Receiver.Recv()
	if err != nil {
		return nil, err
	}
	switch x := data.BytesOneof.(type) {
	case *pb.BytesTuple_Tuple:
		println(x.Tuple)
		return nil, nil
	case *pb.BytesTuple_ControlSignal:
		println(x.ControlSignal)
		return nil, nil
	default:
		err := errors.New("unexpected input stream byte")
		return nil, err
	}
}
