package client

import (
	"fmt"
	"math"
	"time"
)

// Client structure of client
type Client struct {
	name        string    // Имя клиента
	table       int       // Номер стола, за которым сидит клиент (0 - если не сидит за столом)
	arrivalTime time.Time // Arrival time
	leaveTime   time.Time // Leaving time
	isInClub    bool      // Check for being in club
}

// NewClient makes new instance for client
func NewClient(name string) *Client {
	return &Client{
		name:        name,
		table:       0,
		arrivalTime: time.Now(),
		leaveTime:   time.Now(),
	}
}

// Name readonly method
func (c *Client) Name() string {
	return c.name
}

// Table readonly method
func (c *Client) Table() int {
	return c.table
}

// ArrivalTime readonly method
func (c *Client) ArrivalTime() time.Time {
	return c.arrivalTime
}

// IsInClub readonly method
func (c *Client) IsInClub() bool {
	return c.isInClub
}

// LeaveTime readonly method
func (c *Client) LeaveTime() time.Time {
	return c.leaveTime
}

// SetName sets name
func (c *Client) SetName(name string) {
	c.name = name
}

// SetTable sets table
func (c *Client) SetTable(table int) {
	c.table = table
}

// SetArrivalTime sets arrival time
func (c *Client) SetArrivalTime(arrivalTime time.Time) {
	c.arrivalTime = arrivalTime
}

// SetIsInClub sets status
func (c *Client) SetIsInClub(isInClub bool) {
	c.isInClub = isInClub
}

// SetLeaveTime sets arrival time
func (c *Client) SetLeaveTime(leaveTime time.Time) {
	c.leaveTime = leaveTime
}

// IsPlaying checks if client is playing
func (c *Client) IsPlaying() bool {
	if c.table == 0 {
		return false
	}
	return true
}

// SitTable назначает клиенту указанный стол
func (c *Client) SitTable(tableNumber int) {
	c.table = tableNumber
}

// LeaveClub frees table at client and calculates its total sum per day
func (c *Client) LeaveClub(leaveClubTime time.Time, hourlyRate int) (int, time.Time) {
	c.SetTable(0)
	if leaveClubTime.Before(c.arrivalTime) {
		return 0, time.Now()
	}
	c.SetLeaveTime(leaveClubTime)
	c.SetIsInClub(false)
	result, duration := c.CalculateTotalSum(hourlyRate)
	return result, duration
}

// CalculateTotalSum calculates total sum
func (c *Client) CalculateTotalSum(hourlyRate int) (int, time.Time) {
	duration := c.leaveTime.Sub(c.arrivalTime)

	hours := math.Ceil(duration.Hours())

	parsedTime := CalculateSpentTime(duration)
	return int(hours) * hourlyRate, parsedTime
}

// CalculateSpentTime calculates total time which spent for one table
func CalculateSpentTime(duration time.Duration) time.Time {
	hour := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	spentTime := fmt.Sprintf("%02d:%02d", hour, minutes)
	parsedTime, _ := time.Parse("15:04", spentTime)
	return parsedTime
}

// NewDay prepares client to new day
func (c *Client) NewDay() {
	c.SetTable(0)
	c.SetArrivalTime(time.Now())
	c.SetIsInClub(false)
	c.SetLeaveTime(time.Now())
}
