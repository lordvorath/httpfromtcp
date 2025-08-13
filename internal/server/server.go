package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	port          int
	listener      net.Listener
	serverRunning atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to make listener: %v", err)
	}

	s := Server{
		port:     port,
		listener: listener,
	}
	s.serverRunning.Store(true)

	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	s.serverRunning.Store(false)
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("failed to close listener: %v", err)
	}
	return nil
}

func (s *Server) listen() {
	for s.serverRunning.Load() {
		netConn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("failed to establish connection: %v\n", err)
		}

		go s.handle(netConn)

	}
}

func (s *Server) handle(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!"))
	conn.Close()
}
