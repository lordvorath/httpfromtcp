package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		incData := getLinesChannel(netConn)
		for s := range incData {
			fmt.Println(s)
		}
		fmt.Printf("\nThe channel has been closed\n")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var buffer = make([]byte, 8)
	var str string
	var outData = make(chan string)
	go func() {
		defer f.Close()
		defer close(outData)
		for {
			n, err := f.Read(buffer)
			if err != nil {
				if str != "" {
					outData <- str
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			parts := strings.Split(string(buffer[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				str += parts[i]
				outData <- str
				str = ""
			}
			str += parts[len(parts)-1]
		}
	}()

	return outData
}
