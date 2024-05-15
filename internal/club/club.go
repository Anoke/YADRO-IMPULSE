package club

import (
	"fmt"
	"github.com/Anoke/YADRO-IMPULSE/internal/client"
	"github.com/Anoke/YADRO-IMPULSE/internal/output"
	"sort"
	"sync"
	"time"
)

// We could use Singleton package, but external packages are not allowed
var singletonInstance *ComputerClub
var once sync.Once

// Table structure of table
type Table struct {
	Number       int
	Revenue      int
	TimeOccupied time.Time
}

// ComputerClub structure of computer club
type ComputerClub struct {
	tableQuantity  int
	clientsInClub  map[string]*client.Client // Map of clients in club (name -> client entity)
	tableOccupancy map[int]*client.Client    // Map of occupancy of tables (table number -> client entity)
	tables         map[int]*Table            // Massive for tables
	queue          []*client.Client          // Waiting list
	workingHours   [2]time.Time              // Start time and end time of working day
	hourlyRate     int                       // Cost of one hour in club
}

// NewComputerClub singleton for computerClub
func NewComputerClub(tables int, startTime time.Time, endTime time.Time, hourlyRate int) *ComputerClub {
	once.Do(func() {
		singletonInstance = &ComputerClub{
			tableQuantity:  tables,
			tables:         createNewTables(tables),
			clientsInClub:  make(map[string]*client.Client),
			tableOccupancy: createNewTableOccupancy(tables),
			queue:          make([]*client.Client, 0),
			workingHours:   [2]time.Time{startTime, endTime},
			hourlyRate:     hourlyRate,
		}
	})
	return singletonInstance
}

// createNewTables is function for creating new tables for field club.tables
func createNewTables(quantity int) map[int]*Table {
	tables := make(map[int]*Table)

	for i := 1; i <= quantity; i++ {
		table := &Table{
			Number:       i,
			Revenue:      0,
			TimeOccupied: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		tables[i] = table
	}

	return tables
}

// createNewTableOccupancy is function for creating new tables for club.tableOccupancy
func createNewTableOccupancy(quantity int) map[int]*client.Client {
	tables := make(map[int]*client.Client)

	for i := 1; i <= quantity; i++ {
		tables[i] = nil
	}

	return tables
}

// AddClient adds new client to club
func (club *ComputerClub) AddClient(user *client.Client, arrivalTime time.Time) {
	var name string
	name = user.Name()
	user, _ = club.clientsInClub[name]
	if user.IsInClub() {
		club.HandleOutputEventId13(arrivalTime, fmt.Errorf("YouShallNotPass"))
		return
	}

	if arrivalTime.Before(club.workingHours[0]) || arrivalTime.After(club.workingHours[1]) {
		club.HandleOutputEventId13(arrivalTime, fmt.Errorf("NotOpenYet"))
		return
	}

	user.SetIsInClub(true)
}

// FindClient find client
func (club *ComputerClub) FindClient(name string) (*client.Client, bool) {
	user, ok := club.clientsInClub[name]
	return user, ok
}

// CreateClient creates new client
func (club *ComputerClub) CreateClient(name string) *client.Client {
	newClient := client.NewClient(name)
	club.clientsInClub[name] = newClient
	return newClient
}

// AssignTable assigns table to client
func (club *ComputerClub) AssignTable(user *client.Client, tableNumber int, eventTime time.Time) {
	var name string
	name = user.Name()
	var ok bool
	user, ok = club.clientsInClub[name]
	if !ok {
		club.HandleOutputEventId13(eventTime, fmt.Errorf("ClientUnknown"))
		return
	}

	if club.CheckTableAvailability(tableNumber) {
		club.HandleOutputEventId13(eventTime, fmt.Errorf("PlaceIsBusy"))
		return
	}

	var userTable int
	userTable = user.Table()
	if userTable != 0 {
		club.tableOccupancy[userTable] = nil
		//adds to table its duration and revenue for one person
		money, duration := user.LeaveClub(eventTime, club.hourlyRate)
		tyble := club.tables[userTable]
		tyble.Revenue += money
		hour := duration.Hour()
		minute := duration.Minute()
		tyble.TimeOccupied.AddDate(hour, minute, 0)
		user.SetIsInClub(true)

		newUser, err := club.DequeueClient()
		if err == nil {
			club.AssignTable(user, userTable, eventTime)
			club.HandleOutputEventId12(eventTime, newUser)
		}
	}

	club.tableOccupancy[tableNumber] = user
	user.SetTable(tableNumber)
	user.SetArrivalTime(eventTime)
}

// ClientLeaves client leaves club
func (club *ComputerClub) ClientLeaves(user *client.Client, leaveTime time.Time) error {
	var userTable int
	userTable = user.Table()
	if userTable != 0 {
		club.tableOccupancy[userTable] = nil
		money, duration := user.LeaveClub(leaveTime, club.hourlyRate)
		if money == 0 {
			return fmt.Errorf("InvalidFormatTime")
		}
		tyble := club.tables[userTable]
		tyble.Revenue += money
		hour := int(duration.Hour())
		minute := int(duration.Minute())
		tyble.TimeOccupied = tyble.TimeOccupied.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute))

		newUser, err := club.DequeueClient()
		if err == nil {
			club.AssignTable(newUser, userTable, leaveTime)
			club.HandleOutputEventId12(leaveTime, newUser)
		}
	}
	return nil
}

// CheckTableAvailability checks tables availability
func (club *ComputerClub) CheckTableAvailability(tableNumber int) bool {
	return club.tableOccupancy[tableNumber] != nil
}

// EnqueueClient adds client to waiting list
func (club *ComputerClub) EnqueueClient(clientName string, time2 time.Time) error {
	user, ok := club.FindClient(clientName)
	if !ok || !user.IsInClub() {
		return fmt.Errorf("ClientUnknown")
	}

	for table := range club.tableOccupancy {
		if club.tableOccupancy[table] == nil {
			return fmt.Errorf("ICanWaitNoLonger")
		}
	}
	if len(club.queue) >= club.tableQuantity {
		club.HandleOutputEventId11(time2, user)
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

// CalculateTotalSum calculate total revenue. There was mentioned that we need this, but in test it wasn't in test output
func (club *ComputerClub) CalculateTotalSum() int {
	totalRevenue := 0
	for _, table := range club.tables {
		totalRevenue += table.Revenue
	}
	return totalRevenue
}

// EndDay ends every task
func (club *ComputerClub) EndDay() {
	var users []string
	for _, user := range club.tableOccupancy {
		if user != nil {
			err := club.ClientLeaves(user, club.workingHours[1])
			if err != nil {
				return
			}
			users = append(users, user.Name())
		}
	}
	sort.Strings(users)
	for _, i := range users {
		newOldUser, _ := club.FindClient(i)
		club.HandleOutputEventId11(club.workingHours[1], newOldUser)
	}
	output.AddToBuffer(fmt.Sprintf("%s", club.workingHours[1].Format("15:04")))

	// As mentioned earlier (210 line)
	//output.AddToBuffer(fmt.Sprintf("%d", club.CalculateTotalSum()))
	for i := 1; i <= club.tableQuantity; i++ {
		output.AddToBuffer(fmt.Sprintf("%d %d %s", club.tables[i].Number, club.tables[i].Revenue, club.tables[i].TimeOccupied.Format("15:04")))
		club.tables[i].Revenue = 0
		club.tables[i].TimeOccupied = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	}
}

// HandleOutputEventId11 client leaves, also uses in the end of the day
func (club *ComputerClub) HandleOutputEventId11(time time.Time, user *client.Client) {
	output.AddToBuffer(fmt.Sprintf("%s 11 %s", time.Format("15:04"), user.Name()))
}

// HandleOutputEventId12 client sat down
func (club *ComputerClub) HandleOutputEventId12(time time.Time, user *client.Client) {
	output.AddToBuffer(fmt.Sprintf("%s 12 %s %d", time.Format("15:04"), user.Name(), user.Table()))
}

// HandleOutputEventId13 error happened
func (club *ComputerClub) HandleOutputEventId13(time time.Time, err error) {
	output.AddToBuffer(fmt.Sprintf("%s 13 %s", time.Format("15:04"), err))
}

// AddEventToBuffer adds input event to buffer
func (club *ComputerClub) AddEventToBuffer(str string) {
	output.AddToBuffer(str)
}
