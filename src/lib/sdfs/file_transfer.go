package sdfs

import (
	"errors"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"path/filepath"
	"time"
)

func PullKeyFromNode(key, nodeId int) error {
	return nil
}

/*
Transfer file to remote server.

localFilePath: local file path. If empty, deriving from sdfsFilePath using latest version
sdfsFilePath: sdfs file path. If empty, deriving from localFilePath
ip: remote server gRPC address
*/
func FileTransferToNode(ip, localFilePath, sdfsFilePath string) error {
	if sdfsFilePath == "" {
		sdfsFilePath, _ = LfsToSdfs(filepath.Base(localFilePath))
	} else if localFilePath == "" {
		localFilePath = SdfsToLfs(sdfsFilePath, GetFileVersion(sdfsFilePath))
	}
	conn, _ := connect(ip)
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fileClient, err := client.TransferFiles(ctx)
	if err != nil {
		log.Println("Fail to transfer file to: ", ip)
		return err
	}
	file, err := os.Open(localFilePath)
	if err != nil {
		log.Println("Can not find local file path:", localFilePath)
		return err
	}

	message := &pb.FileTransMessage{
		Message: &pb.FileTransMessage_Config{
			Config: &pb.FileCfg{RepNumber: 0, RemoteFilepath: sdfsFilePath}}}

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
		log.Printf("File '%s' Transferred to '%s' as '%s'", localFilePath, ip, sdfsFilePath)
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
