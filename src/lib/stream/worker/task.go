package worker

import (
	"context"
	"errors"
	"fa18cs425mp/src/lib/stream/config"
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
	TaskId               int

	Receiver  pb.StreamProcServices_StreamTuplesClient
	Collector shared.CollectorABC
	Executor  shared.BoltABC
	Spout     shared.SpoutABC
	Sink      shared.SinkABC
}

func NewTask(cfg *pb.TaskCfg) *pb.TaskCfg {
	task := new(Task)

	_ = config.SetupLoadFile(cfg.JobName)
	plug := CompilePlugin(cfg)

	switch cfg.Bolt.BoltType {
	case pb.BoltType_SINK:
		task.Sink = SpawnSinkTask(cfg, plug)
	case pb.BoltType_SPOUT:
		task.Spout = SpawnSpoutTask(cfg, plug)
		task.downStreamSetterSync = make(chan bool, 1)
	default:
		task.Executor = SpawnBoltTask(cfg, plug)
		task.downStreamSetterSync = make(chan bool, 1)
	}
	task.TaskId = GetTMgr().InsertTask(task)
	task.Cfg = cfg
	cfg.TaskId = int64(task.TaskId)
	return cfg
}

// This cfg is the cfg of successor, not self!
func StreamTuple(cfg *pb.TaskCfg, server pb.StreamProcServices_StreamTuplesServer) error {
	id := int(cfg.PredTaskId[0]) // bug fixed: anchor need to specify upstream task id
	task := GetTMgr().Task(id)
	//log.Printf("%v\n", *task.Cfg) // Pretty print the cfg for debug only
	switch task.BoltType() {
	case pb.BoltType_SPOUT:
		go task.RegisterDownStream(cfg, server)
		return task.StreamSpoutTuple() // TODO: change different stream function
	case pb.BoltType_SINK:
		log.Fatalln("Sink should not call this")
		return nil
	default:
		go task.RegisterDownStream(cfg, server)
		return task.StreamBoltTuple()
	}
}

// Connect to data source and anchor the input stream there
func Anchor(cfg *pb.TaskCfg) error {
	id := IdFromCfg(cfg)
	task := GetTMgr().Task(id)
	task.Cfg.PredAddrs = cfg.PredAddrs
	task.Cfg.PredTaskId = cfg.PredTaskId
	switch task.BoltType() {
	case pb.BoltType_SPOUT:
		_ = task.Spout.Init()
		return nil
	case pb.BoltType_SINK:
		_ = task.Sink.Init()
		task.Cfg.PredTaskId = cfg.PredTaskId
		err := task.ConnectUpStream()
		if err != nil {
			return err
		} else {
			// Bug fixed: No one calls sink's streaming, need to invoke at anchor
			go task.StreamSinkTuple()
			return nil
		}
	default:
		_ = task.Executor.Init()
		task.Cfg.PredTaskId = cfg.PredTaskId
		return task.ConnectUpStream()
	}
}

func (s *Task) BoltType() pb.BoltType {
	return s.Cfg.Bolt.BoltType
}

func (s *Task) StreamSinkTuple() {
	// TODO: add mechanism to ack the master.
	for {
		// handle receiver (upstream error)
		arr, control, err := s.GetNextTupleBytes()

		switch control {
		case 0: // checkpoint
			s.Sink.CheckPoint()
		case 1: // stop
			_ = GetTMgr().RemoveTask(s)
			return
		}

		if err != nil {
			log.Println("sink receiving error: ", err)
			return
		} else if arr != nil {
			s.Sink.Execute(arr, s.Collector)
		}
	}
}

func (s *Task) StreamSpoutTuple() error {
	<-s.downStreamSetterSync
	for {
		// handle previous sender side error (downstream error)
		convert, ok := s.Collector.(*Collector)
		if ok && convert.IsLive() != nil {
			return convert.IsLive()
		}

		s.Spout.NextTuple(s.Collector)
	}
}

func (s *Task) StreamBoltTuple() error {
	<-s.downStreamSetterSync
	for {
		// handle previous sender side error (downstream error)
		convert, ok := s.Collector.(*Collector)
		if ok && convert.IsLive() != nil {
			return convert.IsLive()
		}
		// handle receiver (upstream error)
		arr, control, err := s.GetNextTupleBytes()

		switch control {
		case 0: // checkpoint
			s.Collector.IssueCheckPoint()
		case 1: // stop
			s.Collector.IssueStop()
			_ = GetTMgr().RemoveTask(s)
			return nil
		}

		if err != nil {
			log.Println("bolt receiving error: ", err)
			return err
		} else if arr != nil {
			s.Executor.Execute(arr, s.Collector)
		}
	}
}

func (s *Task) ConnectUpStream() error {
	// TODO: to support DAG, need to anchor to multiple upstream
	conn, err := grpc.Dial(s.Cfg.PredAddrs[0], grpc.WithInsecure())

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
	// TODO: add mechanism to register multiple downstream, change Collector to a slice, and sync setting up all downstream before ack the 'true'
	s.Collector = NewCollector(server)

	s.downStreamSetterSync <- true
	close(s.downStreamSetterSync)
	s.downStreamSetterSync = nil
}

func (s *Task) GetNextTupleBytes() ([]byte, int, error) {
	data, err := s.Receiver.Recv()
	if err != nil {
		return nil, -1, err
	}
	switch x := data.BytesOneof.(type) {
	case *pb.BytesTuple_Tuple:
		// Normal tuple byte string
		return x.Tuple, -1, nil
	case *pb.BytesTuple_ControlSignal:
		log.Println("control signal received: ", x.ControlSignal)
		return nil, int(x.ControlSignal), nil
	default:
		err := errors.New("unexpected input stream byte")
		return nil, -1, err
	}
}
