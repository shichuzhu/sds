package main

import (
	"errors"
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
	"log"
	"time"
)

func (s *serviceServer) PullFiles(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
	switch info.FetchType {
	case 0:
		return PutSingleSdfsFile(ctx, info)
	case 1:
		return FetchEntireKey(ctx, info)
	case 2:
		return QueryExistence(ctx, info)
	default:
		return nil, errors.New("Unknown FetchType")
	}
}

func FetchEntireKey(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
	targetID := int(info.MyID)
	targetKey := int(info.FetchKey)

	targetNum := 5
	client, err := sdfs.GetClientOfNodeId(targetID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	fList := sdfs.ListAllFile()
	for e := fList.Front(); e != nil; e = e.Next() {
		if targetFile := e.Value.(string); sdfs.HashToKey(targetFile) == targetKey {
			var numToTransfer int
			versions := sdfs.GetFileVersion(targetFile)
			if targetNum < versions {
				numToTransfer = targetNum
			} else {
				numToTransfer = versions
			}

			for i := 0; i < numToTransfer; i++ {
				localFileName := sdfs.SdfsToLfs(targetFile, versions-i)
				sdfs.FileTransferToNode(client, sdfs.SdfsRootPath+localFileName, "", false)
			}
		}
	}
	log.Printf("### Completed pushing entire key %d to Node %d.", targetKey, targetID)
	return nil, nil
}

// Used by sdfs ls
func QueryExistence(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
	targetFile := info.FileName
	if version := sdfs.GetFileVersion(targetFile); version != 0 {
		log.Println("Queried file exists: ", targetFile)
		return &pb.PullFileInfo{FileExist: true, LatestVersion: int32(version)}, nil
	} else {
		log.Println("Queried file does not exist: ", targetFile)
		return &pb.PullFileInfo{FileExist: false}, nil
	}
}

func PutSingleSdfsFile(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
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
