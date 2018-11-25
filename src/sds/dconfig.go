package main

import (
	"context"
	"fa18cs425mp/src/pb"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"sync"
	"time"
)

func configConnection(conn *grpc.ClientConn, wg *sync.WaitGroup) error {
	defer wg.Done()
	defer conn.Close()
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	text, err := client.ServerConfig(ctx, &pb.ConfigInfo{LogLevel: 1, VmIndex: -1})
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
	for i := 0; i < len(conn); i++ {
		wg.Add(1)
		_ = configConnection(conn[i], &wg)
	}
	wg.Wait()
}

func dswim() {
	conn, err := Connect()
	if err != nil {
		fmt.Println("All the server is closed")
		os.Exit(0)
	}

	var wg sync.WaitGroup

	ArgsCopy = ArgsCopy[1:]
	text := ArgsCopy[0]
	if text == "join" || text == "ls" || text == "leave" {
		for i := 0; i < len(conn); i++ {
			wg.Add(1)
			go actMembership(conn[i], ArgsCopy[:], &wg)
		}
	}
	wg.Wait()
}
