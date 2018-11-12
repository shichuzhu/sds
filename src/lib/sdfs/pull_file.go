package sdfs

import (
	"errors"
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/pb"
	"golang.org/x/net/context"
	"log"
	"time"
)

func FetchEntireKey(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
	targetID := int(info.MyID)
	targetKey := int(info.FetchKey)

	targetNum := 5
	client, err := GetClientOfNodeId(targetID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	fList := ListAllFile()
	for e := fList.Front(); e != nil; e = e.Next() {
		if targetFile := e.Value.(string); HashToKey(targetFile) == targetKey {
			var numToTransfer int
			versions := GetFileVersion(targetFile)
			if targetNum < versions {
				numToTransfer = targetNum
			} else {
				numToTransfer = versions
			}

			for i := 0; i < numToTransfer; i++ {
				localFileName := SdfsToLfs(targetFile, versions-i)
				FileTransferToNode(client, SdfsRootPath+localFileName, "", false)
			}
		}
	}
	log.Printf("### Completed pushing entire key %d to Node %d.", targetKey, targetID)
	return nil, nil
}

// Used by sdfs ls
func QueryExistence(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
	targetFile := info.FileName
	if version := GetFileVersion(targetFile); version != 0 {
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

	versions := GetFileVersion(targetFile)
	if versions == 0 {
		return nil, errors.New("No file on remote server")
	}
	ip := IdToIp(targetID)
	var numToTransfer int
	if targetNum < versions {
		numToTransfer = targetNum
	} else {
		numToTransfer = versions
	}

	for i := 0; i < numToTransfer; i++ {
		localFileName := SdfsToLfs(targetFile, versions-i)
		FileTransferToNodeByIp(ip, SdfsRootPath+localFileName, "", igMT)
	}
	return &pb.PullFileInfo{LatestVersion: int32(versions)}, nil
}

func pullFile(sdsFileName string, ip string, versions int, config *pb.PullFileInfo) *pb.PullFileInfo {
	conn, _ := connect(ip)
	client := pb.NewSdfsServicesClient(conn)
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	info := &pb.PullFileInfo{FileName: sdsFileName, NumOfFile: int32(versions),
		MyID: int32(memlist.MembershipList.MyNodeId)}
	if config != nil {
		info.IgnoreMemtable = config.IgnoreMemtable
	}
	retMessage, err := client.PullFiles(ctx, info)
	if err != nil {
		log.Println("Failure in pull files")
		return nil
	}
	return retMessage
}
