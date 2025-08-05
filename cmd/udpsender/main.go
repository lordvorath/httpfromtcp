package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Errorf("failed to dial UDP address: %v", err)
	}

	buff := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		inp, err := buff.ReadString('\n')
		if err != nil {
			fmt.Errorf("invalid input: %v", err)
		}

		_, err = udpConn.Write([]byte(inp))
		if err != nil {
			fmt.Errorf("failed to write to UDP connection: %v", err)
		}
	}
}
