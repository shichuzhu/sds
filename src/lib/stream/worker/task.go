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
	downStreamSetterSync chan bool
	Cfg                  *pb.TaskCfg

	Receiver  pb.StreamProcServices_StreamTuplesClient
	Collector shared.CollectorABC
	Executor  shared.BoltABC
}

func NewTask() *Task {
	return &Task{downStreamSetterSync: make(chan bool, 1)}
}

func StreamTuple(cfg *pb.TaskCfg, server pb.StreamProcServices_StreamTuplesServer) error {
	// TODO: decide if this is bolt/spout/sink; checkout the task; run the routine of task
	// RegisterDownStream
	id := IdFromCfg(cfg)
	task := GetTMgr().Task(id)
	go task.RegisterDownStream(cfg, server)
	return task.StreamTuple()
}

// Connect to data source and anchor the input stream there
func Anchor(cfg *pb.TaskCfg) error {
	id := IdFromCfg(cfg)
	task := GetTMgr().Task(id)
	return task.ConnectUpStream()
}

func (s *Task) StreamTuple() error {
	<-s.downStreamSetterSync
	for {
		arr, err := s.GetNextTupleBytes()
		if err != nil {
			return err
		}
		s.Executor.Execute(arr, s.Collector)
	}
}

func (s *Task) ConnectUpStream() error {
	conn, err := grpc.Dial(s.Cfg.PredAddrs[0], nil)
	if err != nil {
		log.Printf("fail to dial: %v", err)
		return err
	}
	client := pb.NewStreamProcServicesClient(conn)
	ctx := context.Background()
	s.Receiver, err = client.StreamTuples(ctx, s.Cfg)
	if err != nil {
		log.Println("fail to call grpc StreamTuples: ", err)
	}
	return err
}

// To be called by rpc from downstream
func (s *Task) RegisterDownStream(cfg *pb.TaskCfg, server pb.StreamProcServices_StreamTuplesServer) {
	s.downStreamSetterSync <- true
	s.Collector = NewCollector(server)
	close(s.downStreamSetterSync)
	s.downStreamSetterSync = nil
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
