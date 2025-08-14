package main

import (
	"log"
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

func myHandler(w response.Writer, req *request.Request) {
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
