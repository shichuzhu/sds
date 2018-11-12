package main

import (
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/lib/sdfs"
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/shared/cfg"
	"fa18cs425mp/src/shared/sdfs2fd"
	"flag"
	"google.golang.org/grpc"
	"log"
	"os"
)

var (
	dataPath   = flag.String("datapath", cfg.Cfg.GrepDir, "The path to files to be grep")
	logDir     = flag.String("log", cfg.Cfg.LogDir, "Directory to store the log")
	closeSigs  chan int
	lg         = utils.LogMessage{}
	grpcServer *grpc.Server
)

func main() {
	flag.Parse()
	closeSigs = make(chan int, 1)
	memlist.SetupLocalMemList()

	SetupLogger()
	SetupGRpc()
	memlist.StartFailureDetector()

	sdfs.InitialSdfs()
	defer close(sdfs2fd.Communicate)

	switch <-closeSigs {
	case 1:
		CleanUp()
		log.Println("Server Closed by Client.")
	case 2:
		log.Fatalln("Fatal: Simulating fatal failure of node")
		os.Exit(0)
	}
}
