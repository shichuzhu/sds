package services

import (
	"context"
	"fa18cs425mp/src/lib/stream/master"
	"fa18cs425mp/src/lib/stream/worker"
	"fa18cs425mp/src/pb"
	"fa18cs425mp/src/shared/sdfs2fd"
	"fmt"
	"log"
	"time"
)

type StreamProcServer struct{}

// Master
func (s *StreamProcServer) SubmitJob(ctx context.Context, config *pb.TopoConfig) (*pb.TopoConfig, error) {
	//for sig, ok <- sdfs2fd.Fd2Crane;
	func() {
		for {
			select {
			case _ = <-sdfs2fd.Fd2Crane:
			default:
				return
			}
		}
	}()

	config, _ = master.SubmitJob(config)
	go func() {
		<-sdfs2fd.Fd2Crane
		time.Sleep(3 * time.Second)
		log.Println("Failure detected, restarting job!")
		_, _ = master.SubmitJob(config)
	}()
	return config, nil
}

// Standby Master
func (s *StreamProcServer) SyncMasterState(ctx context.Context, config *pb.TopoConfig) (*pb.TopoConfig, error) {
	return nil, nil
}

// Worker
func (s *StreamProcServer) SpawnTask(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
	log.Printf("to spawn bolt %d as type %v", cfg.Bolt.BoltId, cfg.Bolt.BoltType)
	cfg = worker.NewTask(cfg)
	return cfg, nil
}

func (s *StreamProcServer) Anchor(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
	err := worker.Anchor(cfg)
	return cfg, err
}

func (s *StreamProcServer) CheckPoint(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
	return nil, nil
}

func (s *StreamProcServer) Terminate(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
	return nil, nil
}

// This function should not return until streaming stops
func (s *StreamProcServer) StreamTuples(cfg *pb.TaskCfg, stream pb.StreamProcServices_StreamTuplesServer) error {
	err := worker.StreamTuple(cfg, stream)
	fmt.Println("stream terminated")
	return err
}
