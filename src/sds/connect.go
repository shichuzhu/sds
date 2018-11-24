package main

import (
	"errors"
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/shared/cfg"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

var (
	opts []grpc.DialOption
)

func init() {
	opts = append(opts, grpc.WithInsecure())
}

func Connect() ([]*grpc.ClientConn, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	strAddr := IP + ":" + strconv.Itoa(port)
	conn, err := grpc.DialContext(ctx, strAddr, opts...)
	if err != nil {
		message := fmt.Sprintf("CAN NOT CONNECT TO IP %v", strAddr)
		log.Println(message)
		return nil, errors.New(message)
	}
	return conn, nil
}

func ConnectLocal() (*grpc.ClientConn, error) {
	localIp := memlist.GetOutboundIP().String()
	return helper(localIp, cfg.Cfg.DefaultTCPPort)
}
