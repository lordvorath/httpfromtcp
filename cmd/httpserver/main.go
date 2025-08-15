package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/lordvorath/httpfromtcp/internal/request"
	"github.com/lordvorath/httpfromtcp/internal/response"
	"github.com/lordvorath/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, myHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func myHandler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		handleChunked(w, req)
		return
	}
	var code response.StatusCode = response.StatusOk
	var message string = "OK"
	var title string = "Success!"
	var body string = "Your request was an absolute banger."

	if req.RequestLine.RequestTarget == "/yourproblem" {
		code = response.StatusBadRequest
		message = "Bad Request"
		title = "Bad Request"
		body = "Your request honestly kinda sucked."
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		code = response.StatusInternalError
		message = "Internal Server Error"
		title = "Internal Server Error"
		body = "Okay, you know what? This one is on me.."
	}

	html := `<html><head><title>$CODE $MESSAGE</title></head><body><h1>$TITLE</h1><p>$BODY</p></body></html>`
	html = strings.Replace(html, "$CODE", strconv.Itoa(int(code)), -1)
	html = strings.Replace(html, "$MESSAGE", message, -1)
	html = strings.Replace(html, "$TITLE", title, -1)
	html = strings.Replace(html, "$BODY", body, -1)

	w.WriteStatusLine(code)

	headers := response.GetDefaultHeaders(len(html))
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)

	_, err := w.WriteBody([]byte(html))
	if err != nil {
		return
	}
}

const chunkSize = 1024

func handleChunked(w *response.Writer, req *request.Request) {
	//get data
	//data := []byte(`"Host": "httpbin.org"`) //CHEATING, because httpbin.org is down

	// From the solution because httpbin.org is down
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target
	fmt.Println("Proxying to", url)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	data := make([]byte, chunkSize)

	//send response
	var code response.StatusCode = response.StatusOk
	headers := response.GetDefaultHeaders(0)
	headers.Set("Content-Type", "text/html")
	headers.Set("Transfer-Encoding", "chunked")
	headers.Set("Trailer", "X-Content-SHA256, X-Content-Length")

	w.WriteStatusLine(code)
	w.WriteHeaders(headers)

	fullBody := make([]byte, 0)
	for {
		n, err := resp.Body.Read(data)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			_, err := w.WriteChunkedBody(data[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
			fullBody = append(fullBody, data[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}

	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Printf("failed to write end of chunked body: %v", err)
		return
	}

	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	headers.Set("X-Content-SHA256", sha256)
	ll := fmt.Sprintf("%d", len(fullBody))
	headers.Set("X-Content-Length", ll)

	err = w.WriteTrailers(headers)
	if err != nil {
		fmt.Printf("error writing trailers: %v", err)
	}

}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}
