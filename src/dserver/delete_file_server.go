package main

import (
	"fa18cs425mp/src/lib/sdfs"
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
	"log"
	"os"
)

func (s *serviceServer) DeleteFile(ctx context.Context, message *pb.StringMessage) (*pb.IntMessage, error) {
	fileName := message.Mesg
	versions := sdfs.GetFileVersion(fileName)
	ret := 1
	for i := versions; i != 0; i-- {
		localName := sdfs.SdfsToLfs(fileName, i)
		err := os.Remove(sdfs.SdfsRootPath + localName)
		if err != nil {
			log.Println("Error in deleting file")
			ret = -1
		}
	}

	sdfs.DeleteFileFromTable(fileName)
	return &pb.IntMessage{Mesg: int32(ret)}, nil
}
