package worker

import (
	"fa18cs425mp/src/lib/stream/shared"
	"log"
	"plugin"
	"testing"
)

func TestSpawnBolt(t *testing.T) {
	spout1 := SpawnTask()
	spout1.Init()
	spout2 := SpawnTask()
	spout2.Init()
}

func SpawnTask() shared.SpoutABC {
	filepath := "/home/shichu/usr/gopath/src/fa18cs425mp/exclaimation.so"
	plug, err := plugin.Open(filepath)
	if err != nil {
		log.Println(err)
		return nil
	}
	sym, err := plug.Lookup("NewSpout")
	if err != nil {
		log.Println(err)
		return nil
	}
	var spout shared.SpoutABC
	spout = sym.(func() shared.SpoutABC)()
	return spout
}
