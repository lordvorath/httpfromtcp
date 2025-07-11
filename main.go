package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	filename := "messages.txt"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open %s: %v", filename, err)
		return
	}
	defer f.Close()

	var s = make([]byte, 8)
	for {
		n, err := f.Read(s)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("unexpected error while reading file: %v", err)
			break
		}

		str := s[:n]
		fmt.Printf("read: %s\n", str)
	}
}
