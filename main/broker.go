package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

func logf(s string, p ...interface{}) {
	log.Printf(s+"\n", p...)
}

var chanWorker = make(chan net.Conn)

func addWorker(w net.Conn) {
	chanWorker <- w
}

func do() error {
	workerListener, err := net.Listen("tcp", brokerWorkerAddress)
	if err != nil {
		return err
	}
	clientListener, err := net.Listen("tcp", brokerClientAddress)
	if err != nil {
		return err
	}

	// Queue Workers
	go func() {
		for {
			workerConn, err := workerListener.Accept()
			if err != nil {
				logf("Accepting worker connection: %s", err.Error())
				continue
			}
			logf("Got Worker [%s]", workerConn.RemoteAddr())
			go addWorker(workerConn)
		}
	}()

	// Loop Clients
	for {
		clientConn, err := clientListener.Accept()
		if err != nil {
			return err
		}
		logf("Got Client [%s]", clientConn.RemoteAddr())
		match(clientConn)
	}

}

func match(clientConn net.Conn) {
	workerConn := <-chanWorker // Only accept a client request once it has a worker
	logf("Match C[%s] to W[%s]", clientConn.RemoteAddr(), workerConn.RemoteAddr())
	clientConn.Write([]byte{0x00})
	go linkConnections(clientConn, workerConn)
	go linkConnections(workerConn, clientConn)
}

func linkConnections(connA net.Conn, connB net.Conn) {
	io.Copy(connA, connB)
	// The error or number of bytes isn't really important
	// Close both if either fail / close.
	connA.Close()
	connB.Close()
}
