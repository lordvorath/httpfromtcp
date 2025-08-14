package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/lordvorath/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOk            StatusCode = 200
	StatusError         StatusCode = 400
	StatusInternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine string
	switch statusCode {
	case StatusOk:
		statusLine = "HTTP/1.1 200 OK"
	case StatusError:
		statusLine = "HTTP/1.1 400 Bad Request"
	case StatusInternalError:
		statusLine = "HTTP/1.1 500 Internal Server Error"
	default:
		code := strconv.Itoa(int(statusCode))
		statusLine = "HTTP/1.1 " + code + " "
	}

	statusLine = statusLine + "\r\n"
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("failed to write status line: %v", err)
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Headers{}
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		hh := k + ": " + v + "\r\n"
		_, err := w.Write([]byte(hh))
		if err != nil {
			return fmt.Errorf("failed to write header line: %v", err)
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("failed to write end of headers: %v", err)
	}
	return nil
}

func WriteBody(w io.Writer, body []byte) error {
	_, err := w.Write(body)
	if err != nil {
		return fmt.Errorf("failed to write body: %v", err)
	}
	return nil
}
