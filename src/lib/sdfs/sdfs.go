package sdfs

import (
	"errors"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const REPLICA = 4

var SdfsRootPath string

func SdfsPut(localFileName, sdfsFilename string) {
	key := HashToKey(sdfsFilename)
	nodeId := FindNodeId(key, 0)
	ip := nodeId.Addr()
	if err := FileTransferToNode(ip, localFileName, sdfsFilename); err != nil {
		log.Println("Initial transfer to master failed")
		return
	}

	conn, _ := connect(ip)
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	retMessage, err := client.PutFile(ctx, &pb.StringMessage{Mesg: sdfsFilename})
	if err != nil {
		log.Println("Failure during in putting file")
		return
	}
	if retMessage.Mesg == 1 {
		log.Println("Successfully put file into replicas")
	}
}
func SdfsGet(sdfsFilename, localFilename string) {
	key := HashToKey(sdfsFilename)
	fileMaster := FindNodeId(key, 0)
	/*
		TODO:Find IP for the particular key
	*/

	var ip string
	ret := pullFile(sdfsFilename, ip, 1)
	if ret == -1 {
		log.Println("Cannot find this file under name: " + sdfsFilename)
		return
	}

	fileVersion := GetFileVersion(sdfsFilename)
	currentLocalName := SdfsToLfs(sdfsFilename, fileVersion)
	err := os.Rename(SdfsRootPath+currentLocalName, SdfsRootPath+localFilename)
	if err != nil {
		log.Println("Error in change file name")
		return
	}

	log.Println("Successfully get file " + sdfsFilename + " From system")
	return

}
func SdfsDelete(sdfsFilename string) {
}
func SdfsLs(fileName string) {
}

func SdfsStore() {
	listOfFile := listAllFile()
	for e := listOfFile.Front(); e != nil; e = e.Next() {
		log.Println(e.Value)
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
