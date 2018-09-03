# FA18CS425MP

## GRPC

### Installation

Install grpc with
`go get google.golang.org/grpc`, and protoc-gen-go with `go get -u github.com/golang/protobuf/protoc-gen-go`

Compile the `.proto` file using

`protoc -I protobuf/ protobuf/grep_log.proto --go_out=plugins=grpc:protobuf`

## How to run

Run the server on VM:

`go run mp1/src/server/server.go`

Then run the the client:

`go run mp1/src/server/client.go`

## Tricks

1. To return the line number, need to add `-n`.
2. To return the file name, need to add another null file (e.g. /dev/null)

## TODO:

- Add json file to support input IP address.
- Write code to autosend IP & port info from server to the client.