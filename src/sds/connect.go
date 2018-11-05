package main

import (
	"encoding/json"
	"errors"
	"fa18cs425mp/src/lib/membership"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

type Configuration struct {
	Addrs []struct {
		IP   string
		Port int
	}
}

var configFileName = "cfg.json"
var config Configuration
var opts []grpc.DialOption

func loadConfig() {
	fileContent, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Println("Cannot read the config file")
	}
	if err := json.Unmarshal(fileContent, &config); err != nil {
		log.Println("Fail to parse the JSON config file")
	}
}

func Connect() ([]*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*3))

	var ret []*grpc.ClientConn

	for _, i := range TargetNodes {
		conn, err := helper(config.Addrs[i].IP, config.Addrs[i].Port)
		if err != nil {
			continue
		}
		ret = append(ret, conn)
	}

	if len(ret) == 0 {
		return nil, errors.New("all connection failed")
	}
	return ret, nil
}

func helper(IP string, port int) (*grpc.ClientConn, error) {
	strAddr := IP + ":" + strconv.Itoa(port)
	conn, err := grpc.Dial(strAddr, opts...)
	if err != nil {
		message := fmt.Sprintf("CAN NOT CONNECT TO IP %v", strAddr)
		log.Println(message)
		return nil, errors.New(message)
	}
	return conn, nil
}

func ConnectLocal() (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*3))
	localIp := membership.GetOutboundIP().String()
	samplePort := 10001 // TODO: Change to 10000 when at remote VMs!!!!!
	if len(config.Addrs) == 0 {
	} else {
		samplePort = config.Addrs[0].Port
	}
	return helper(localIp, samplePort)
}
