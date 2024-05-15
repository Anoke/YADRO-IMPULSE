package input

import (
	"bufio"
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/club"
	"github.com/Anoke/YADRO-IMPULSE/internal/eventHandler"
	"github.com/Anoke/YADRO-IMPULSE/internal/output"
	"github.com/Anoke/YADRO-IMPULSE/internal/validation"
	"os"
	"strconv"
	"strings"
	"time"
)

// Input parses input and calls event handlers
func Input(fileName string, test bool) (*bufio.Writer, *os.File) {
	var whereToPut *os.File
	if test {
		var err error
		whereToPut, err = os.Create("temp.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		whereToPut = os.Stdout
	}
	outputBuffer1 := bufio.NewWriter(whereToPut)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("File open error:", err)
		return nil, whereToPut
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	outputBuffer := output.CreateBufferOutput(whereToPut)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	var numTables int
	var startTime, endTime time.Time
	var hourlyRate int
	var compClub *club.ComputerClub
	firstTime := true
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		switch lineNumber {
		// Parse quantity of tables in club
		case 1:
			numTables, err = strconv.Atoi(line)
			if err != nil {
				output.AddToBuffer(line)
				return outputBuffer, whereToPut
			}
		// Parses work hours
		case 2:
			startTime, endTime, err = validation.ParseWorkHours(line)
			if err != nil {
				output.AddToBuffer(line)
				return outputBuffer, whereToPut
			}
		// Parses cost of one hour
		case 3:
			hourlyRate, err = strconv.Atoi(line)
			if err != nil || hourlyRate <= 0 {
				output.AddToBuffer(line)
				return outputBuffer, whereToPut
			}
		// Event processing and execution
		default:
			//Validation of event
			if validation.IsEventStringValid(line) != nil {
				_, _ = outputBuffer1.WriteString(line + "\n")
				return outputBuffer1, whereToPut
			}

			compClub = club.NewComputerClub(numTables, startTime, endTime, hourlyRate)
			if firstTime {
				output.AddToBuffer(fmt.Sprintf("%s", startTime.Format("15:04")))
				firstTime = false
			}
			parts := strings.Fields(line)
			eventTime := parts[0]
			eventId := parts[1]
			eventBody := strings.Join(parts[2:], " ")
			//Execution of event
			eventHandler.HandleEvent(eventTime, eventId, eventBody, compClub)

		}
		if err = scanner.Err(); err != nil {
			err1 := fmt.Errorf("error scanning file: %v", err)
			if err1 != nil {
				_, _ = outputBuffer1.WriteString(err1.Error() + "\n")
				return outputBuffer1, whereToPut
			}
		}
	}

	if compClub == nil {
		return nil, whereToPut
	}
	// Ends day at club
	compClub.EndDay()
	return outputBuffer, whereToPut
}
