package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"regexp"
	"strings"
)

var (
	host = flag.String("host", "localhost", "host")
	port = flag.String("port", "8080", "port")
	// this abomination is needed, because go's Json.Number or big.Int's
	// setString actually accepts non number strings like '[123]' :) *note* the
	// '[]' brackets
	NUMBER, _ = regexp.Compile(`^[+-]?([0-9]*[.])?[0-9]+$`)
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
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected: ", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("[%s] client closed connection.\n", conn.RemoteAddr())
				return // EOF
			}
			fmt.Println("Closing connection, read failed, err: ", err.Error())
			return // error
		}
		request, err := parseRequest(data)
		if err != nil {
			malformedReq := fmt.Sprintf("error: malformed request: %s\n", err.Error())
			fmt.Printf("Got malformed request. Closing connection. \n%+v\n%s\n", request, err.Error())
			conn.Write([]byte(malformedReq))
			return
		}
		ans := process(*request)
		response := serialize(ans)
		fmt.Printf("[%s] Request: %s\n", conn.RemoteAddr(), *request)
		fmt.Printf("[%s] Response: %+v\n", conn.RemoteAddr(), ans)
		conn.Write(fmt.Appendf(response, "\n"))
	}
}

type Request struct {
	Method string          `json:"method"`
	Number json.RawMessage `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func serialize(response Response) []byte {
	data, _ := json.Marshal(response)
	return data
}

func process(req Request) Response {
	numStr := string(req.Number)
	falseResp := Response{Method: "isPrime", Prime: false}
	// if float return false
	if strings.ContainsAny(numStr, ".eE") {
		return falseResp
	}
	n := new(big.Int)
	// SetString returns (n, ok) — ok is false if the string isn't a valid integer
	if _, ok := n.SetString(string(req.Number), 10); !ok {
		// Defensive: treat unparseable as not prime, though parseRequest should catch it
		return falseResp
	}
	return Response{Method: "isPrime", Prime: n.ProbablyPrime(20)}
}

func parseRequest(data []byte) (*Request, error) {
	var body Request
	fmt.Printf("Parsing body %s", data)
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("Failed to parse request body: %s", err.Error())
	}
	// validations
	if body.Method != "isPrime" {
		return nil, fmt.Errorf("invalid method: %s", body)
	}

	if matched := NUMBER.Match(body.Number); !matched {
		return nil, fmt.Errorf("invalid number, should be an int or float, got: %s", body)
	}
	return &body, nil
}
