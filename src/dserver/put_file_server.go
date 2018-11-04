package main

import (
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
)

func (s *serviceServer) PutFile(ctx context.Context, putMessage *pb.StringMessage) (*pb.IntMessage, error) {
	sdfsName := putMessage.Mesg
	version := sdfs.GetFileVersion(sdfsName)
	localName := sdfs.SdfsToLfs(sdfsName, version)

	fileKey := sdfs.HashToKey(sdfsName)
	for i := 1; i <= 3; i++ {
		tmp := sdfs.FindNodeId(fileKey, i)
		ip := tmp.Addr()
		//ip := sdfs.FindNodeId(fileKey, i).Addr()
		sdfs.FileTransferToNode(ip, localName)
	}

	return &pb.IntMessage{Mesg: 1}, nil
}
