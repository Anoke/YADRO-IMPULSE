package club

import (
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/client"
	"github.com/Anoke/YADRO-IMPULSE/internal/eventHandler"
	"sync"
	"time"
)

var singletonInstance *ComputerClub
var once sync.Once

type ComputerClub struct {
	tables         int
	clientsInClub  map[string]*client.Client // Map of clients in club (name -> client entity)
	tableOccupancy map[int]*client.Client    // Map of occupancy of tables (table number -> client entity)
	queue          []*client.Client          // Waiting list
	workingHours   [2]time.Time              // Start time and end time of working day
	hourlyRate     int                       // Cost of one hour in club
	tableMoney     map[int]int               // Revenue from each table
}

// NewComputerClub makes new instance of club
func NewComputerClub(tables int, startTime time.Time, endTime time.Time, hourlyRate int) *ComputerClub {
	once.Do(func() {
		singletonInstance = &ComputerClub{
			tables:         tables,
			clientsInClub:  make(map[string]*client.Client),
			tableOccupancy: make(map[int]*client.Client, tables),
			queue:          make([]*client.Client, 0),
			workingHours:   [2]time.Time{startTime, endTime},
			hourlyRate:     hourlyRate,
		}
	})
	return singletonInstance
}

// AddClient adds new client to club for ID1
func (club *ComputerClub) AddClient(user *client.Client, arrivalTime time.Time) (*client.Client, error) {
	var name string
	name = user.Name()
	user, _ = club.clientsInClub[name]
	if user.IsInClub() {
		return nil, fmt.Errorf("YouShallNotPass")
	}

	if arrivalTime.Before(club.workingHours[0]) || arrivalTime.After(club.workingHours[1]) {
		return nil, fmt.Errorf("NotOpenYet")
	}

	user.SetIsInClub(true)
	return user, nil
}

//// FindOrCreateClient tries to find client or creates new one
//func (club *ComputerClub) FindOrCreateClient(name string) *client.Client {
//	user, ok := club.clientsInClub[name]
//	if ok {
//		return user
//	}
//	return club.CreateClient(user.Name())
//}

// FindClient find client
func (club *ComputerClub) FindClient(name string) (*client.Client, bool) {
	user, ok := club.clientsInClub[name]
	return user, ok
}

// CreateClient creates new user
func (club *ComputerClub) CreateClient(name string) *client.Client {
	return client.NewClient(name)
}

// AssignTable assigns table to client ID2
func (club *ComputerClub) AssignTable(user *client.Client, tableNumber int) error {
	var name string
	name = user.Name()
	var ok bool
	user, ok = club.clientsInClub[name]
	if !ok {
		return fmt.Errorf("ClientUnknown")
	}

	if club.CheckTableAvailability(tableNumber) {
		return fmt.Errorf("PlaceIsBusy")
	}

	var userTable int
	userTable = user.Table()
	if userTable != 0 {
		club.tableOccupancy[userTable] = nil
		user.SetTable(0)
	}

	club.tableOccupancy[tableNumber] = user
	user.SetTable(tableNumber)
	return nil
}

// ClientLeaves client leaves club ID4
func (club *ComputerClub) ClientLeaves(user *client.Client, leaveTime time.Time) error {
	var name string
	name = user.Name()
	var ok bool
	user, ok = club.clientsInClub[name]
	if !ok {
		return fmt.Errorf("ClientUnknown")
	}

	tableNumber := user.Table()
	club.tableOccupancy[tableNumber] = nil
	club.tableMoney[tableNumber] += user.LeaveClub(leaveTime, club.hourlyRate)
	user2, err := club.DequeueClient()
	if err == nil {
		err2 := club.AssignTable(user2, tableNumber)
		if err2 != nil {
			return err2
		}
	}
	eventHandler.HandleOutputEventId11(leaveTime)
	return nil
}

// CheckTableAvailability checks tables availability
func (club *ComputerClub) CheckTableAvailability(tableNumber int) bool {
	return club.tableOccupancy[tableNumber] != nil
}

// EnqueueClient adds client to waiting list D3
func (club *ComputerClub) EnqueueClient(clientName string, time2 time.Time) error {
	user, ok := club.FindClient(clientName)
	if !ok {
		return fmt.Errorf("ClientUnknown")
	}
	for table := range club.tableOccupancy {
		if club.tableOccupancy[table] == nil {
			return fmt.Errorf("ICanWaitNoLonger")
		}
	}
	if len(club.queue) >= club.tables {
		eventHandler.HandleOutputEventId11(time2)
		return nil
	}
	club.queue = append(club.queue, user)

	return nil
}

// DequeueClient drops client from waiting list
func (club *ComputerClub) DequeueClient() (*client.Client, error) {
	if len(club.queue) == 0 {
		return nil, fmt.Errorf("waiting queue is empty")
	}

	user := club.queue[0]
	club.queue = club.queue[1:]
	return user, nil
}

// CalculateTotalSum calculate revenue
func (club *ComputerClub) CalculateTotalSum() int {
	total := 0
	for _, money := range club.tableMoney {
		total += money
	}
	return total
}

// GetMoney get tables and money from club
func (club *ComputerClub) GetMoney() map[int]int {
	return club.tableMoney
}
