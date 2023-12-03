package reader

import (
	"encoding/binary"
	"go_log_reader/src/enums"
	"sync"
)

func ReadSize(buf *[]byte, offset int) (uint32, error) {
	if len(*buf) < offset+8 {
		return 0, enums.ErrInvalidData
	}

	segments := binary.LittleEndian.Uint32((*buf)[offset:])
	if len(*buf) < offset+4+int(segments*4) {
		return 0, enums.ErrInvalidData
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var totalSize uint32

	var err error // To store errors during processing

	for localIndex := 4; localIndex < 4+int(segments*4); localIndex += 4 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()

			if len(*buf) < offset+idx+4 {
				// Handle error condition
				err = enums.ErrInvalidData
				return
			}
			segSize := binary.LittleEndian.Uint32((*buf)[offset+idx:])
			totalSize += (segSize * 8)
		}(localIndex)
	}

	wg.Wait()

	// Account for header size
	totalSize += 8

	if err != nil {
		return 0, err
	}

	return totalSize, nil
}

func ReadMessage(buf *[]byte, offset int) ([]byte, error) {
	size, err := ReadSize(buf, offset)
	if err != nil {
		return nil, err
	}

	if size == 0 || len(*buf) < offset+int(size) {
		return nil, enums.ErrInsufficientData
	}
	return (*buf)[offset : offset+int(size)], nil
}
