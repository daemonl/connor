package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/daemonl/connor"
)

var brokerClientAddress string
var brokerWorkerAddress string

func init() {
	flag.StringVar(&brokerClientAddress, "client", ":5556", "broker's client address")
	flag.StringVar(&brokerWorkerAddress, "worker", ":5555", "broker's worker address")
}

func main() {
	err := do()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
	fmt.Println("Clean Exit")
}

var chanWorker = make(chan net.Conn)
var chanClient = make(chan net.Conn)

func do() error {
	connor.Logger = connor.FuncLogger(func(m string) {
		log.Println(m)
	})

	// Queue Workers
	go connor.TinyHandshakeListen(brokerWorkerAddress, chanWorker)
	go connor.TinyHandshakeListen(brokerClientAddress, chanClient)

	// Loop Clients
	for {
		clientConn := <-chanClient
		log.Printf("Got Client [%s]\n", clientConn.RemoteAddr())
		workerConn := <-chanWorker // Only accept a client request once it has a worker
		log.Printf("Match C[%s] to W[%s]\n", clientConn.RemoteAddr(), workerConn.RemoteAddr())
		connor.BindConnections(clientConn, workerConn)
	}

}
