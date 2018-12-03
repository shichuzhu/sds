package main

import (
	"context"
	"fa18cs425mp/src/lib/utils"
	"fa18cs425mp/src/pb"
	"fmt"
	"log"
	"path"
)

func crane() {
	// Connect to local gRPC server, always.
	ArgsCopy = ArgsCopy[1:]
	jobFullPath := ArgsCopy[0]
	if len(jobFullPath) > 1 && jobFullPath[len(jobFullPath)-1] == '/' {
		jobFullPath = jobFullPath[:len(jobFullPath)-1]
	}
	if conn, err := ConnectLocal(); err != nil {
		log.Panic(err)
	} else {
		client := pb.NewStreamProcServicesClient(conn)
		ctx := context.Background()

		err := func() error {
			jobPath := path.Dir(jobFullPath)
			jobName := path.Base(jobFullPath)
			_ = utils.RunShellString(fmt.Sprintf("zip -rj %s/%s.zip %s", jobPath, jobName, jobFullPath))
			if err != nil {
				log.Panic(err)
			}
			_ = utils.RunShellString(fmt.Sprintf("sds sdfs %s/%s.zip %s.zip", jobPath, jobName, jobName))
			if err != nil {
				log.Panic(err)
			}
			_, err = client.SubmitJob(ctx, &pb.TopoConfig{JobName: jobName})
			return err
		}()

		if err != nil {
			log.Panic(err)
		} else {
			// Wait for result/ack.
			fmt.Printf("Job %s is completed.", jobFullPath)
		}
	}
}
