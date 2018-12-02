package master

import (
	"context"
	"errors"
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/pb"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

var BoltNodeMap map[int]int
var BoltTaskMap map[int]int

/*
	In this Spawn task function, I will firstly send spawn task message for each bolt
	It will traverse all the node in node list, and then traverse back
	TODO: Node that this method DIDN'T implement to avoid standby master
*/
func SpawnTaskMaster(configFileName string) {
	BoltNodeMap = make(map[int]int)
	BoltTaskMap = make(map[int]int)
	fileConfig := ReadConfig(configFileName)
	_, _, nodeIDList := memlist.GetListElement()
	masterID := memlist.MemList.MyNodeId

	nodeListIndex := 0
	listLength := len(nodeIDList)
	Bolts := fileConfig.Bolts
	for i := 0; i < len(Bolts); i++ {
		if nodeIDList[nodeListIndex%listLength] == masterID {
			nodeListIndex += 1
		}
		nodeID := nodeIDList[nodeListIndex%listLength]
		err := sendSpawnMessage(nodeID, &Bolts[i], listLength)
		if err != nil {
			log.Println("Have error in sending spawn message ", Bolts[i].ID)
		}
		nodeListIndex += 1
	}

}

/*
	This function is used to transfer the spawn message to NodeID
	NodeID is the input index of node
	Return error if transfer message is not successful
	I will map the Bolt id to task id and node id
	Also it will send message and wait for response
	TODO:In the filed specification field, missing taskID, cp_id, listen_addr
*/
func sendSpawnMessage(nodeID int, bolt *Bolt, boltsLength int) error {
	client, err := GetClientofNodeId(nodeID)
	if err != nil {
		return err
	}
	var spawnMess pb.TaskCfg
	if bolt.ID == 0 {
		spawnMess = pb.TaskCfg{
			Bolt: &pb.Bolt{
				BoltType: pb.BoltType_SPOUT,
				BoltId:   int64(bolt.ID),
				Name:     bolt.Name,
				Preds:    []int64{},
			},
		}
	} else if bolt.ID < boltsLength-1 {
		spawnMess = pb.TaskCfg{
			Bolt: &pb.Bolt{
				BoltType: pb.BoltType_BOLT,
				BoltId:   int64(bolt.ID),
				Name:     bolt.Name,
				Preds:    []int64{},
			},
		}
	} else {
		spawnMess = pb.TaskCfg{
			Bolt: &pb.Bolt{
				BoltType: pb.BoltType_SINK,
				BoltId:   int64(bolt.ID),
				Name:     bolt.Name,
				Preds:    []int64{},
			},
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	taskResponse, _ := (*client).SpawnTask(ctx, &spawnMess)

	/*
		Here we will specific the two map, task map and node
	*/

	BoltNodeMap[bolt.ID] = nodeID //initialize nodeID for the bolt
	BoltTaskMap[bolt.ID] = int(taskResponse.TaskId)

	return nil

}

/*
	This function will send anchor message for each bolt
	The anchor message will specify the NODEID and Taskid, and split by ","
*/

func sendAnchorMessage(bolt *Bolt) error {
	boltID := bolt.ID
	predList := bolt.Pred
	predAddr := make([]string, len(predList))
	for i := 0; i < len(predList); i++ {
		nodeID := BoltNodeMap[boltID]
		taskID := BoltTaskMap[boltID]
		tmp := fmt.Sprintf("%d,%d", nodeID, taskID)
		predAddr[i] = tmp
	}

	client, err := GetClientofNodeId(BoltNodeMap[boltID])
	if err != nil {
		log.Println("Error in initialize the client. In send Anchor function")
		return err
	}

	spawnMess := pb.TaskCfg{
		Bolt: &pb.Bolt{
			BoltId: int64(boltID),
		},
		PredAddrs: predAddr,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	retMess, _ := (*client).Anchor(ctx, &spawnMess)

	if retMess.Success == true {
		return nil
	}
	/*
		TODO: How to deal with not successful Anchor ?
	*/

	return nil

}

func GetClientofNodeId(nodeID int) (*pb.StreamProcServicesClient, error) {
	nodeIP := memlist.NextNofId(0, nodeID).Addr()
	if conn, err := connect(nodeIP); err != nil {
		log.Println("Failure in connection Node ", nodeID)
		return nil, err
	} else {
		client := pb.NewStreamProcServicesClient(conn)
		return &client, nil
	}
}

func connect(IP string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, IP, grpc.WithInsecure())
	if err != nil {
		message := fmt.Sprintf("CAN NOT CONNECT TO IP %v", IP)
		log.Println(message)
		return nil, errors.New(message)
	}
	return conn, nil
}
