package mp3

import (
	"errors"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"time"
)

func testTransferFile(ip string) {
	FileTransferToNode(ip, "Hello.txt")
}

func FileTransferToNode(ip string, filePath string) error {
	conn, _ := connect(ip)
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fileClient, err := client.TransferFiles(ctx)
	if err != nil {
		fmt.Println("Fail to transfer file to client")
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can not find file path" + filePath)
		return err
	}

	message := &pb.FileTransMessage{
		Message: &pb.FileTransMessage_Config{
			Config: &pb.FileCfg{RepNumber: 0, RemoteFilepath: filePath}}}

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
		fmt.Println("File has been successfully transfer to ")
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

func TestFileTransfer(t *testing.T) {
	testTransferFile("127.0.0.1:10001")
}
