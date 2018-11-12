package main

import (
	"bufio"
	"fa18cs425mp/src/lib/membership"
	"fa18cs425mp/src/lib/sdfs"
	"fa18cs425mp/src/lib/utils"
	pb "fa18cs425mp/src/protobuf"
	"fa18cs425mp/src/shared/sdfs2fd"
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
	port       = flag.Int("port", membership.DefaultTcpPort, "The server port")
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

type serviceServer struct {
	//object states defined here
}

func (s *serviceServer) ServerConfig(ctx context.Context, info *pb.ConfigInfo) (*pb.StringMessage, error) {
	if info.LogLevel >= 0 && info.LogLevel < 4 {
		logLevel = info.LogLevel
	}

	if info.VmIndex == -1 {
		vmIndex = int32(membership.MembershipList.MyNodeId)
	} else {
		vmIndex = info.VmIndex
	}

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
		log.Println("Server Closed by Client.")
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
	log.Println("This is the last log before closing the log file.")
	lg.Close()
}

func main() {
	closeSigs = make(chan int)

	RegisterFdFlags() // Also call the flag.Parse() inside
	SetupGrpc()
	if *nodeId != -1 {
		membership.MembershipList.MyNodeId = *nodeId
	}
	SetupLogger()
	SpawnFailureDetector()

	sdfs.SdfsRootPath = *sdfsPath
	InitialSdfs()
	defer close(sdfs2fd.Communicate)

	if <-closeSigs == 1 {
		CleanUp()
	}
}
