package sdfs

import (
	"fa18cs425mp/src/lib/membership"
	pb "fa18cs425mp/src/protobuf"
	"golang.org/x/net/context"
	"log"
)

func pullFile(sdsFileName string, ip string, versions int, config *pb.PullFileInfo) *pb.PullFileInfo {
	conn, _ := connect(ip)
	client := pb.NewServerServicesClient(conn)
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	info := &pb.PullFileInfo{FileName: sdsFileName, NumOfFile: int32(versions),
		MyID: int32(membership.MembershipList.MyNodeId)}
	if config != nil {
		info.IgnoreMemtable = config.IgnoreMemtable
	}
	retMessage, err := client.PullFiles(ctx, info)
	if err != nil {
		log.Println("Failure in pull files")
		return nil
	}
	return retMessage
}
