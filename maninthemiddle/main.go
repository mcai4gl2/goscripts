package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"sync"
)

func main() {
	address := flag.String("peer", "", "peer host:port")
	flag.Parse()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	self := "localhost:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	log.Println("Listening on: " + self)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(connection, *address)
	}
}

func serve(localConn net.Conn, remote string) {
	log.Printf("[%v] - estblished", localConn.RemoteAddr())
	remoteConn, err := net.Dial("tcp", remote)
	if err != nil {
		log.Fatal(err)
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		defer localConn.Close()
		defer remoteConn.Close()
		for {
			buffer := make([]byte, 1024)
			n, err := localConn.Read(buffer)
			if err != nil {
				break
			}
			log.Printf("Outgoing message: %s", string(buffer[:n]))
			n, err = remoteConn.Write(buffer[:n])
			if err != nil {
				break
			}
		}
	}()

	go func() {
		defer waitGroup.Done()
		defer localConn.Close()
		defer remoteConn.Close()
		for {
			buffer := make([]byte, 1024)
			n, err := remoteConn.Read(buffer)
			if err != nil {
				break
			}
			log.Printf("Incoming message: %s", string(buffer[:n]))
			n, err = localConn.Write(buffer[:n])
			if err != nil {
				break
			}
		}
	}()

	waitGroup.Wait()
	log.Printf("[%v] - closed", localConn.RemoteAddr())
}
