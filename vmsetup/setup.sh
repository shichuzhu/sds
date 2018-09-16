#!/bin/sh
# Run using
# ssh szhu28@fa18-cs425-g44-01.cs.illinois.edu "bash -s" < setup.sh
# For bulk run, use 
# Sadly, the input will be directed to only the first instance.
# Therefore, need to use subprocess.Popen, with shell=True, and quote the '< setup.sh'
# ./deploy For 2 Each "bash -s" '< setup.sh'

mkdir go
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> .bashrc
go get google.golang.org/grpc
go get github.com/golang/protobuf/protoc-gen-go
mkdir mp