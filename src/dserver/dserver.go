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
	logFile  = flag.String("log", "data/mp2/output.log", "Filepath to store the log")
	//configFile = flag.String("configFile", "remotecfg.json", "The json file containing IP/port info of VMs")
	closeSigs  chan int
	logLevel   int32
	vmIndex    int32
	lg         = cl.LogMessage{}
	grpcServer *grpc.Server
)

type serviceServer struct {
	//object states defined here
}

func (s *serviceServer) ServerConfig(ctx context.Context, info *pb.ConfigInfo) (*pb.StringMessage, error) {
	if info.LogLevel >= 0 && info.LogLevel < 4 {
		logLevel = info.LogLevel
	}

	vmIndex = info.VmIndex
	//lg = new(cl.LogMessage)
	lg.Init(vmIndex, 1)

	message := fmt.Sprintf("Config information receive successfully by Server %v", vmIndex)
	return &pb.StringMessage{Mesg: message}, nil
}

func (s *serviceServer) ReturnMatches(theCmd *pb.StringArray, stream pb.ServerServices_ReturnMatchesServer) error {
	cmd := exec.Command("/bin/sh", "-c", strings.Join(theCmd.GetMesgs(), " "))
	cmd.Dir = *dataPath
	log.Printf("From \"%s\" executing: %s", cmd.Dir, strings.Join(cmd.Args, " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if err := stream.Send(&pb.StringMessage{Mesg: scanner.Text()}); err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceServer) CloseServer(_ context.Context, closeMessage *pb.IntMessage) (*pb.StringMessage, error) {
	var message string
	if closeType := closeMessage.GetMesg(); closeType == 0 {
		log.Fatalln("Fatal: Simulating fatal failure of node")
	} else {
		fmt.Println("Server Closed by Client.")
		closeSigs <- 1
		message = fmt.Sprintf("Server %v Successfully closed", vmIndex)

	}
	return &pb.StringMessage{Mesg: message}, nil
}

func SetupGrpc() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer = grpc.NewServer()
	pb.RegisterServerServicesServer(grpcServer, &serviceServer{})
	go grpcServer.Serve(lis)
}

func CleanUp() {
	grpcServer.GracefulStop()
	lg.Close()
}

func main() {
	closeSigs = make(chan int)

	RegisterFdFlags() // Also call the flag.Parse() inside
	SetupLogger()

	SetupGrpc()
	SpawnFailureDetector()
	InitialSdfs()

	if <-closeSigs == 1 {
		CleanUp()
	}
}
