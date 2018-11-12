package sdfs

import (
	"errors"
	"fa18cs425mp/src/pb"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
)

type SdfsServer struct{}

func (s *SdfsServer) DeleteFile(ctx context.Context, message *pb.StringMessage) (*pb.IntMessage, error) {
	fileName := message.Mesg
	versions := GetFileVersion(fileName)
	ret := 1
	for i := versions; i != 0; i-- {
		localName := SdfsToLfs(fileName, i)
		err := os.Remove(SdfsRootPath + localName)
		if err != nil {
			log.Println("Error in deleting file")
			ret = -1
		}
	}

	DeleteFileFromTable(fileName)
	return &pb.IntMessage{Mesg: int32(ret)}, nil
}

func (s *SdfsServer) TransferFiles(stream pb.SdfsServices_TransferFilesServer) error {
	message, err := stream.Recv()
	if err != nil {
		return err
	}
	config := message.GetConfig()
	fileName := config.RemoteFilepath
	version := int(config.FileVersion)
	if version == 0 {
		version = GetFileVersion(fileName) + 1
	}
	localName := SdfsToLfs(fileName, version)
	file, err := os.Create(SdfsRootPath + localName)
	for {
		message, err = stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		content := message.GetChunk()
		if content != nil {
			file.Write(content)
		}
	}
	if config.IgnoreMemtable {
		log.Println("Dummy transfer, not updating MemTable!!!")
	} else {
		InsertFileVersion(fileName, version)
	}
	ret := pb.IntMessage{Mesg: 1}
	err = stream.SendAndClose(&ret)

	if err != nil {
		log.Println("Error when receiving file ", fileName, " : ", err)
		return err
	}
	log.Println("Received file " + fileName)
	return nil
}

func (s *SdfsServer) PullFiles(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
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

func (s *SdfsServer) PutFile(ctx context.Context, putMessage *pb.StringMessage) (*pb.IntMessage, error) {
	sdfsName := putMessage.Mesg
	version := GetFileVersion(sdfsName)
	localName := SdfsRootPath + SdfsToLfs(sdfsName, version)

	fileKey := HashToKey(sdfsName)
	for i := 1; i <= 3; i++ {
		tmp := FindNodeMember(fileKey, i)
		ip := tmp.Addr()
		if err := FileTransferToNodeByIp(ip, localName, "", false); err != nil {
			return nil, err
		}
	}
	return &pb.IntMessage{Mesg: 1}, nil
}
