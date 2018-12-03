all: onetest

onetest:
	sds sdfs put sample.txt sampletxt
	sleep 0.3
	sds sdfs put sample.txt sampletxt
	sleep 0.3

#all: deploy build run

deploy:
	vmsetup/deploy Copy src @go/src/fa18cs425mp/

clean:
	vmsetup/deploy For Each 'go clean -i fa18cs425mp/...'

build:
	vmsetup/deploy For Each 'go install fa18cs425mp/...'

run:
	ssh szhu28@fa18-cs425-g44-01.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 1 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-02.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 2 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-03.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 3 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-04.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 4 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-05.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 5 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-06.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 6 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-07.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 7 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-08.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 8 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-09.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 9 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	ssh szhu28@fa18-cs425-g44-10.cs.illinois.edu 'shopt -s huponexit ; dserver -nodeid 0 -port 10000 -datapath "data/mp2" > log.out 2> log.err < /dev/null' &
	sleep 2
	#	sds grep -n 1,2,3 -c ailure '*'
	sds grep -c 123456 '../mp1/*'
# The first blank before -c to let flag render it as non-flag argument
#	vmsetup/deploy Spawn Each '-port 10000 -dataPath "data/mp2"'

join:
	for i in $$(seq 0 8); do sds -n $$i swim join 172.22.156.148:11000 ; sleep 0.5 ; done

localsetup:
#	go get -u google.golang.org/grpc
#	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I src/pb/ src/pb/server_services.proto --go_out=plugins=grpc:src/pb

buildlocal:
	go install fa18cs425mp/src/...

runlocal: buildlocal
	dserver -port 10001 -pfd 11001 -nodeid 1 &
	dserver -port 10002 -pfd 11002 -nodeid 2 &
	dserver -port 10003 -pfd 11003 -nodeid 3 &
	dserver -port 10004 -pfd 11004 -nodeid 4 &
#	dserver -port 10005 -pfd 11005 -nodeid 5 &
#	dserver -port 10006 -pfd 11006 -nodeid 6 &
#	dserver -port 10007 -pfd 11007 -nodeid 7 &
#	dserver -port 10008 -pfd 11008 -nodeid 8 &
#	dserver -port 10009 -pfd 11009 -nodeid 9 &
#	dserver -port 10010 -pfd 11010 -nodeid 0 &
	sleep 0.5
	sds swim join 128.174.245.229:11001
	sleep 0.5
	sds crane examples/streamProcessing/exclamation
	#sds grep -c 123456 '*'

testbuild:
	go build -buildmode=plugin -o data/mp4/exclamation/plugin/exclamation.so fa18cs425mp/examples/streamProcessing/exclamation

test:
	test/mp1/runtest

.PHONY: clean all test
