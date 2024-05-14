package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	timeFormatRegex = regexp.MustCompile(`^\d{2}:\d{2}$`)
	clientNameRegex = regexp.MustCompile(`^[a-z0-9_-]+$`)
)

func ReadAndValidateInputFile(line string) ([]string, error) {
	return nil, nil
}
func ParseWorkHours(str string) (time.Time, time.Time, error) {
	parts := strings.Fields(str)
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf(str)
	}

	startTime, err := ParseTimeFormat(parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err := ParseTimeFormat(parts[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startTime, endTime, nil
}

func ParseTimeFormat(str string) (time.Time, error) {
	if !timeFormatRegex.MatchString(str) {
		return time.Time{}, fmt.Errorf("invalid format of time")
	}
	return time.Parse("15:04", str)
}

func IsEventStringValid(str string) error {
	parts := strings.Fields(str)

	if len(parts) < 3 {
		return fmt.Errorf(str)
	}

	eventTime := parts[0]
	eventId := parts[1]
	eventBody := strings.Join(parts[2:], " ")

	if _, err := ParseTimeFormat(eventTime); err != nil {
		return fmt.Errorf(str)
	}

	switch eventId {
	case "1":
		if err := ValidateClient(eventBody); err != nil {
			return fmt.Errorf(str)
		}
	case "2":
		if err := ValidateClientSitDownMiAmor(eventBody); err != nil {
			return fmt.Errorf(str)
		}
	case "3":
		if err := ValidateClient(eventBody); err != nil {
			return fmt.Errorf(str)
		}
	case "4":
		if err := ValidateClient(eventBody); err != nil {
			return fmt.Errorf(str)
		}
	default:
		return fmt.Errorf(str)
	}
	return nil
}

func ValidateClient(line string) error {
	if !clientNameRegex.MatchString(line) {
		return fmt.Errorf(line)
	}
	return nil
}

func ValidateClientSitDownMiAmor(line string) error {
	parts := strings.Fields(line)

	if len(parts) != 2 {
		return fmt.Errorf(line)
	}

	clientName := parts[0]
	if !clientNameRegex.MatchString(clientName) {
		return fmt.Errorf(line)
	}

	tableNumber, err := strconv.Atoi(parts[1])
	if err != nil || tableNumber <= 0 {
		return fmt.Errorf(line)
	}

	return nil
}
