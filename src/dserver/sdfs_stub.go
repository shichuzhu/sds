package main

import (
	"errors"
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
)

func (s *serviceServer) SdfsCall(_ context.Context, argsMsgs *pb.StringArray) (*pb.StringArray, error) {
	args := argsMsgs.GetMesgs()
	var response []string
	switch text := args[0]; text {
	case "put":
		if len(args) == 3 {
			sdfs.SdfsPut(args[1], args[2])
		}
	case "get":
		if len(args) == 3 {
			sdfs.SdfsGet(args[1], args[2])
		}
	case "delete":
		if len(args) == 2 {
			sdfs.SdfsDelete(args[1])
		}
	case "ls":
		if len(args) == 2 {
			sdfs.SdfsLs(args[1])
		}
	case "store":
		sdfs.SdfsStore()
	default:
		fmt.Println("Invalid input.")
		return nil, errors.New("Invalid input")
	}
	return &pb.StringArray{Mesgs: response}, nil
}

func InitialSdfs() {
	sdfs.MemTableIntial()
}
