package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqSlice, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqParts := strings.Split(string(reqSlice), "\r\n")

	reqLine, err := parseRequestLine(reqParts[0])
	if err != nil {
		return nil, err
	}

	return &Request{reqLine}, nil
}

func parseRequestLine(line string) (RequestLine, error) {
	reqLine := strings.Split(line, " ")
	if len(reqLine) != 3 {
		return RequestLine{}, fmt.Errorf("bad request-line: expected 3 parts, got %d", len(reqLine))
	}

	method := reqLine[0]
	target := reqLine[1]
	version := strings.TrimPrefix(reqLine[2], "HTTP/")

	if (strings.ToUpper(method) != method) || (strings.ContainsAny(method, "0123456789")) {
		return RequestLine{}, fmt.Errorf("invalid method: %s", method)
	}

	if version != "1.1" {
		return RequestLine{}, fmt.Errorf("unsupported HTTP version. Only 1.1 is supported")
	}

	return RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, nil
}
