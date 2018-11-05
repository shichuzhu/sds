package main

import (
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"io"
	"log"
	"os"
)

func (s *serviceServer) TransferFiles(stream pb.ServerServices_TransferFilesServer) error {
	message, err := stream.Recv()
	if err != nil {
		return err
	}
	config := message.GetConfig()
	fileName := config.RemoteFilepath
	version := int(config.FileVersion)
	if version == 0 {
		version = sdfs.GetFileVersion(fileName) + 1
	}
	localName := sdfs.SdfsToLfs(fileName, version)
	file, err := os.Create(sdfs.SdfsRootPath + localName)
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
		sdfs.InsertFileVersion(fileName, version)
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
