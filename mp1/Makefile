all: build run

build: src/server src/client
	go build -o bin/server src/server/server.go
	go build -o bin/dgrep src/client/client.go

run: bin/server bin/dgrep
	bin/server -port 10000 &
	bin/server -port 10001 &
	bin/server -port 10002 &
	sleep 1
	bin/dgrep -n "'#4'"

clean:
	rm -f bin/*

.PHONY: clean all
