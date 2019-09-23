package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
)

func echo(s net.Conn) {
	defer s.Close()

	log.Printf("[%v <-> %v]\n", s.LocalAddr(), s.RemoteAddr())
	b := bufio.NewReader(s)
	for {
		line, e := b.ReadBytes('\n')
		if e != nil {
			break
		}
		log.Printf("[%v <-> %v] - %s", s.LocalAddr(), s.RemoteAddr(), line)
		s.Write(line)
	}
	log.Printf("closed\n")
}

func main() {
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
		go echo(connection)
	}
}
