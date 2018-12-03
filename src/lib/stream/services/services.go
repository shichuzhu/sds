package services

import (
	"context"
	"fa18cs425mp/src/lib/stream/master"
	"fa18cs425mp/src/lib/stream/worker"
	"fa18cs425mp/src/pb"
)

type StreamProcServer struct{}

// Master
/*
	TODO: We need to get file config name here
*/
func (s *StreamProcServer) SubmitJob(ctx context.Context, config *pb.TopoConfig) (*pb.TopoConfig, error) {
	jsonName := config.JobName
	err := master.SpawnTaskMaster(jsonName)

	/*
		return message here is original message
	*/
	return config, err
}

// Standby Master
func (s *StreamProcServer) SyncMasterState(ctx context.Context, config *pb.TopoConfig) (*pb.TopoConfig, error) {
	return nil, nil
}

// Worker
func (s *StreamProcServer) SpawnTask(ctx context.Context, cfg *pb.TaskCfg) (*pb.TaskCfg, error) {
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
	return worker.StreamTuple(cfg, stream)
}
