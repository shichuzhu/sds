# FA18CS425MP

## GRPC

### Installation

Install grpc with
`go get google.golang.org/grpc`, and protoc-gen-go with `go get -u github.com/golang/protobuf/protoc-gen-go`

Compile the `.proto` file using

`protoc -I protobuf/ protobuf/grep_log.proto --go_out=plugins=grpc:protobuf`

## How to run

1. Modify the `mp1/config.json` file to add the appropriate IP addresses and port numbers of the VMs.
2. Modify the `mp1/Makefile` to spawn the corresponding VMs.
3. `make`

## Tricks

1. To return the line number, need to add `-n`.
2. To return the file name, need to add another null file (e.g. /dev/null)
3. The wildcard "*" will not work in exec.Command. Have to do the filepath.Glob manually, or invoke a shell `/bin/sh -c ...`
4. Invoking sh -c results in "#" signs unquoted.
5. Go install will make the binary filename based on package name (which is the name of the directory containing "main"), not the name of file containing "main" function.
6. To uninstall binaries, use `go clean -i mp/...`

## TODO:

- [Done] Add json file to support input IP address.
- [Done] Allow user-input grep pattern. (Currently hard-coded)
- Write command that will cleanup the running service instead of manually cancel each server.
- Write code to auto-send IP & port info from server to the client.
- [Done] Add build functionality instead of running `go run` all the time.
