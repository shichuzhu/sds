package main

import (
	pb "mp/src/protobuf"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	opts       []grpc.DialOption
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

func (s *Dispatcher) dispatch(addr string, port int) {
	defer s.wg.Done()
	strAddr := addr + ":" + strconv.Itoa(port)
	conn, err := grpc.Dial(strAddr, s.opts...)
	if err != nil {
		log.Printf("fail to dial: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewServerServicesClient(conn)
	s.distGrep(client, &pb.Cmd{Cmd: "grep " + strings.Join(os.Args[1:], " ")})
}

func loadConfig() {
	fileContent, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatalln("Cannot read the config file")
	}
	if err := json.Unmarshal(fileContent, &config); err != nil {
		log.Fatalln("Fail to parse the JSON config file")
	}
}

func main() {
	loadConfig()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*3))
	dispatcher := Dispatcher{opts: opts}

	for i := 0; i < len(config.Addrs); i++ {
		dispatcher.wg.Add(1)
		go dispatcher.dispatch(config.Addrs[i].IP, config.Addrs[i].Port)
	}
	dispatcher.wg.Wait()
}
