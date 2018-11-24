package main

import (
	"bufio"
	"errors"
	"fa18cs425mp/src/lib/memlist"
	"fa18cs425mp/src/pb"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os/exec"
	"strings"
)

type serviceServer struct {
	//object states defined here
}

func (s *serviceServer) ServerConfig(ctx context.Context, info *pb.ConfigInfo) (*pb.StringMessage, error) {
	return &pb.StringMessage{Mesg: "Log level set"}, nil
}

func (s *serviceServer) ReturnMatches(theCmd *pb.StringArray, stream pb.ServerServices_ReturnMatchesServer) error {
	cmd := exec.Command("/bin/sh", "-c", strings.Join(theCmd.GetMesgs(), " "))
	cmd.Dir = *dataPath
	log.Printf("From \"%s\" executing: %s", cmd.Dir, strings.Join(cmd.Args, " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if err := stream.Send(&pb.StringMessage{Mesg: scanner.Text()}); err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceServer) CloseServer(_ context.Context, closeMessage *pb.IntMessage) (*pb.StringMessage, error) {
	var message string
	closeSigs <- int(closeMessage.GetMesg())
	return &pb.StringMessage{Mesg: message}, nil
}

func (s *serviceServer) ActMembership(_ context.Context, argsMsgs *pb.StringArray) (*pb.StringMessage, error) {
	args := argsMsgs.GetMesgs()
	var response string
	switch text := args[0]; text {
	case "ls":
		response = memlist.FormListInfo()
	case "leave":
		memlist.LeaveGroup()
	case "join":
		memlist.JoinByIntroducer(args[1])
	default:
		fmt.Println("Invalid input.")
		return nil, errors.New("invalid input")
	}
	return &pb.StringMessage{Mesg: response}, nil
}
