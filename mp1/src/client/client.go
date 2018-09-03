package main

import (
	pb "../protobuf"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	serverAddr = flag.String("addr", "127.0.0.1", "The server ip address")
	serverPort = flag.Int("port", 10000, "The server port number")
)

type Dispatcher struct {
	writerLock sync.Mutex
	wg         sync.WaitGroup
	opts       []grpc.DialOption
}

func (s *Dispatcher) distGrep(client pb.GrepLogClient, cmd *pb.Cmd) {
	// The maximum time a client will be waiting for the server to respond.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stream, err := client.ReturnMatches(ctx, cmd)
	if err != nil {
		log.Printf("%v.ReturnMatches(_) = _, %v", client, err)
		return
	}
	for {
		grepLine, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v.ReturnMathces(_) = _, %v", client, err)
			return
		}
		s.writerLock.Lock()
		fmt.Println(grepLine.GetLine())
		s.writerLock.Unlock()
	}
}

func (s *Dispatcher) dispatch(addr string, port int) {
	defer s.wg.Done()
	strAddr := addr + ":" + strconv.Itoa(port)
	conn, err := grpc.Dial(strAddr, s.opts...)
	if err != nil {
		log.Printf("fail to dial: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewGrepLogClient(conn)
	s.distGrep(client, &pb.Cmd{Cmd: "grep"})
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*3))
	dispatcher := Dispatcher{opts: opts}

	for i := 0; i < 2; i++ {
		dispatcher.wg.Add(1)
		go dispatcher.dispatch(*serverAddr, *serverPort+i)
	}
	dispatcher.wg.Wait()
}
