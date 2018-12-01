package services

import (
	"context"
	"fa18cs425mp/src/pb"
)

type StreamProcServer struct{}

// Master
func (s *StreamProcServer) SubmitJob(ctx context.Context, config *pb.TopoConfig) (*pb.TopoConfig, error) {
	return nil, nil
}

// Standby Master
func (s *StreamProcServer) SyncMasterState(ctx context.Context, config *pb.TopoConfig) (*pb.TopoConfig, error) {
	return nil, nil
}

// Worker
func (s *StreamProcServer) SpawnTask(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
	return nil, nil
}

func (s *StreamProcServer) Anchor(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
	// TODO: call task.StreamTuple
	return nil, nil
}

func (s *StreamProcServer) CheckPoint(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
	return nil, nil
}

func (s *StreamProcServer) Terminate(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
	return nil, nil
}

// This function should not return
func (s *StreamProcServer) StreamTuples(cfg *pb.TaskCfg, stream pb.StreamProcServices_StreamTuplesServer) error {
	return nil
}
