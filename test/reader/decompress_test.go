package reader

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"go_log_reader/src/reader"
	"log"
	"testing"

	"github.com/dsnet/compress/bzip2"
)

func compress2Bzip2(input []byte) ([]byte, error) {
	var compressed bytes.Buffer

	writer, err := bzip2.NewWriter(&compressed, nil)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(input)
	if err != nil {
		writer.Close()
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}

func compress2Gzip(input []byte) ([]byte, error) {
	var compressed bytes.Buffer

	writer := gzip.NewWriter(&compressed)

	_, err := writer.Write(input)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}

func readAndPrintStream(callback func(func([]byte))) {
	callback(func(msg []byte) {
		fmt.Println("Received message:", string(msg))
		if msg == nil {
			log.Fatal("no data received")
		}
	})

}

func TestDecompress(t *testing.T) {
	inputData := []byte("SampleInput")
	selectedStream, err := reader.DecompressStream(&inputData)
	if err != nil && selectedStream != nil {
		t.Fatal(err)
	}
	readAndPrintStream(selectedStream)

	inputData = []byte("")
	decompressedStream, err := reader.DecompressStream(&inputData)

	if err != nil && decompressedStream != nil {
		t.Fatal(err)
	}
	readAndPrintStream(decompressedStream)

}

func TestDecompress_Bzip2(t *testing.T) {

	var sampleInput []byte = nil
	textBzip2 := "This is a sample input string."
	sampleInput = append(sampleInput, []byte(textBzip2)...)

	// Compress the data using Bzip2
	compressedData, err := compress2Bzip2(sampleInput)
	if err != nil {
		log.Fatal("Error compressing data:", err)
	}

	if len(compressedData) == 0 {
		log.Fatal("Compressed data is empty")
	}

	fmt.Printf("Compressed data length: %d bytes\n", len(compressedData))

	decompressedStream, err := reader.DecompressStream(&compressedData)

	if err != nil && decompressedStream != nil {
		t.Fatal(err)
	}
}

func TestDecompress_GZ(t *testing.T) {

	var sampleInput []byte = nil
	textGZ := "This is a sample input string."
	sampleInput = append(sampleInput, []byte(textGZ)...)

	compressedData, err := compress2Gzip(sampleInput)
	if err != nil {
		log.Fatal("Error compressing data:", err)
	}

	if len(compressedData) == 0 {
		log.Fatal("Compressed data is empty")
	}

	fmt.Printf("Compressed data length: %d bytes\n", len(compressedData))

	decompressedStream, err := reader.DecompressStream(&compressedData)

	if err != nil && decompressedStream != nil {
		t.Fatal(err)
	}
}
