package input

import (
	"bufio"
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/validation"
	"os"
	"strconv"
	"time"
)

// ValidateInputFileFirst validates input values of first 3 strings to get basic info about club
func ValidateInputFileFirst(filename string) (int, time.Time, time.Time, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, time.Time{}, time.Time{}, 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	var numTables int
	var startTime, endTime time.Time
	var hourlyRate int

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		switch lineNumber {
		case 1:
			numTables, err = strconv.Atoi(line)
			if err != nil {
				return 0, time.Time{}, time.Time{}, 0, fmt.Errorf(line)
			}
		case 2:
			startTime, endTime, err = validation.ParseWorkHours(line)
			if err != nil {
				return 0, time.Time{}, time.Time{}, 0, fmt.Errorf(line)
			}
		case 3:
			hourlyRate, err = strconv.Atoi(line)
			if err != nil || hourlyRate <= 0 {
				return 0, time.Time{}, time.Time{}, 0, fmt.Errorf(line)
			}
		default:
			if validation.IsEventStringValid(line) != nil {
				return 0, time.Time{}, time.Time{}, 0, fmt.Errorf(line)
			}

			if lineNumber > 3 {
				return numTables, startTime, endTime, hourlyRate, nil
			}

		}
	}

	if err := scanner.Err(); err != nil {
		return 0, time.Time{}, time.Time{}, 0, fmt.Errorf("error scanning file: %v", err)
	}

	return numTables, startTime, endTime, hourlyRate, nil
}

func ValidateEventLine(line string, lineNumber int) (int, time.Time, time.Time, int, error) {
	var numTables int
	var startTime, endTime time.Time
	var hourlyRate int
	var err error

	switch lineNumber {
	case 1:
		numTables, err = strconv.Atoi(line)
		if err != nil {
			return 0, time.Time{}, time.Time{}, 0, fmt.Errorf(line)
		}
	case 2:
		startTime, endTime, err = validation.ParseWorkHours(line)
		if err != nil {
			return 0, time.Time{}, time.Time{}, 0, fmt.Errorf(line)
		}
	case 3:
		hourlyRate, err = strconv.Atoi(line)
		if err != nil || hourlyRate <= 0 {
			return 0, time.Time{}, time.Time{}, 0, fmt.Errorf(line)
		}
	default:
		if validation.IsEventStringValid(line) != nil {
			return 0, time.Time{}, time.Time{}, 0, fmt.Errorf(line)
		}
	}
	return numTables, startTime, endTime, hourlyRate, nil
}
