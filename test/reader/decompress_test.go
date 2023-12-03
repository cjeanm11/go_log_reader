package reader

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"go_log_reader/src/reader"
	"io"
	"log"
	"os"
	"testing"

	"github.com/dsnet/compress/bzip2"
)

func compress2Bzip2(input []byte) ([]byte, error) {
	var compressed bytes.Buffer

	// Create a Bzip2 writer with default compression settings
	writer, err := bzip2.NewWriter(&compressed, nil)
	if err != nil {
		return nil, err
	}

	// Write the input data to the writer
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

	// Create a Gzip writer with default compression settings
	writer := gzip.NewWriter(&compressed)

	// Write the input data to the writer
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

func TestDecompress(t *testing.T) {
	inputData := []byte("SampleInput") // Replace this with your actual input data
	selectedStream, err := reader.DecompressStream(&inputData)
	size, err := io.Copy(os.Stdout, selectedStream)
	fmt.Printf("size : %d\n", size)
	if err != nil {
		t.Fatal(err)
	}

	inputData = []byte("") // Replace this with your actual input data
	selectedStream, err = reader.DecompressStream(&inputData)
	size, err = io.Copy(os.Stdout, selectedStream)
	fmt.Printf("size : %d\n", size)
	if err != nil {
		t.Fatal(err)
	}

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

	size, err := io.Copy(os.Stdout, decompressedStream)
	fmt.Printf("size : %d \n", size)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecompress_GZ(t *testing.T) {

	var sampleInput []byte = nil
	textGZ := "This is a sample input string."
	sampleInput = append(sampleInput, []byte(textGZ)...)

	// Compress the data using Bzip2
	compressedData, err := compress2Gzip(sampleInput)
	if err != nil {
		log.Fatal("Error compressing data:", err)
	}

	if len(compressedData) == 0 {
		log.Fatal("Compressed data is empty")
	}

	fmt.Printf("Compressed data length: %d bytes\n", len(compressedData))

	decompressedStream, err := reader.DecompressStream(&compressedData)

	size, err := io.Copy(os.Stdout, decompressedStream)
	fmt.Printf("size : %d \n", size)
	if err != nil {
		t.Fatal(err)
	}
}
