package main

import (
	"fa18cs425mp/src/lib/stream"
	"flag"
)

var (
	BoltId         *int
	configFileName *string
)

func init() {
	BoltId = flag.Int("boltId", -1, "Bolt ID of the current process")
	configFileName = flag.String("topo", "topo.json", "Path to topology json file")
	flag.Parse()
}

func GetBolt() stream.BoltABC {
	var bolt stream.BoltABC
	switch *BoltId {
	case 0:
		bolt = &Spout{}
	case 1:
		bolt = &ExclaimAdder{}
	case 2:
		bolt = &Halver{}
	default:
		panic("no bolt specified") // TODO: debug purpose only
		return nil
	}
	return bolt
}

func main() {
	bolt := GetBolt()
	bolt.Init()
	bolt.NextTuple()
}
