package main

import (
	pb "../protobuf"
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
)

func distGrep(client pb.GrepLogClient, cmd *pb.Cmd) {
	//log.Printf("Looking for features within %v", cmd)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ReturnMatches(ctx, cmd)
	if err != nil {
		log.Fatalf("%v.ReturnMatches(_) = _, %v", client, err)
	}
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ReturnMathces(_) = _, %v", client, err)
		}
		log.Println(line)
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewGrepLogClient(conn)
	distGrep(client, &pb.Cmd{Cmd: "grep"})
}
