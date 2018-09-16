package main

import (
	co "fa18cs425mp/src/lib/connect"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Configuration struct {
	Addrs []struct {
		IP   string
		Port int
	}
}

var configFileName = "config.json"
var config Configuration

type Dispatcher struct {
	writerLock sync.Mutex
	wg         sync.WaitGroup
}

func (s *Dispatcher) distGrep(client pb.ServerServicesClient, cmd *pb.Cmd) {
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

func (s *Dispatcher) dispatch(conn *grpc.ClientConn) {
	defer s.wg.Done()
	defer conn.Close()
	client := pb.NewServerServicesClient(conn)
	s.distGrep(client, &pb.Cmd{Cmd: "grep " + strings.Join(os.Args[1:], " ")})
}

func main() {
	dispatcher := Dispatcher{}
	connLists, err := co.Connect()
	if err != nil {
		log.Println("No Server can be connected")
		os.Exit(1)
	}

	for i := 0; i < len(connLists); i++ {
		dispatcher.wg.Add(1)
		go dispatcher.dispatch(connLists[i])
	}
	dispatcher.wg.Wait()
}
