package main

import (
	"bufio"
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/club"
	"github.com/Anoke/YADRO-IMPULSE/internal/eventHandler"
	"github.com/Anoke/YADRO-IMPULSE/internal/validation"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myprogram <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Открытие файла для чтения
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Error opening input file: %v", err)
	}
	defer file.Close()

	// Буфер для накопления вывода
	var outputBuffer strings.Builder

	// Сканер для построчного чтения файла
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	var numTables int
	var startTime, endTime time.Time
	var hourlyRate int
	var compClub *club.ComputerClub
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		switch lineNumber {
		case 1:
			numTables, err = strconv.Atoi(line)
			if err != nil {
				fmt.Println(line)
			}
		case 2:
			startTime, endTime, err = validation.ParseWorkHours(line)
			if err != nil {
				fmt.Println(line)
			}
		case 3:
			hourlyRate, err = strconv.Atoi(line)
			if err != nil || hourlyRate <= 0 {
				fmt.Println(line)
			}
		default:
			if validation.IsEventStringValid(line) != nil {
				fmt.Println(line)
			}

			compClub = club.NewComputerClub(numTables, startTime, endTime, hourlyRate)
			//тут надо написать метод для выполнения события
			parts := strings.Fields(line)
			eventTime := parts[0]
			eventId := parts[1]
			eventBody := strings.Join(parts[2:], " ")
			err = eventHandler.HandleEvent(eventTime, eventId, eventBody, compClub)

		}
		if err := scanner.Err(); err != nil {
			fmt.Errorf("error scanning file: %v", err)
		}
	}

	// Проверка наличия ошибок
	if outputBuffer.Len() > 0 {
		fmt.Println(outputBuffer.String())
	} else {
		fmt.Println("All events processed successfully.")
	}
	// Инициализация компьютерного клуба
	myClub := club.NewComputerClub(tables, startTime, endTime, hourlyRate)

	// Обработка каждого события
	for _, event := range events {
		err := eventhandler.HandleEvent(myClub, event)
		if err != nil {
			log.Printf("Error handling event: %v", err)
		}
	}

	// Вывод результатов
	fmt.Println("Start time:", myClub.StartTime)
	fmt.Println("Events for the day:")
	for _, event := range events {
		fmt.Println(event)
	}
	fmt.Println("End time:", myClub.EndTime)

	fmt.Println("Revenue and usage per table:")
	for _, table := range myClub.Tables {
		fmt.Printf("Table %d: Revenue=%d, UsageTime=%s\n", table.Number, table.Revenue, table.UsageTime)
	}
}
