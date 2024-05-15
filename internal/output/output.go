package output

import (
	"bufio"
	"os"
)

// field for buffer
var outputBuffer *bufio.Writer

// CreateBufferOutput creates new instance of buffer
func CreateBufferOutput(whereToPut *os.File) *bufio.Writer {
	outputBuffer = bufio.NewWriter(whereToPut)
	return outputBuffer
}

// AddToBuffer adds string to buffer
func AddToBuffer(str string) {
	_, _ = outputBuffer.WriteString(str + "\n")
}
