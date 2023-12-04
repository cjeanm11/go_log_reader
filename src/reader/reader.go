package reader

import (
	"io"
	"log"
)

// TODO general : add more specific error messsage (errors.go)

// StreamHandler represents a generic stream handler interface
type StreamHandler interface {
	OnMessage(callback func([]byte))
}

// GenericStream is a generic stream handler
type GenericStream struct {
	MessageChannel chan []byte
}

// NewGenericStream creates a new GenericStream
func NewGenericStream() *GenericStream {
	return &GenericStream{
		MessageChannel: make(chan []byte),
	}
}

// OnMessage listens for incoming messages on the stream
func (gs *GenericStream) OnMessage(callback func([]byte)) {
	for msg := range gs.MessageChannel {
		callback(msg)
	}
}

// StreamReader processes the input stream using a provided StreamHandler
func StreamReader(inputStream io.Reader, streamHandler StreamHandler) func(func([]byte)) {
	var isStarted bool

	return func(fn func([]byte)) {
		if !isStarted {
			isStarted = true
			go func() {
				defer close(streamHandler.(*GenericStream).MessageChannel)
				buf := make([]byte, 1024)
				for {
					n, err := inputStream.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Fatal(err)
						}
						break
					}
					if n == 0 {
						break
					}
					data := make([]byte, n)
					copy(data, buf[:n])
					streamHandler.(*GenericStream).MessageChannel <- data
				}
			}()
		}

		streamHandler.OnMessage(fn)
	}
}
