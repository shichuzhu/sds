run:
	go run mp1/src/server/server.go -port 10000 &
	go run mp1/src/server/server.go -port 10001 &
	go run mp1/src/server/server.go -port 10002 &
	sleep 1
	go run mp1/src/client/client.go
