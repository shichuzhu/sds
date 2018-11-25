package main

import (
	"fa18cs425mp/src/lib/stream"
	"flag"
)

var (
	BoltId         *int
	configFileName *string
	IsSink         *bool
)

// TODO: change to grpc instead of flag. Since arguments may change upon failure
func init() {
	BoltId = flag.Int("boltId", -1, "Bolt ID of the current process")
	IsSink = flag.Bool("sink", false, "Current bolt is a sink")
	configFileName = flag.String("topo", "topo.json", "Path to topology json file")
	flag.Parse()
}

func GetBolt() stream.BoltABC {
	var bolt stream.BoltABC
	switch *BoltId {
	case 1:
		bolt = &ExclaimAdder{}
	case 2:
		bolt = &Halver{}
	default:
		return nil
	}
	return bolt
}

func GetSpout() stream.SpoutABC {
	return &Spout{}
}

func GetSink() stream.SinkABC {
	return &Halver{}
}

func isSpout() bool {
	return *BoltId == 0
}

func isSink() bool {
	return *IsSink
}

func main() {
	if isSpout() {
		spout := GetSpout()
		spout.NextTuple()
	} else if isSink() {
		sink := GetSink()
		sink.Init()
		sink.Execute()
		sink.CheckPoint()
	} else {
		bolt := GetBolt()
		bolt.Init()
		bolt.Execute()
	}
}
