package main

import (
	"bufio"
	cl "fa18cs425mp/src/lib/loggenerator"
	pb "fa18cs425mp/src/protobuf"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/exec"
	"strings"
)

var (
	port     = flag.Int("port", 10000, "The server port")
	dataPath = flag.String("dataPath", "data", "The path to files to be grep")
	//configFile = flag.String("configFile", "remotecfg.json", "The json file containing IP/port info of VMs")
	closeSigs chan int
	logLevel  int32
	vmIndex   int32
	lg        = cl.LogMessage{}
)

type serviceServer struct {
	//object states defined here
}

func (s *serviceServer) ServerConfig(ctx context.Context, info *pb.ConfigInfo) (*pb.Info, error) {
	if info.LogLevel >= 0 && info.LogLevel < 4 {
		logLevel = info.LogLevel
	}

	vmIndex = info.VMIndex
	//lg = new(cl.LogMessage)
	lg.Init(vmIndex, 1)

	message := fmt.Sprintf("Config information receive successfully by Server %v", vmIndex)
	return &pb.Info{Info: message}, nil
}

func (s *serviceServer) ReturnMatches(theCmd *pb.Cmd, stream pb.ServerServices_ReturnMatchesServer) error {
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

func (s *serviceServer) CloseServer(_ context.Context, closeMessage *pb.CloseMessage) (*pb.Info, error) {
	var message string
	if closeType := closeMessage.GetCloseType(); closeType == 0 {
		log.Fatalln("Fatal: Simulating fatal failure of node")
	} else {
		fmt.Println("Server Closed by Client.")
		closeSigs <- 1
		message = fmt.Sprintf("Server %v Successfully closed", vmIndex)

	}
	return &pb.Info{Info: message}, nil
}

func main() {
	flag.Parse()
	closeSigs = make(chan int)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterServerServicesServer(grpcServer, &serviceServer{})
	go grpcServer.Serve(lis)
	if <-closeSigs == 1 {
		grpcServer.GracefulStop()
		lg.Close()
	}
}
