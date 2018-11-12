package main

import (
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/lib/sdfs"
	"fa18cs425mp/src/pb"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func SetupGRpc() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *memlist.TcpPort))
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
	vmIndex = int32(memlist.MemList.MyNodeId)
	lg.Init(vmIndex, 1, *logDir)
}
