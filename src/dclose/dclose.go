package main

import (
	co "fa18cs425mp/src/lib/connect"
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func closeConnection(conn *grpc.ClientConn) error {
	defer wg.Done()
	defer conn.Close() // optional
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	text, err := client.CloseServer(ctx, &pb.IntMessage{Mesg: 1})
	if err != nil {
		log.Printf("failure to close server: %s\n", err)
	} else {
		log.Println(*text)
	}
	return nil
}

func main() {
	conn, err := co.Connect()
	if err != nil {
		log.Println("All the server is closed")
		os.Exit(1)
	}

	for i := 0; i < len(conn); i++ {
		wg.Add(1)
		go closeConnection(conn[i])
	}
	wg.Wait()
}
