package sdfs

import (
	"errors"
	"fa18cs425mp/src/lib/membership"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"path/filepath"
	"time"
)

func IdToIp(nodeId int) string {
	tmp := membership.NextNofId(0, nodeId)
	return tmp.Addr()
}

func PullKeyFromNode(key, nodeId int) error {
	info := &pb.PullFileInfo{FetchType: 1,
		MyID:     int32(membership.MembershipList.MyNodeId),
		FetchKey: int32(key)}
	client, err := GetClientOfNodeId(nodeId)
	if err != nil {
		return err
	}
	_, err = (*client).PullFiles(context.Background(), info)
	if err != nil {
		return err
	}
	return nil
}

func FileTransferToNodeByIp(ip, localFilePath, sdfsFilePath string, igMT bool) error {
	if conn, err := connect(ip); err != nil {
		log.Println("Fail to connect to node ", ip)
		return err
	} else {
		defer conn.Close()
		log.Println("Connect to node ", ip)
		client := pb.NewServerServicesClient(conn)
		return FileTransferToNode(&client, localFilePath, sdfsFilePath, igMT)
	}
}

func GetClientOfNodeId(nodeId int) (*pb.ServerServicesClient, error) {
	ip := IdToIp(nodeId)
	if conn, err := connect(ip); err != nil {
		log.Println("Fail to connect to node ", nodeId)
		return nil, err
	} else {
		log.Println("Connect to node ", nodeId)
		client := pb.NewServerServicesClient(conn)
		return &client, nil
	}
}

/*
Transfer file to remote server.

localFilePath: local file path. If empty, deriving from sdfsFilePath using latest version
sdfsFilePath: sdfs file path. If empty, deriving from localFilePath
ip: remote server gRPC address
*/
func FileTransferToNode(client *pb.ServerServicesClient, localFilePath, sdfsFilePath string, igMT bool) error {
	version := 0
	if sdfsFilePath == "" {
		sdfsFilePath, version = LfsToSdfs(filepath.Base(localFilePath))
	} else if localFilePath == "" {
		version = GetFileVersion(sdfsFilePath)
		localFilePath = SdfsToLfs(sdfsFilePath, version)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fileClient, err := (*client).TransferFiles(ctx)
	if err != nil {
		return err
	}
	file, err := os.Open(localFilePath)
	if err != nil {
		log.Println("Can not find local file path:", localFilePath)
		return err
	}

	message := &pb.FileTransMessage{
		Message: &pb.FileTransMessage_Config{
			Config: &pb.FileCfg{RepNumber: 0, FileVersion: int32(version),
				RemoteFilepath: sdfsFilePath, IgnoreMemtable: igMT}}}

	fileClient.Send(message)

	buf := make([]byte, 1024)

	n, _ := file.Read(buf)
	for n != 0 {
		message = &pb.FileTransMessage{
			Message: &pb.FileTransMessage_Chunk{Chunk: buf[0:n]},
		}
		fileClient.Send(message)
		n, _ = file.Read(buf)
	}

	recv, err := fileClient.CloseAndRecv()
	if err != nil {
		return err
	}
	if recv.GetMesg() == 1 {
		log.Printf("File '%s' Transferred as '%s'", localFilePath, sdfsFilePath)
	}
	return nil
}

func connect(IP string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*3))
	conn, err := grpc.Dial(IP, opts...)
	if err != nil {
		message := fmt.Sprintf("CAN NOT CONNECT TO IP %v", IP)
		log.Println(message)
		return nil, errors.New(message)
	}
	return conn, nil
}
