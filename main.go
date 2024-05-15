package main

import (
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/input"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid input. Input should be: ./task.exe <file_name>")
		return
	}
	fileName := os.Args[1]
	buffer, file := input.Input(fileName, false)
	err := buffer.Flush()
	if err != nil {
		fmt.Printf("Error with buffer flush: %v\n", err)
	}
	file.Close()
}
