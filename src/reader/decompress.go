package reader

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
)

const (
	bx2     = "bx2"
	gz      = "gz"
	unknown = "unknown"
)

var decompressionMap = map[string]DecompressionFunc{
	bx2: decompressBZ2,
	gz:  decompressGZ,
}

type DecompressionFunc func(data *[]byte) ([]byte, error)

func getFileType(chunk *[]byte) string {
	if bytes.HasPrefix(*chunk, []byte{0x42, 0x5A}) {
		return bx2
	} else if bytes.HasPrefix(*chunk, []byte{0x1F, 0x8B}) {
		return gz
	} else {
		if len(*chunk) > 2 {
			return fmt.Sprintf("% x", (*chunk)[:2])
		}
		return unknown
	}
}

func DecompressStream(inputData *[]byte) (io.Reader, error) {
	fileType := getFileType(inputData)
	decompressionFunc, exists := decompressionMap[fileType]

	if !exists {
		return bytes.NewReader(*inputData), nil
	}

	decompressedData, err := decompressionFunc(inputData)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(decompressedData), nil
}

func decompressBZ2(data *[]byte) ([]byte, error) {
	buf := bytes.NewBuffer(*data)
	decompressed := bytes.NewBuffer(nil)

	reader := bzip2.NewReader(buf)
	if _, err := io.Copy(decompressed, reader); err != nil {
		return nil, err
	}

	return decompressed.Bytes(), nil
}

func decompressGZ(data *[]byte) ([]byte, error) {
	buf := bytes.NewBuffer(*data)
	decompressed := bytes.NewBuffer(nil)

	reader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(decompressed, reader); err != nil {
		return nil, err
	}

	return decompressed.Bytes(), nil
}
