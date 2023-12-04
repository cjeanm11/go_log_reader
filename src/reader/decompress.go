package reader

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
)

// TODO general : add more specific error messsage (errors.go)

const (
	bx2 = "bx2"
	gz  = "gz"
	// TODO add 7z format
	unknown = "unknown"
)

var decompressionMap = map[string]DecompressionFunc{
	bx2: decompressBZ2,
	gz:  decompressGZ,
}

type DecompressionFunc func(data *[]byte) (io.Reader, error)

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

func DecompressStream(inputData *[]byte) (func(func([]byte)), error) {
	fileType := getFileType(inputData)
	decompressionFunc, exists := decompressionMap[fileType]

	if !exists {
		buf := bytes.NewBuffer(*inputData)
		return StreamReader(buf, NewGenericStream()), nil
	}

	decompressedData, err := decompressionFunc(inputData)
	if err != nil {
		return nil, err
	}

	return StreamReader(decompressedData, NewGenericStream()), nil
}

func decompressBZ2(data *[]byte) (io.Reader, error) {
	buf := bytes.NewBuffer(*data)
	decompressed := bytes.NewBuffer(nil)

	reader := bzip2.NewReader(buf)
	if _, err := io.Copy(decompressed, reader); err != nil {
		return nil, err
	}

	return decompressed, nil
}

func decompressGZ(data *[]byte) (io.Reader, error) {
	buf := bytes.NewBuffer(*data)
	decompressed := bytes.NewBuffer(nil)

	reader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(decompressed, reader); err != nil {
		return nil, err
	}

	return decompressed, nil
}
