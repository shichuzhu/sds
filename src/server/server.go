package main

import (
	"bufio"
	pb "fa18cs425mp/src/protobuf"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

var (
	port     = flag.Int("port", 10000, "The server port")
	dataPath = flag.String("dataPath", "data", "The path to files to be grep")
)

type grepLogServer struct {
	//object states defined here
}

func (s *grepLogServer) ReturnMatches(theCmd *pb.Cmd, stream pb.ServerServices_ReturnMatchesServer) error {
	//log.Printf("New request received with pattern: %s\n", theCmd)
	cmd := exec.Command("/bin/sh", "-c", theCmd.GetCmd())
	cmd.Dir = *dataPath
	log.Printf("From \"%s\" executing: %s", cmd.Dir, strings.Join(cmd.Args, " "))
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
	pb.RegisterServerServicesServer(grpcServer, &grepLogServer{})
	//grpcServer.Serve(lis)
	go grpcServer.Serve(lis)
	time.Sleep(time.Second * 10)
	grpcServer.GracefulStop()
}
