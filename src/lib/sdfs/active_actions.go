package sdfs

import (
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
	"log"
	"time"
)

func callDeleteFile(sdfsFileNmae string, nodeID int) int {
	client, _ := GetClientOfNodeId(FindNodeId(nodeID, 0))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	retMessage, err := (*client).DeleteFile(ctx, &pb.StringMessage{sdfsFileNmae})
	if err != nil {
		log.Println("Error in calling delete file")
		return -1
	}

	return int(retMessage.Mesg)

}
