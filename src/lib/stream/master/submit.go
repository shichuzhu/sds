package master

import (
	"context"
	"errors"
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/lib/stream/config"
	"fa18cs425mp/src/pb"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

var BoltNodeMap map[int]int
var BoltTaskMap map[int]int

func SubmitJob(topo *pb.TopoConfig) (ret *pb.TopoConfig, err error) {
	ret = topo
	err = config.SetupLoadFile(topo.JobName)
	if err != nil {
		log.Println("job submission failed: ", err)
		return
	}
	return topo, SpawnTaskMaster(topo)
}

/*
	In this Spawn task function, I will firstly send spawn task message for each bolt
	It will traverse all the node in node list, and then traverse back
	TODO: Node that this method DIDN'T implement to avoid standby master
*/
func SpawnTaskMaster(topo *pb.TopoConfig) (err error) {
	dirPath := config.DirPath(topo.JobName)
	configFileName := fmt.Sprintf("%ssrc/topo.json", dirPath)
	BoltNodeMap = make(map[int]int)
	BoltTaskMap = make(map[int]int)
	fileConfig := ReadConfig(configFileName)
	_, _, nodeIDList := memlist.GetListElement()
	masterID := memlist.MemList.MyNodeId

	nodeListIndex := 0
	listLength := len(nodeIDList)
	Bolts := fileConfig.Bolts
	cfgs := make([]*pb.TaskCfg, len(Bolts))
	for i := 0; i < len(Bolts); i++ {
		if nodeIDList[nodeListIndex%listLength] == masterID {
			nodeListIndex += 1
		}
		nodeID := nodeIDList[nodeListIndex%listLength]
		log.Printf("assigning bolt %d to Node %d\n", Bolts[i].ID, nodeID)
		cfgs[i], err = sendSpawnMessage(nodeID, &Bolts[i], len(Bolts), topo.JobName)
		if err != nil {
			log.Println("error spawning task", Bolts[i].ID, err)
			return err
		}
		nodeListIndex += 1
	}

	for i := 0; i < len(Bolts); i++ {
		err := sendAnchorMessage(&Bolts[i], cfgs[i])
		if err != nil {
			log.Println("Error in sending anchor message" + err.Error())
		}
	}

	return nil
}

/*
	This function is used to transfer the spawn message to NodeID
	NodeID is the input index of node
	Return error if transfer message is not successful
	I will map the Bolt id to task id and node id
	Also it will send message and wait for response
*/
func sendSpawnMessage(nodeID int, bolt *Bolt, boltsLength int, jobName string) (*pb.TaskCfg, error) {
	client, err := GetClientOfNodeId(nodeID)
	if err != nil {
		return nil, err
	}
	var spawnMess pb.TaskCfg
	//log.Println("Comparison!!!!!!", bolt.ID, boltsLength)
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
	spawnMess.JobName = jobName
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	taskResponse, err := (*client).SpawnTask(ctx, &spawnMess)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	/*
		Here we will specific the two map, task map and node
	*/

	BoltNodeMap[bolt.ID] = nodeID //initialize nodeID for the bolt
	BoltTaskMap[bolt.ID] = int(taskResponse.TaskId)

	return taskResponse, nil
}

/*
	This function will send anchor message for each bolt
	The anchor message will specify the NODEID and Taskid, and split by ","
*/

func sendAnchorMessage(bolt *Bolt, cfg *pb.TaskCfg) error {
	boltID := bolt.ID
	predList := bolt.Pred
	predAddr := make([]string, len(predList))
	taskList := make([]int64, len(predList))
	//nodeList := make([]int, len(predList))
	for i := 0; i < len(predList); i++ {
		predID := int(predList[i])
		nodeID := BoltNodeMap[predID]
		taskID := BoltTaskMap[predID]
		taskList[i] = int64(taskID)
		predAddr[i] = memlist.NextNofId(0, nodeID).Addr()
	}

	client, err := GetClientOfNodeId(BoltNodeMap[boltID])
	if err != nil {
		log.Println("Error in initialize the client. In send Anchor function", err)
		return err
	}

	cfg.PredTaskId = taskList
	cfg.PredAddrs = predAddr
	cfg.Bolt.BoltId = int64(boltID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = (*client).Anchor(ctx, cfg)
	if err != nil {
		return err
	}
	/*
		TODO: How to deal with not successful Anchor ?
	*/

	return nil

}

func GetClientOfNodeId(nodeID int) (*pb.StreamProcServicesClient, error) {
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
