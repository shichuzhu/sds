package main

import (
	pb "../protobuf"
	"bufio"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/exec"
	"time"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type grepLogServer struct {
	//object states defined here
}

func (s *grepLogServer) ReturnMatches(theCmd *pb.Cmd, stream pb.GrepLog_ReturnMatchesServer) error {
	log.Println("New request received!")
	var dir = "mp1/src/toys/"
	var python3 = "/usr/bin/python3"
	cmd := exec.Command(python3, dir+"gen.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if err := stream.Send(&pb.GrepLine{Line: scanner.Text()}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGrepLogServer(grpcServer, &grepLogServer{})
	go grpcServer.Serve(lis)
	time.Sleep(time.Second * 5)
	grpcServer.GracefulStop()
}
