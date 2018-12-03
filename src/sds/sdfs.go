package main

import (
	"context"
	"fa18cs425mp/src/pb"
	"fmt"
	"log"
	"os"
	"time"
)

func dsdfs() {
	// Connect to local gRPC server, always.
	tmpCopy := make([]string, len(ArgsCopy))
	copy(tmpCopy, ArgsCopy)
	tmpCopy[0], _ = os.Getwd()

	ArgsCopy = ArgsCopy[1:]
	if conn, err := ConnectLocal(); err != nil {
		log.Panic(err)
	} else {
		client := pb.NewSdfsServicesClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// copy the instructions to the corresponding local gRPC server call.
		if mesg, err := client.SdfsCall(ctx, &pb.StringArray{Mesgs: tmpCopy}); err != nil {
			log.Panic(err)
		} else {
			// Wait for result/ack.
			fmt.Println(mesg.GetMesgs())
		}
	}
}
