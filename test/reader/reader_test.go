package reader

import (
	"fmt"
	"go_log_reader/src/reader"
	"log"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	file, err := os.Open("../resources/files/input.log")
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	inputStream := file

	genericStream := reader.NewGenericStream()

	streamReader := reader.StreamReader(inputStream, genericStream)

	streamReader(func(msg []byte) {
		fmt.Println("Received message:", string(msg))
		if msg == nil {
			log.Fatal("no data received")
		}
	})
}
