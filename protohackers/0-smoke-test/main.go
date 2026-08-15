package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	port  = flag.String("port", ":8080", "port for listening")
	debug = flag.Bool("debug", false, "set debug logging")
)

func main() {
	flag.Parse()

	tcp, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to create listener, err: %s", err)
	}
	defer tcp.Close()
	fmt.Printf("Listening on port - %s...\n", *port)
	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Printf("failed to accept a connection, err: %s", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Handling connection - %s\n", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err == io.EOF { // closing client sends io.EOF
				fmt.Printf("Closing client connection %s\n", conn.RemoteAddr())
				return
			}
			fmt.Println("failed to read data, err:", err)
			return
		}
		fmt.Printf("Client said: %s", bytes)
		prefix := ""
		if *debug {
			prefix = "Server replied with - "
		}
		response := fmt.Sprintf("%s%s", prefix, bytes)
		conn.Write([]byte(response))
	}
}
