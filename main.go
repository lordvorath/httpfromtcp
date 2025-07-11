package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	filename := "messages.txt"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open %s: %v", filename, err)
		return
	}
	defer f.Close()

	var buffer = make([]byte, 8)
	var str string
	for {
		n, err := f.Read(buffer)
		if err != nil {
			if str != "" {
				fmt.Printf("read: %s\n", str)
				str = ""
			}
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		parts := strings.Split(string(buffer[:n]), "\n")

		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", str, parts[i])
			str = ""
		}
		str += parts[len(parts)-1]
	}
}
