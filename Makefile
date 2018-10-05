all: deploy build run

deploy:
	vmsetup/deploy Copy src @go/src/fa18cs425mp/

clean:
	vmsetup/deploy For Each 'go clean -i fa18cs425mp/...'

build:
	vmsetup/deploy For Each 'go install fa18cs425mp/...'

run: build
	vmsetup/deploy Spawn Each '-port 10000 -dataPath "data/mp1"'
	sleep 2
# The first blank before -c to let flag render it as non-flag argument
	dgrep -n 1,2,3 ' -c 515922 * /dev/null'

localsetup:
#	go get -u google.golang.org/grpc
#	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I src/protobuf/ src/protobuf/server_services.proto --go_out=plugins=grpc:src/protobuf

buildlocal: src/dserver src/dgrep
	go install fa18cs425mp/...

runlocal: buildlocal
	server -port 10000 -dataPath "data/mp1" &
	server -port 10001 -dataPath "data/mp1" &
	server -port 10002 -dataPath "data/mp1" &
	sleep 1
	dgrep -n 1,2,3  ' -c 123456 * /dev/null'

test:
	test/mp1/runtest

.PHONY: clean all test
