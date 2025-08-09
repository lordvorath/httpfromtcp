package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const crlf = "\r\n"
const bufferSize = 8

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	ParserState requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0

	req := &Request{
		ParserState: requestStateInitialized,
	}

	for req.ParserState != requestStateDone {
		if readToIndex >= len(buffer) {
			newbuf := make([]byte, 2*len(buffer))
			copy(newbuf, buffer)
			buffer = newbuf
		}

		n, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.ParserState = requestStateDone
				break
			}
			return nil, err
		}

		readToIndex += n

		n, err = req.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[n:])
		readToIndex -= n
	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	line := string(data[:idx])
	reqLine, err := requestLineFromString(line)
	if err != nil {
		return nil, 0, err
	}

	return reqLine, idx + 2, nil
}

func requestLineFromString(line string) (*RequestLine, error) {
	reqLine := strings.Split(line, " ")
	if len(reqLine) != 3 {
		return nil, fmt.Errorf("bad request-line: expected 3 parts, got %d\n%s", len(reqLine), reqLine)
	}

	method := reqLine[0]
	target := reqLine[1]
	version := strings.TrimPrefix(reqLine[2], "HTTP/")

	if (strings.ToUpper(method) != method) || (strings.ContainsAny(method, "0123456789")) {
		return nil, fmt.Errorf("invalid method: %s", method)
	}

	if version != "1.1" {
		return nil, fmt.Errorf("unsupported HTTP version. Only 1.1 is supported")
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParserState {
	case requestStateInitialized:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.ParserState = requestStateDone
		r.RequestLine = *rl
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown parser state")
	}
}
