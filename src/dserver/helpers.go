package main

import (
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/lib/sdfs"
	"fa18cs425mp/src/pb"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func SetupGrpc() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer = grpc.NewServer()
	pb.RegisterServerServicesServer(grpcServer, &serviceServer{})
	pb.RegisterSdfsServicesServer(grpcServer, &sdfs.SdfsServer{})
	go grpcServer.Serve(lis)
}

func CleanUp() {
	grpcServer.GracefulStop()
	log.Println("This is the last log before closing the log file.")
	lg.Close()
}

func SetupLogger() {
	vmIndex = int32(memlist.MembershipList.MyNodeId)
	lg.Init(vmIndex, 1, *logDir)
}

func RegisterFdFlags() {
	port := flag.Int("pfd", 11000, "port number to use for the failure detector")
	drop := flag.Float64("d", 0.0, "Simulated packet drop rate.")

	flag.Parse()

	memlist.PacketDrop.SetDropRate(float32(*drop))
	memlist.MembershipList.MyPort = *port
}

func SpawnFailureDetector() {
	memlist.InitInstance()
	memlist.StartFailureDetector()
}
