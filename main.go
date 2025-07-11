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

	incData := getLinesChannel(f)
	for s := range incData {
		fmt.Printf("read: %s\n", s)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var buffer = make([]byte, 8)
	var str string
	var outData = make(chan string)
	go func() {
		defer f.Close()
		defer close(outData)
		for {
			n, err := f.Read(buffer)
			if err != nil {
				if str != "" {
					outData <- str
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
				str += parts[i]
				outData <- str
				str = ""
			}
			str += parts[len(parts)-1]
		}
	}()

	return outData
}
