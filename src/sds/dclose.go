package main

import (
	"context"
	"fa18cs425mp/src/pb"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
	"time"
)

func closeConnection(conn *grpc.ClientConn, wg *sync.WaitGroup) error {
	defer wg.Done()
	defer conn.Close() // optional
	client := pb.NewServerServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	text, err := client.CloseServer(ctx, &pb.IntMessage{Mesg: 0})
	if err != nil {
		log.Printf("Error on close RPC: %s\n", err)
	} else {
		log.Println(*text)
	}
	return nil
}

func dclose() {
	conn, err := Connect()
	if err != nil {
		log.Println("All the server is closed")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for i := 0; i < len(conn); i++ {
		wg.Add(1)
		go closeConnection(conn[i], &wg)
	}
	wg.Wait()
}
