package worker

import (
	"fa18cs425mp/src/lib/stream/config"
	"fa18cs425mp/src/lib/utils"
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
	SetupDirectories(cfg)
	_ = utils.RunShellString("zip -rj data/mp4/exclamation/src/exclamation.zip examples/streamProcessing/exclamation")
	//_ = utils.RunShellString("cp test/mp4/user_code/exclamation.zip data/mp4/exclamation/src")
	plug := CompilePlugin(cfg)
	col := &Collector{}

	spout1 := SpawnSpoutTask(cfg, plug)
	_ = spout1.Init()
	spout2 := SpawnSpoutTask(cfg, plug)
	_ = spout2.Init()

	cfg.Bolt.Name = "ExclaimAdder"
	bolt1 := SpawnBoltTask(cfg, plug)

	bolt1.Execute(nil, col)
	spout1.NextTuple(col)
	bolt1.Execute(nil, col)
	spout1.NextTuple(col)
	bolt1.Execute(nil, col)
	spout2.NextTuple(col)
	bolt1.Execute(nil, col)
}

type Collector struct {
	state []byte
}

func (s *Collector) Emit(arr []byte) {
	if arr != nil {
		s.state = arr
	}
	fmt.Println(s.state)
}
