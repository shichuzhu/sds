package main

import (
	co "fa18cs425mp/src/lib/connect"
	ag "fa18cs425mp/src/lib/parseargs"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func configConnection(conn *grpc.ClientConn, index int32) error {
	defer wg.Done()
	defer conn.Close()
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	text, err := client.ServerConfig(ctx, &pb.ConfigInfo{LogLevel: 1, VMIndex: index})
	if err != nil {
		fmt.Println("Failure in config server at ")
	}

	fmt.Println(text)
	return nil
}

func actMembership(conn *grpc.ClientConn, args []string) error {
	defer wg.Done()
	defer conn.Close()
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	text, err := client.ActMembership(ctx, &pb.StringArray{Msgs: args})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(text.GetMesg())
	}
	return nil
}

func main() {
	conn, err := co.Connect()
	if err != nil {
		fmt.Println("All the server is closed")
		os.Exit(0)
	}
	if len(ag.ArgsCopy) == 0 {
		for i := 0; i < len(conn); i++ {
			wg.Add(1)
			configConnection(conn[i], int32(i))
		}
	} else {
		text := ag.ArgsCopy[0]
		if text == "join" || text == "ls" || text == "leave" {
			for i := 0; i < len(conn); i++ {
				wg.Add(1)
				actMembership(conn[i], ag.ArgsCopy[:])
			}
		}
	}
}
