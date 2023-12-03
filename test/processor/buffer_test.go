package processor

import (
	"fmt"
	"go_log_reader/src/processor"
	"testing"
)

func TestBuffer(t *testing.T) {

	// Little-endian format
	buffer := []byte{
		0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, // 16 bytes
		0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00,
	}

	offset := 0 // Replace this with the desired offset

	var want uint32 = 16

	size := processor.ReadSize(buffer, offset)
	fmt.Printf("Message size: %d\n", size) // 16 int32
	if !(want == size) {
		t.Fatalf(`ReadSize(buffer, offset) = %q, want match for %#q`, size, want)
	}

	message := processor.ReadMessage(buffer, offset)
	if message != nil {
		fmt.Printf("Message content: % x\n", message)
	} else {
		fmt.Println("Failed to read message. Invalid buffer or offset.")
	}

	offset += int(size)
	message = processor.ReadMessage(buffer, offset)
	if message != nil {
		fmt.Printf("Message content: % x\n", message)
	} else {
		fmt.Println("Failed to read message. Invalid buffer or offset.")
	}
}
