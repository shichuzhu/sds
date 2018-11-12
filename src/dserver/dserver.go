package main

import (
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/lib/sdfs"
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/shared/cfg"
	"fa18cs425mp/src/shared/sdfs2fd"
	"flag"
	"google.golang.org/grpc"
)

var (
	dataPath   = flag.String("datapath", cfg.Cfg.GrepDir, "The path to files to be grep")
	logDir     = flag.String("log", cfg.Cfg.LogDir, "Directory to store the log")
	closeSigs  chan int
	logLevel   int32
	vmIndex    int32
	lg         = utils.LogMessage{}
	grpcServer *grpc.Server
)

func main() {
	flag.Parse()
	closeSigs = make(chan int, 1)

	SetupGRpc()
	SetupLogger()
	memlist.SpawnFailureDetector()

	sdfs.InitialSdfs()
	defer close(sdfs2fd.Communicate)

	if <-closeSigs == 1 {
		CleanUp()
	}
}
