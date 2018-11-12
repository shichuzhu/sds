package main

import (
	"errors"
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/shared/cfg"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

var opts []grpc.DialOption

func Connect() ([]*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*3))

	var ret []*grpc.ClientConn

	for _, i := range TargetNodes {
		conn, err := helper(cfg.Cfg.Addrs[i].IP, cfg.Cfg.Addrs[i].Port)
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
	localIp := memlist.GetOutboundIP().String()
	return helper(localIp, cfg.Cfg.DefaultTCPPort)
}
