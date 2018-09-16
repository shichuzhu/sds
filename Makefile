all: build run

deploy:
	vmsetup/deploy Copy src @go/src/fa18cs425mp/

localsetup:
	go get google.golang.org/grpc
	go get github.com/golang/protobuf/protoc-gen-go
	protoc -I src/protobuf/ src/protobuf/server_services.proto --go_out=plugins=grpc:src/protobuf

build: src/server src/dgrep
	go install fa18cs425mp/...

run: build
	server -port 10000 &
	server -port 10001 &
	server -port 10002 &
	sleep 1
	dgrep '-n "#4" * /dev/null'

clean:
	rm -f bin/*

.PHONY: clean all
