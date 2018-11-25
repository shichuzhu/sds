package sdfs

import (
	"context"
	"errors"
	"fa18cs425mp/src/pb"
	"io"
	"log"
	"os"
)

type Server struct{}

func (s *Server) DeleteFile(ctx context.Context, message *pb.StringMessage) (*pb.IntMessage, error) {
	fileName := message.Mesg
	versions := GetFileVersion(fileName)
	ret := 1
	for i := versions; i != 0; i-- {
		localName := SdfsnameToLfs(fileName, i)
		err := os.Remove(RootPath + localName)
		if err != nil {
			log.Println("Error in deleting file")
			ret = -1
		}
	}

	DeleteFileFromTable(fileName)
	return &pb.IntMessage{Mesg: int32(ret)}, nil
}

func (s *Server) TransferFiles(stream pb.SdfsServices_TransferFilesServer) error {
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
	localName := SdfsnameToLfs(fileName, version)
	file, err := os.Create(RootPath + localName)
	for {
		message, err = stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		content := message.GetChunk()
		if content != nil {
			_, _ = file.Write(content)
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

func (s *Server) PullFiles(ctx context.Context, info *pb.PullFileInfo) (*pb.PullFileInfo, error) {
	switch info.FetchType {
	case 0:
		return PutSingleSdfsFile(info)
	case 1:
		return FetchEntireKey(ctx, info)
	case 2:
		return QueryExistence(info)
	default:
		return nil, errors.New("unknown FetchType")
	}
}

func (s *Server) PutFile(ctx context.Context, putMessage *pb.StringMessage) (*pb.IntMessage, error) {
	sdfsName := putMessage.Mesg
	version := GetFileVersion(sdfsName)
	localName := RootPath + SdfsnameToLfs(sdfsName, version)

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
