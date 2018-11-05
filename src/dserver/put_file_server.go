package main

import (
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
)

func (s *serviceServer) PutFile(ctx context.Context, putMessage *pb.StringMessage) (*pb.IntMessage, error) {
	sdfsName := putMessage.Mesg
	version := sdfs.GetFileVersion(sdfsName)
	localName := sdfs.SdfsRootPath + sdfs.SdfsToLfs(sdfsName, version)

	fileKey := sdfs.HashToKey(sdfsName)
	for i := 1; i <= 3; i++ {
		tmp := sdfs.FindNodeMember(fileKey, i)
		ip := tmp.Addr()
		if err := sdfs.FileTransferToNodeByIp(ip, localName, "", false); err != nil {
			return nil, err
		}
	}
	return &pb.IntMessage{Mesg: 1}, nil
}
