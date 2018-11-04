package sdfs

import (
	"errors"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"time"
)

func FileTransferToNode(ip string, port int, filePath string) {
	conn, _ := connect(ip, port)
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fileClient, _ := client.TransferFiles(ctx)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can not find file path" + filePath)
	}

	message := &pb.FileTransMessage{
		FileTransMessage: &pb.FileTransMessage_Config_{
			Config: &pb.FileTransMessage_Config{
				RemoteFilepath: filePath,
				RepNumber:      0,
			},
		},
	}
	fileClient.Send(message)

	buf := make([]byte, 1024)

	n, _ := file.Read(buf)
	for n != 0 {
		message = &pb.FileTransMessage{
			FileTransMessage: &pb.FileTransMessage_Chunk{Chunk: buf[0:n]},
		}
		fileClient.Send(message)
		n, _ = file.Read(buf)
	}

	recv, err := fileClient.CloseAndRecv()
	if recv.GetMesg() == 1 {
		fmt.Println("File has been successfully transfer to ")
	}

}

func connect(IP string, port int) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*3))
	strAddr := IP + ":" + strconv.Itoa(port)
	conn, err := grpc.Dial(strAddr, opts...)
	if err != nil {
		message := fmt.Sprintf("CAN NOT CONNECT TO IP %v", strAddr)
		log.Println(message)
		return nil, errors.New(message)
	}
	return conn, nil
}
