package main

import (
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"io"
	"os"
)

func (s *serviceServer) TransferFiles(stream pb.ServerServices_TransferFilesServer) error {
	message, err := stream.Recv()
	if err != nil {
		return err
	}
	fileName := message.GetConfig().RemoteFilepath

	/*Here we need to get version nmmber to create new file)
	TODO: Change the protocol buffer naming system to adapt the change
	*/
	version := sdfs.GetFileVersion(fileName)
	localName := sdfs.SdfsToLfs(fileName, version+1)
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
	sdfs.InsertFileVersion(fileName, version+1)
	ret := pb.IntMessage{Mesg: 1}
	err = stream.SendAndClose(&ret)

	if err != nil {
		fmt.Println("This client has successfully receive file " + fileName)
	}

	return nil
}
