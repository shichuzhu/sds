package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	wg   sync.WaitGroup
	lock sync.Mutex
)

func testDial() {
	conn1, err1 := net.Dial("tcp", "127.0.0.1:10000")
	defer conn1.Close()
	if err1 != nil {
		fmt.Println("CANNOT CONNECT TO SERVER")
	} else {
		fmt.Println("DIAL CORRECTLY")
	}
	time.Sleep(time.Second * 3)

	conn2, err2 := net.Dial("tcp", "127.0.0.1:10000")
	defer conn2.Close()
	if err2 != nil {
		fmt.Println("RECONNECTION ERROR")
	} else {
		fmt.Println("DIAL CORRECTLY FOR REDIAL")
	}
}

func connect(i int, sig *int) {
	defer wg.Done()
	conn, err := net.Dial("tcp", "127.0.0.1:10000")
	defer conn.Close()
	if err != nil {
		fmt.Printf("CONNECTION FAILURE IN CLIENT : %d", i)
	} else {
		lock.Lock()
		*sig += 1
		lock.Unlock()
	}
}

func massiveClient() {
	numClient := 500
	ch := 0
	wg.Add(numClient)
	for i := 0; i < numClient; i++ {
		go connect(i, &ch)
	}
	wg.Wait()
	if ch < numClient {
		fmt.Println("FAIL MASSIVE CONNECTION TEST.")
	} else {
		fmt.Println("PASS MASSIVE CONNECTION TEST.")
	}
}

func Testmain() {
	testDial()
	massiveClient()

}
