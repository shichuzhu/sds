package main

import (
	"errors"
	ms "fa18cs425mp/src/lib/membership"
	"fa18cs425mp/src/pb"
	"flag"
	"fmt"
	"golang.org/x/net/context"
)

func SetupLogger() {
	vmIndex = int32(ms.MembershipList.MyNodeId)
	lg.Init(vmIndex, 1, *logDir)
}

func RegisterFdFlags() {
	port := flag.Int("pfd", 11000, "port number to use for the failure detector")
	drop := flag.Float64("d", 0.0, "Simulated packet drop rate.")

	flag.Parse()

	ms.PacketDrop.SetDropRate(float32(*drop))
	ms.MembershipList.MyPort = *port
}

func SpawnFailureDetector() {
	ms.InitInstance()
	ms.StartFailureDetector()
}

func (s *serviceServer) ActMembership(_ context.Context, argsMsgs *pb.StringArray) (*pb.StringMessage, error) {
	args := argsMsgs.GetMesgs()
	var response string
	switch text := args[0]; text {
	case "ls":
		response = ms.FormListInfo()
	case "leave":
		ms.LeaveGroup()
	case "join":
		ms.JoinByIntroducer(args[1])
	default:
		fmt.Println("Invalid input.")
		return nil, errors.New("Invalid input")
	}
	return &pb.StringMessage{Mesg: response}, nil
}
