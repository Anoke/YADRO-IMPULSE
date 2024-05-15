package eventHandler

import (
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/club"
	"github.com/Anoke/YADRO-IMPULSE/internal/validation"
	"strconv"
	"strings"
)

// HandleEvent handles event for each ID
func HandleEvent(eventTime string, eventId string, eventBody string, club *club.ComputerClub) {
	// Adds to buffer event itself
	event := eventTime + " " + eventId + " " + eventBody
	club.AddEventToBuffer(event)
	evTime, _ := validation.ParseTimeFormat(eventTime)
	var err error
	// For each event makes its tasks
	switch eventId {
	// Event 1 -- client enters club
	case "1":
		user, ok := club.FindClient(eventBody)
		if !ok {
			user = club.CreateClient(eventBody)
		}
		club.AddClient(user, evTime)
	// Event 2 -- client sits to table
	case "2":
		parts := strings.Fields(eventBody)
		user, ok := club.FindClient(parts[0])
		if !ok {
			// Generated outputEvent
			club.HandleOutputEventId13(evTime, fmt.Errorf("ClientUnknown"))
		}
		tableNumber, _ := strconv.Atoi(parts[1])
		club.AssignTable(user, tableNumber, evTime)
	// Client who was in club enters queue
	case "3":
		err = club.EnqueueClient(eventBody, evTime)
		if err != nil {
			// Handles output event
			club.HandleOutputEventId13(evTime, err)
		}
	// Client leaves club
	case "4":
		user, ok := club.FindClient(eventBody)
		if !ok {
			// Handles output event
			club.HandleOutputEventId13(evTime, fmt.Errorf("ClientUnknown"))
		}
		err = club.ClientLeaves(user, evTime)
		if err != nil {
			club.HandleOutputEventId13(evTime, err)
		}
	}
}
