package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/lordvorath/httpfromtcp/internal/request"
	"github.com/lordvorath/httpfromtcp/internal/response"
)

type Server struct {
	port          int
	listener      net.Listener
	handler       Handler
	serverRunning atomic.Bool
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to make listener: %v", err)
	}

	s := Server{
		port:     port,
		listener: listener,
		handler:  handler,
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
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		he := &HandlerError{
			StatusCode: int(response.StatusBadRequest),
			Message:    err.Error(),
		}
		he.Write(conn)
		return
	}

	writer := response.Writer{
		W: conn,
	}

	s.handler(writer, req)

}

func (he HandlerError) Write(conn net.Conn) error {
	code := strconv.Itoa(int(he.StatusCode))
	message := "HTTP/1.1 " + code + " " + he.Message + "\r\n"
	h := response.GetDefaultHeaders(len(message))
	err := response.WriteStatusLine(conn, response.StatusCode(he.StatusCode))
	if err != nil {
		return err
	}
	err = response.WriteHeaders(conn, h)
	if err != nil {
		return err
	}
	err = response.WriteBody(conn, []byte(message))
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}
