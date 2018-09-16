#!/bin/sh
# Run using
# ssh $CS425NETID@fa18-cs425-g$CS425GROUPID-01.cs.illinois.edu "bash -s" < setup.sh
# For bulk run, use 
# Sadly, the input will be directed to only the first instance.
# Therefore, need to use subprocess.Popen, with shell=True, and quote the '< setup.sh'
# ./deploy For 2 Each "bash -s" '< setup.sh'

cd
mkdir -p data/mp1
mkdir logs
mkdir -p go/src/fa18cs425mp
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> .bashrc
go get google.golang.org/grpc
go get github.com/golang/protobuf/protoc-gen-go
