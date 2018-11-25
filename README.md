# FA18CS425MP

## MP4

### Design

Entities:
1. Master per system
2. Stand-by master per system
3. Supervisor per VM
4. Task
5. Client. Submit job. Local code, SDFS data sources (both input data and file databases).

* Master and supervisor detect each other's failure.
* Supervisor detect Tasks' failure.
i.e. each tasks owns a Unix domain socket on disk and supervisor 'connect/search' for its existence.
  * If any socket exists, supervisor needs to report its failure upon starting.
  * TODO: Current implementation doesn't support node rejoining.
  Therefore, always assume all tasks fail simultaneously
  * TODO: In the sds close method, need to close all the bolts tasks as well.

## MP3

### Installation

Deploy the code, cfg.json and compile the codes on VMs by running `ansible-playbook roles/deploy.yml` under the `ansible` directory. Also run `make buildlocal` to make sure the binaries are available locally.

### Run demo

1. `cp remotecfg.json cfg.json` to switch to the remote (VM) mode.
2. `make run` will spawn the processes on every VMs and setup the gRPC server and Failure detection UDP ports.
3. `make join` to connect all VMs together. Each VM will now have the full membership list.
4. To use the SDFS functionality, go to any live VM and run `sds sdfs put/get/ls/store/get-versions` commands.
5. To fail any nodes, say, VM1 and VM2, run `sds -n 0,1 close`.
6. To tear down the systems, run `sds close`.

### Notes

1. Local files are stored under data/mp3/ directories of each VM.
2. Logs are recorded under data/mp2/ directories.

## MP2

### Installation

Deploy the source files to VMs and compile using the Makefile. The procedure is described in the MP1 session. After compilation, an extra command `dswim` will be installed to `$GOPATH/bin/dswim`. This will be the command that we'll use for this MP.

### Run demo

The steps needed to play the demo is below:
1. ssh to one VM and run `dswim -i`. The flag `-i` here implying the dswim will be the introducer. Use `ls` for the stdin input to retrieve the id, which will be used by other nodes to join later.
2. ssh to other VM and run `dswim`. This will block on user input. type `join <id of the introducer>` to join the group. Repeat this process on as many VMs as you like. You can type `ls` in any VM to get its local membership list and its own ID (i.e. IP:port)
3. Ctrl-C to kill one of the dswim process to simulate a node failing.
4. type leave to actively leave the group. Instead of failure that requires a timeout to be detected, active leave will disseminate its leave immediately.

## MP1

### Installation

Install grpc with
`go get google.golang.org/grpc`, and protoc-gen-go with `go get -u github.com/golang/protobuf/protoc-gen-go`

Compile the `.proto` file using

`protoc -I protobuf/ protobuf/grep_log.proto --go_out=plugins=grpc:protobuf`

## How to run

0. Copy SOURCEME.sh.example to SOURCEME.sh and change the value in latter.
1. Modify the `remotecfg.json` file to add the appropriate IP addresses and port numbers of the VMs.
2. Refer to [vmsetup/README.md](vmsetup/README.md) to spawn the corresponding VMs.
3. Add 'CS425NETID' environmental variable by `source SOURCEME.sh`.
4. `make`

## Tricks

1. To return the line number, need to add `-n`.
2. To return the file name, need to add another null file (e.g. /dev/null)
3. The wildcard "*" will not work in exec.Command. Have to do the filepath.Glob manually, or invoke a shell `/bin/sh -c ...`
4. Invoking sh -c results in "#" signs unquoted.
5. Go install will make the binary filename based on package name (which is the name of the directory containing "main"), not the name of file containing "main" function.
6. To uninstall binaries, use `go clean -i mp/...`
7. gRPC naming conventions [link](https://developers.google.com/protocol-buffers/docs/style)

## TODO:

- [Done] Add json file to support input IP address.
- [Done] Allow user-input grep pattern. (Currently hard-coded)
- Write command that will cleanup the running service instead of manually cancel each server.
- Write code to auto-send IP & port info from server to the client.
- [Done] Add build functionality instead of running `go run` all the time.
