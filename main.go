package main

import (
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/input"
	"os"
)

var file string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid input. Input should be: ./task.exe <file_name>")
		return
	}
	fileName := os.Args[1]
	if file != "" {
		fileName = file
	}
	buffer, f := input.Input(fileName, false)
	err := buffer.Flush()
	if err != nil {
		fmt.Printf("Error with buffer flush: %v\n", err)
	}
	f.Close()
}
