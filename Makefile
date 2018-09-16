all: build run

deploy:
	vmsetup/deploy Copy src @go/src/fa18cs425mp/

build:
	vmsetup/deploy For Each 'go install fa18cs425mp/...'

localsetup:
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I src/protobuf/ src/protobuf/server_services.proto --go_out=plugins=grpc:src/protobuf

buildlocal: src/dserver src/dgrep
	go install fa18cs425mp/...

runlocal: buildlocal
	server -port 10000 -dataPath "data/mp1" &
	server -port 10001 -dataPath "data/mp1" &
	server -port 10002 -dataPath "data/mp1" &
	sleep 1
	dgrep '-n "#4" * /dev/null'

run: build
	vmsetup/deploy Spawn Each '-port 10000 -dataPath "data/mp1"'
	sleep 2
	dgrep '-n 515922 * /dev/null'

clean:
	rm -f bin/*

.PHONY: clean all
