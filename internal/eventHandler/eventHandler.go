package eventHandler

import (
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/club"
	"github.com/Anoke/YADRO-IMPULSE/internal/validation"
	"strconv"
	"strings"
	"time"
)

func HandleEvent(eventTime string, eventId string, eventBody string, club *club.ComputerClub) error {
	time, _ := validation.ParseTimeFormat(eventTime)
	var err error
	switch eventId {
	case "1":
		user, ok := club.FindClient(eventBody)
		if !ok {
			user = club.CreateClient(eventBody)
		}
		user, err = club.AddClient(user, time)
		if err != nil {
			return err
		}
	case "2":
		parts := strings.Fields(eventBody)
		user, ok := club.FindClient(parts[0])
		if !ok {
			return fmt.Errorf("ClientUnknown")
		}
		tableNumber, err2 := strconv.Atoi(parts[1])
		if err2 != nil {
			return err2
		}
		err = club.AssignTable(user, tableNumber)
		if err != nil {
			return err
		}
	case "3":
		err = club.EnqueueClient(eventBody)
		if err != nil {
			return err
		}
	case "4":
	}
}

// HandleOutputEventId11 client leaves, also uses in the end of the day
func HandleOutputEventId11(time2 time.Time) {

}
func HandleOutputEventId12(time time.Time, err error, tableNumber int) {

}
func HandleOutputEventId13(time time.Time, err error) {

}
