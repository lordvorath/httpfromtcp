package main

import (
	"fmt"
	"log"
	"net"

	"github.com/lordvorath/httpfromtcp/internal/request"
)

func main() {
	tcpListener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("failed to make listener: %v", err)
	}
	defer tcpListener.Close()

	for {
		netConn, err := tcpListener.Accept()
		if err != nil {
			log.Fatalf("failed to establish connection: %v", err)
		}

		fmt.Println("Connection established from:", netConn.RemoteAddr())

		req, err := request.RequestFromReader(netConn)
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))

	}
}
