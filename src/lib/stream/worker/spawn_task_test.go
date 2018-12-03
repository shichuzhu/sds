package worker

import (
	"fa18cs425mp/src/lib/stream/config"
	"fa18cs425mp/src/lib/stream/shared"
	"fa18cs425mp/src/pb"
	"flag"
	"fmt"
	"testing"
)

func TestSpawnBolt(t *testing.T) {
	flag.Parse()
	config.InitialCrane()
	cfg := &pb.TaskCfg{JobName: "exclamation",
		Bolt: &pb.Bolt{Name: "Spout"}}
	_ = config.SetupLoadFile(cfg.JobName)
	plug := CompilePlugin(cfg)

	col := new(TestCollector)

	spout := SpawnSpoutTask(cfg, plug)
	_ = spout.Init()
	spout2 := SpawnSpoutTask(cfg, plug)
	_ = spout2.Init()

	cfg.Bolt.Name = "ExclaimAdder"
	bolt1 := SpawnBoltTask(cfg, plug)
	_ = bolt1.Init()
	cfg.Bolt.Name = "Halver"
	bolt2 := SpawnSinkTask(cfg, plug)
	_ = bolt2.Init()

	spout.NextTuple(col)
	bolt1.Execute(col.state, col)
	bolt2.Execute(col.state, col)

	spout.NextTuple(col)
	bolt1.Execute(col.state, col)
	bolt2.Execute(col.state, col)

	spout.NextTuple(col)
	bolt1.Execute(col.state, col)
	bolt2.Execute(col.state, col)

	spout2.NextTuple(col)
	bolt1.Execute(col.state, col)
	bolt2.Execute(col.state, col)

	sink, ok := bolt2.(shared.SinkABC)
	if ok {
		sink.CheckPoint()
	}
}

type TestCollector struct {
	state []byte
}

func (s *TestCollector) Emit(arr []byte) {
	if arr != nil {
		s.state = arr
	}
	fmt.Println(s.state)
}

func (s *TestCollector) IssueStop() {
	return
}

func (s *TestCollector) IssueCheckPoint() {
	return
}
