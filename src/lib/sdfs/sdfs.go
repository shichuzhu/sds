package sdfs

import (
	"errors"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func SdfsPut(localFileName, sdfsFilename string) {
	key := HashToKey(sdfsFilename)
	nodeId := FindNodeId(key, 0)
	/*
		TODO: May Node ID to IP and port number
	*/
	var ip string //IP of NODEID
	FileTransferToNode(ip, localFileName)
	conn, _ := connect(ip)
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	retMessage, err := client.PutFile(ctx, &pb.StringMessage{Mesg: sdfsFilename})

	if err != nil {
		fmt.Println("Failure during in putting file")
	}

	if retMessage.Mesg == 1 {
		fmt.Println("Successfully put file into 4 replicas")
	}
}
func SdfsGet(sdfsFilename, localFilename string) {
}
func SdfsDelete(sdfsFilename string) {
}
func SdfsLs(fileName string) {
}

func SdfsStore() {
	listOfFile := listAllFile()
	for e := listOfFile.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func connectPut(IP string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*15))
	conn, err := grpc.Dial(IP, opts...)
	if err != nil {
		message := fmt.Sprintf("CAN NOT CONNECT TO IP %v", IP)
		log.Println(message)
		return nil, errors.New(message)
	}
	return conn, nil
}
