package main

import (
	ms "fa18cs425mp/src/lib/membership"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

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
