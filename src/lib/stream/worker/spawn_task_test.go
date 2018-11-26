package worker

import (
	"fa18cs425mp/src/lib/stream/config"
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/pb"
	"flag"
	"testing"
)

func TestSpawnBolt(t *testing.T) {
	flag.Parse()
	config.InitialCrane()
	cfg := &pb.TaskCfg{JobName: "exclamation"}
	SetupDirectories(cfg)
	_ = utils.RunShellString("cp test/mp4/user_code/exclamation.zip data/mp4/exclamation/src")
	plug := CompilePlugin(cfg)
	spout1 := SpawnTask(cfg, plug)
	spout1.Init()
	spout2 := SpawnTask(cfg, plug)
	spout2.Init()
}
