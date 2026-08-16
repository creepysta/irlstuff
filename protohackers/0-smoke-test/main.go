package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	host = flag.String("host", "localhost", "host")
	port = flag.String("port", "8080", "port")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", *host, *port)
	tcp, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to create listener, err: %s", err)
	}
	defer tcp.Close()
	fmt.Printf("Listening on port - %s...\n", addr)
	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Printf("failed to accept a connection, err: %s", err)
			continue
		}
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	defer conn.Close()
	// if _, err := io.Copy(conn, conn); err != nil {
	// 	fmt.Println("failed to copy with error: ", err.Error())
	// }

	fmt.Println("Client connected: ", conn.RemoteAddr())
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			fmt.Printf("Client sent: %s\n", buf[:n])
			if _, werr := conn.Write(buf[:n]); werr != nil {
				fmt.Println("Closing connection, write failed, err: ", err.Error())
				return
			}
		}
		if err != nil {
			fmt.Println("Closing connection, read failed, err: ", err.Error())
			return // EOF or error
		}
	}
}
