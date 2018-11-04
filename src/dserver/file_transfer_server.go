package main

import (
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"io"
	"os"
)

func (s *serviceServer) TransferFiles(stream pb.ServerServices_TransferFilesServer) error {
	message, err := stream.Recv()
	fileName := message.GetConfig().RemoteFilepath

	/*Here we need to get version nmmber to create new file)
	TODO: Version number create
	*/
	file, err := os.Create(fileName)
	for {
		message, err = stream.Recv()
		if err == io.EOF {
			break
		}
		content := message.GetChunk()
		if content != nil {
			file.Write(content)
		}
	}

	ret := pb.IntMessage{Mesg: 1}
	err = stream.SendAndClose(&ret)

	if err != nil {
		fmt.Println("This client has successfully receive file " + fileName)
	}

	return nil
}
