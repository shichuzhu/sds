package Connect_Pack

import (
	"encoding/json"
	"errors"
	pa "fa18cs425mp/src/lib/parseargs"
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
		log.Fatalln("Cannot read the config file")
	}
	if err := json.Unmarshal(fileContent, &config); err != nil {
		log.Fatalln("Fail to parse the JSON config file")
	}
}

func Connect() ([]*grpc.ClientConn, error) {
	loadConfig()
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*3))

	var ret []*grpc.ClientConn

	set := pa.RegisterNodeArgs(nil)
	if !pa.ParseArgs(set) {
		log.Panicln("Fail to parse Node")
	}
	if pa.ServerIds == nil {
		for i := 0; i < len(config.Addrs); i++ {
			pa.ServerIds = append(pa.ServerIds, i)
		}
	}

	for _, i := range pa.ServerIds {
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
