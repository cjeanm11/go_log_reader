package processor

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"log"
)

// string fileTypes
const (
	bx2     string = "bx2"
	gz             = "gz"
	unknown        = "unknown"
)

var decompressionMap = map[string]DecompressionFunc{
	bx2: decompressBZ2,
	gz:  decompressGZ,
	// Add other decompression functions as needed
}

type DecompressionFunc func(data []byte) []byte

// Checks for byte order mark
func getFileType(chunk []byte) string {
	// Implement additional checks for other file types if needed
	// If no known magic bytes match, return "unknown"
	if bytes.HasPrefix(chunk, []byte{0x42, 0x5A}) { // Check for the BZ2 magic bytes
		return bx2
	} else if bytes.HasPrefix(chunk, []byte{0x1F, 0x8B}) { // Check for the GZ magic bytes
		return gz
	} else {
		if len(chunk) > 2 {
			return fmt.Sprintf("% x", chunk[:2])
		}
		return unknown
	}
}

func DecompressStream(inputData []byte, options map[string]interface{}) io.Reader {
	var selectedStream io.Reader

	fileType := getFileType(inputData)
	decompressionFunc, exists := decompressionMap[fileType]

	if exists {
		selectedStream = bytes.NewReader(decompressionFunc(inputData))
	} else {
		fmt.Println("Unknown file format :", fileType)
		selectedStream = bytes.NewReader(inputData)
	}

	return selectedStream
}

func decompressBZ2(data []byte) []byte {
	buf := bytes.NewBuffer(data)
	decompressed := bytes.NewBuffer(nil)

	reader := bzip2.NewReader(buf)
	if _, err := io.Copy(decompressed, reader); err != nil {
		log.Fatal(err)
	}

	return decompressed.Bytes()
}

func decompressGZ(data []byte) []byte {
	buf := bytes.NewBuffer(data)
	decompressed := bytes.NewBuffer(nil)

	reader, err := gzip.NewReader(buf)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(decompressed, reader); err != nil {
		log.Fatal(err)
	}

	return decompressed.Bytes()
}
