package main

import (
	"errors"
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
)

func (s *serviceServer) PullFiles(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
	targetID := int(info.MyID)
	targetFile := info.FileName
	targetNum := int(info.NumOfFile)
	igMT := info.IgnoreMemtable

	versions := sdfs.GetFileVersion(targetFile)
	if versions == 0 {
		return nil, errors.New("No file on remote server")
	}
	ip := sdfs.IdToIp(targetID)
	var numToTransfer int
	if targetNum < versions {
		numToTransfer = targetNum
	} else {
		numToTransfer = versions
	}

	for i := 0; i < numToTransfer; i++ {
		localFileName := sdfs.SdfsToLfs(targetFile, versions-i)
		sdfs.FileTransferToNodeByIp(ip, sdfs.SdfsRootPath+localFileName, "", igMT)
	}

	return &pb.PullFileInfo{LatestVersion: int32(versions)}, nil
}
