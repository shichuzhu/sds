package main

import (
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
)

func (s *serviceServer) PullFiles(ctx context.Context, info *pb.PullFileInfo) (*pb.IntMessage, error) {
	targetID := int(info.MyID)
	targetFile := info.FileName
	targetNum := int(info.NumOfFile)

	versions := sdfs.GetFileVersion(targetFile)
	if versions == 0 {
		return &pb.IntMessage{Mesg: -1}, nil
	}
	/*
		TODO: transfer id to ip (in tcp)
	*/
	ip := sdfs.IdToIp(targetID)
	var numToTransfer int
	if targetNum < versions {
		numToTransfer = targetNum
	} else {
		numToTransfer = versions
	}

	for i := 0; i < numToTransfer; i++ {
		localFileNmae := sdfs.SdfsToLfs(targetFile, versions-i)
		sdfs.FileTransferToNodeByIp(ip, localFileNmae, "")
	}

	return &pb.IntMessage{Mesg: int32(numToTransfer)}, nil
}
