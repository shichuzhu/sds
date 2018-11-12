package main

import (
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/lib/sdfs"
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/shared/sdfs2fd"
	"flag"
	"google.golang.org/grpc"
)

var (
	port       = flag.Int("port", memlist.DefaultTcpPort, "The server port")
	dataPath   = flag.String("dataPath", "data", "The path to files to be grep")
	logDir     = flag.String("log", "data/mp2", "Directory to store the log")
	nodeId     = flag.Int("nodeid", -1, "The nodeid if not randomized")
	sdfsPath   = flag.String("sdfsPath", "data/mp3/", "The path to sdfs files to be stored")
	closeSigs  chan int
	logLevel   int32
	vmIndex    int32
	lg         = utils.LogMessage{}
	grpcServer *grpc.Server
)

func main() {
	closeSigs = make(chan int)

	RegisterFdFlags() // Also call the flag.Parse() inside
	SetupGrpc()
	if *nodeId != -1 {
		memlist.MembershipList.MyNodeId = *nodeId
	}
	SetupLogger()
	SpawnFailureDetector()

	sdfs.SdfsRootPath = *sdfsPath
	sdfs.InitialSdfs()
	defer close(sdfs2fd.Communicate)

	if <-closeSigs == 1 {
		CleanUp()
	}
}
