package sdfs

import (
	"errors"
	pb "fa18cs425mp/src/protobuf"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const REPLICA = 4

var SdfsRootPath string

func SdfsPut(localFileName, sdfsFilename string) {
	key := HashToKey(sdfsFilename)
	client, _ := GetClientOfNodeId(FindNodeId(key, 0))
	if err := FileTransferToNode(client, localFileName, sdfsFilename, false); err != nil {
		log.Println("Initial transfer to master failed")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	retMessage, err := (*client).PutFile(ctx, &pb.StringMessage{Mesg: sdfsFilename})
	if err != nil {
		log.Println("Failure during in putting file")
		return
	}
	if retMessage.Mesg == 1 {
		log.Println("Successfully put file into replicas")
	}
}

func SdfsGet(sdfsFilename, localFilename string) {
	key := HashToKey(sdfsFilename)
	fileMaster := FindNodeId(key, 0)
	ip := IdToIp(fileMaster)

	ret := pullFile(sdfsFilename, ip, 1, &pb.PullFileInfo{IgnoreMemtable: true})
	if ret == nil {
		log.Println("Cannot find this file under name: " + sdfsFilename)
		return
	}

	fileVersion := GetFileVersion(sdfsFilename)
	currentLocalName := SdfsToLfs(sdfsFilename, fileVersion)
	err := exec.Command("cp", SdfsRootPath+currentLocalName, localFilename).Run()
	if err != nil {
		log.Println("Error in copy file")
		return
	}

	log.Println("Successfully get file " + sdfsFilename + " From system")
	return
}

func SdfsDelete(sdfsFilename string) {
	fileKey := HashToKey(sdfsFilename)
	for i := 0; i < 4; i++ {
		NodeID := FindNodeId(fileKey, i)
		ret := callDeleteFile(sdfsFilename, NodeID)

		if ret == -1 {
			log.Println("Error in delete file")
		}
	}
}

func SdfsLs(fileName string) []string {
	retStr := make([]string, 0, 4)
	fileKey := HashToKey(fileName)
	for i := 0; i < 4; i++ {
		retStr = append(retStr, strconv.Itoa(FindNodeId(fileKey, i)))
	}
	return retStr
}

func SdfsStore() []string {
	listOfFile := listAllFile()
	retStr := make([]string, 0)
	for e := listOfFile.Front(); e != nil; e = e.Next() {
		//log.Println(e.Value)
		retStr = append(retStr, e.Value.(string))
	}
	return retStr
}

func connectPut(IP string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(time.Second*15))
	conn, err := grpc.Dial(IP, opts...)
	if err != nil {
		message := fmt.Sprintf("CAN NOT CONNECT TO IP %v", IP)
		log.Println(message)
		return nil, errors.New(message)
	}
	return conn, nil
}

func SdfsGetVersions(sdfsFilename string, numVersions int, localfilename string) {
	file, err := os.Create(localfilename)
	defer file.Close()

	if err != nil {
		log.Println("Cannot create local file")
		return
	}
	key := HashToKey(sdfsFilename)
	fileMaster := FindNodeId(key, 0)
	ip := IdToIp(fileMaster)

	retConfig := pullFile(sdfsFilename, ip, numVersions, &pb.PullFileInfo{IgnoreMemtable: true})
	if retConfig == nil {
		return
	}
	currVersions := int(retConfig.LatestVersion)
	endVersion := currVersions - numVersions
	if endVersion < 0 {
		endVersion = 0
	}
	for i := currVersions; i != endVersion; i-- {
		localFileName := SdfsToLfs(sdfsFilename, i)
		localFile, err := os.Open(SdfsRootPath + localFileName)
		if err != nil {
			log.Println("Cannot get local file" + localFileName)
			continue
		}
		index := strconv.Itoa(i)
		file.WriteString("File version: " + index + "\n\n")

		buffer := make([]byte, 1024)
		n, _ := localFile.Read(buffer)

		for n != 0 {
			file.Write(buffer[0:n])
			n, _ = localFile.Read(buffer)
		}
		file.WriteString("\n\n")
	}

	return
}
