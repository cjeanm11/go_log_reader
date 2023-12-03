package processor

import (
	"encoding/binary"
)

func ReadSize(buf []byte, offset int) uint32 {
	if len(buf) < offset+8 {
		return 0
	}

	segments := binary.LittleEndian.Uint32(buf[offset:])
	if len(buf) < offset+4+int(segments*4) {
		return 0
	}

	localIndex := 4
	var totalSize uint32

	for i := uint32(0); i < segments; i++ {
		if len(buf) < offset+localIndex+4 {
			return 0
		}

		segSize := binary.LittleEndian.Uint32(buf[offset+localIndex:])
		totalSize += (segSize * 8)
		localIndex += 4
	}

	// Account for header size
	totalSize += 8

	return totalSize
}

func ReadMessage(buf []byte, offset int) []byte {
	size := ReadSize(buf, offset)
	if size == 0 || len(buf) < offset+int(size) {
		return nil
	}
	return buf[offset : offset+int(size)]
}
