package input

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestInput1(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  "../../testData/test1",
			out: "../../testData/result1",
		},
		{
			in:  "../../testData/test2",
			out: "../../testData/result2",
		},
		{
			in:  "../../testData/test3",
			out: "../../testData/result3",
		},
		{
			in:  "../../testData/test4",
			out: "../../testData/result4",
		},
	}

	var file *os.File
	for _, test := range tests {
		var outputBuffer *bufio.Writer
		outputBuffer, file = Input(test.in, true)
		if outputBuffer == nil {
			t.Fatalf("Expected non-nil output buffer")
		}
		err := outputBuffer.Flush()
		if err != nil {
			t.Fatalf("Error flushing output buffer: %v", err)
		}

		content1, err := ioutil.ReadFile("temp.txt")
		if err != nil {
			fmt.Println("Ошибка при чтении файла 1:", err)
			return
		}

		content2, err := ioutil.ReadFile(test.out)
		if err != nil {
			fmt.Println("Ошибка при чтении файла 2:", err)
			return
		}

		str1 := string(content1)
		str2 := string(content2)

		if str1 != str2 {
			t.Errorf("Содержимое файлов различается.")
		}
	}
	file.Close()
}
