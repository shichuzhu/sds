package main

import (
	"errors"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
)

func (s *serviceServer) SdfsCall(_ context.Context, argsMsgs *pb.StringArray) (*pb.StringArray, error) {
	args := argsMsgs.GetMesgs()
	var response []string
	switch text := args[0]; text {
	case "leave":

	case "join":

	default:
		fmt.Println("Invalid input.")
		return nil, errors.New("Invalid input")
	}
	return &pb.StringArray{Mesgs: response}, nil
}
