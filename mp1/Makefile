all: build run

deploy:
	go get google.golang.org/grpc
	go get github.com/golang/protobuf/protoc-gen-go
	protoc -I src/protobuf/ src/protobuf/server_services.proto --go_out=plugins=grpc:src/protobuf

build: src/server src/client
	mkdir -p bin
	go build -o bin/server src/server/server.go
	go build -o bin/dgrep src/client/dgrep.go

run: bin/server bin/dgrep
	bin/server -port 10000 &
	bin/server -port 10001 &
	bin/server -port 10002 &
	sleep 1
	bin/dgrep -n "'#4'"

clean:
	rm -f bin/*

.PHONY: clean all
