package worker_test

import (
	"context"
	"fa18cs425mp/src/lib/stream/config"
	"fa18cs425mp/src/lib/stream/services"
	"fa18cs425mp/src/pb"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
	"time"
)

func TestBoltWithNetwork(t *testing.T) {
	flag.Parse()
	config.InitialCrane()
	cfg0 := &pb.TaskCfg{JobName: "exclamation",
		Bolt: &pb.Bolt{Name: "Spout",
			BoltType: pb.BoltType_SPOUT},
		PredAddrs: []string{"localhost:12345"}}
	cfg1 := &pb.TaskCfg{JobName: "exclamation",
		Bolt: &pb.Bolt{Name: "ExclaimAdder",
			BoltType: pb.BoltType_BOLT},
		PredAddrs: []string{"localhost:12345"}}
	cfg2 := &pb.TaskCfg{JobName: "exclamation",
		Bolt: &pb.Bolt{Name: "Halver",
			BoltType: pb.BoltType_SINK},
		PredAddrs: []string{"localhost:12345"}}
	//_ = utils.RunShellString("cp test/mp4/user_code/exclamation.zip data/mp4/exclamation/src")

	// setup gRpc service
	grpcServer := grpc.NewServer()
	pb.RegisterStreamProcServicesServer(grpcServer, &services.StreamProcServer{})
	lis, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go grpcServer.Serve(lis)

	// setup gRpc client
	conn, _ := grpc.Dial("localhost:12345", grpc.WithInsecure())
	client := pb.NewStreamProcServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cfg0, err = client.SpawnTask(ctx, cfg0)
	if err != nil {
		t.Fatal(err)
	}
	println(cfg0.TaskId)
	cfg1, err = client.SpawnTask(ctx, cfg1)
	if err != nil {
		t.Fatal(err)
	}
	println(cfg1.TaskId)
	cfg2, err = client.SpawnTask(ctx, cfg2)
	if err != nil {
		t.Fatal(err)
	}
	println(cfg2.TaskId)

	// Stream!
	cfg0, err = client.Anchor(ctx, cfg0)
	if err != nil {
		t.Fatal(err)
	}
	println(cfg0.TaskId)

	cfg1.PredTaskId = []int64{0}
	cfg1, err = client.Anchor(ctx, cfg1)
	if err != nil {
		t.Fatal(err)
	}
	println(cfg1.TaskId)

	cfg2.PredTaskId = []int64{1}
	cfg2, err = client.Anchor(ctx, cfg2)
	if err != nil {
		t.Fatal(err)
	}
	println(cfg2.TaskId)

	time.Sleep(1 * time.Second)
	grpcServer.Stop()
}
