package main

import (
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"sync"
	"time"
)

func configConnection(conn *grpc.ClientConn, index int32, wg *sync.WaitGroup) error {
	defer wg.Done()
	defer conn.Close()
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	text, err := client.ServerConfig(ctx, &pb.ConfigInfo{LogLevel: 1, VmIndex: index})
	if err != nil {
		fmt.Println("Failure in config server at ")
	}

	fmt.Println(text)
	return nil
}

func actMembership(conn *grpc.ClientConn, args []string, wg *sync.WaitGroup) error {
	defer wg.Done()
	defer conn.Close()
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	text, err := client.ActMembership(ctx, &pb.StringArray{Mesgs: args})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(text.GetMesg())
	}
	return nil
}

func dconfig() {
	conn, err := Connect()
	if err != nil {
		fmt.Println("All the server is closed")
		os.Exit(0)
	}

	var wg sync.WaitGroup
	if len(ArgsCopy) == 0 {
		for i := 0; i < len(conn); i++ {
			wg.Add(1)
			configConnection(conn[i], int32(i), &wg)
		}
	} else {
		text := ArgsCopy[0]
		if text == "join" || text == "ls" || text == "leave" {
			for i := 0; i < len(conn); i++ {
				wg.Add(1)
				actMembership(conn[i], ArgsCopy[:], &wg)
			}
		}
	}
}
