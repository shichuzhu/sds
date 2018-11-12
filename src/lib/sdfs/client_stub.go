package sdfs

import (
	"errors"
	"fa18cs425mp/src/pb"
	"fa18cs425mp/src/shared/sdfs2fd"
	"golang.org/x/net/context"
	"log"
	"os"
	"strconv"
)

func (s *SdfsServer) SdfsCall(_ context.Context, argsMsgs *pb.StringArray) (*pb.StringArray, error) {
	args := argsMsgs.GetMesgs()
	var response []string
	switch text := args[0]; text {
	case "put":
		if len(args) == 3 {
			SdfsPut(args[1], args[2])
		}
	case "get":
		if len(args) == 3 {
			SdfsGet(args[1], args[2])
		}
	case "delete":
		if len(args) == 2 {
			SdfsDelete(args[1])
		}
	case "ls":
		if len(args) == 2 {
			response = SdfsLs(args[1])
		}
	case "store":
		response = SdfsStore()
	case "get-versions":
		if len(args) == 4 {
			numVar, err := strconv.Atoi(args[2])
			if err != nil {
				return nil, errors.New("Please input integer for version number")
			}
			SdfsGetVersions(args[1], numVar, args[3])
		}
	default:
		log.Println("Invalid input.")
		return nil, errors.New("Invalid input")
	}
	return &pb.StringArray{Mesgs: response}, nil
}

func InitialSdfs() {
	SdfsRootPath = *SdfsRootPathp
	sdfs2fd.Fd2Sdfs = make(chan int)
	sdfs2fd.Sdfs2Fd = make(chan int)
	os.RemoveAll(SdfsRootPath)
	os.Mkdir(SdfsRootPath, os.ModePerm)
	MemTableIntial()
	go ReReplicateHandler()
}
