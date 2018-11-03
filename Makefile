all: runlocal

#all: deploy build run

deploy:
	vmsetup/deploy Copy src @go/src/fa18cs425mp/

clean:
	vmsetup/deploy For Each 'go clean -i fa18cs425mp/...'

build:
	vmsetup/deploy For Each 'go install fa18cs425mp/...'

run:
	vmsetup/deploy Spawn Each '-port 10000 -dataPath "data/mp2"'
	sleep 2
# The first blank before -c to let flag render it as non-flag argument
	sds grep -n 1,2,3 -c ailure '*'
	sds grep -c 123456 '../mp1/*'

localsetup:
#	go get -u google.golang.org/grpc
#	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I src/protobuf/ src/protobuf/server_services.proto --go_out=plugins=grpc:src/protobuf

buildlocal:
	go install fa18cs425mp/...

runlocal: buildlocal
	dserver -port 10000 -pfd 11000 -dataPath "data/mp1" &
	dserver -port 10001 -pfd 11001 -dataPath "data/mp1" &
	dserver -port 10002 -pfd 11002 -dataPath "data/mp1" &
	#sleep 1
	#sds grep -c 123456 '*'

test:
	test/mp1/runtest

.PHONY: clean all test
