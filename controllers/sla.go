package controllers

import (
	"SLA/models"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func isValidTime(inputTime time.Time) bool {
	// Define working hours and break time
	workingStartTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 8, 59, 0, 0, inputTime.Location())
	workingEndTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 18, 0, 0, 0, inputTime.Location())
	breakStartTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 12, 0, 0, 0, inputTime.Location())
	breakEndTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 13, 0, 0, 0, inputTime.Location())

	// Check if the input time is a working day (Monday to Friday)
	if inputTime.Weekday() >= time.Monday && inputTime.Weekday() <= time.Friday {
		// Check if the input time is within working hours and outside the break time
		if inputTime.After(workingStartTime) && inputTime.Before(workingEndTime) && !(inputTime.After(breakStartTime) && inputTime.Before(breakEndTime)) {
			return true
		}
	}
	return false
}

func CalculateSLA(c *gin.Context) {
	var i models.SLARequest
	var res models.SLAResponse

	if err := c.ShouldBindQuery(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Wrong input"})
		return
	}

	//Validate Start Time Input
	if !isValidTime(i.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Start time input should Be On Working Hours And Working Day"})
		return
	}

	//Define 50,75,100
	SLA50 := float64(i.SLA) * 0.5
	SLA75 := float64(i.SLA) * 0.75
	SLA100 := float64(i.SLA) * 1

	var wg sync.WaitGroup

	SLA50ch, SLA75ch, SLA100ch := make(chan string), make(chan string), make(chan string)
	// SLA50ch, _, _ := make(chan string), make(chan string), make(chan string)
	go CalculateDate(&wg, i.StartTime, SLA50, SLA50ch)
	go CalculateDate(&wg, i.StartTime, SLA75, SLA75ch)
	go CalculateDate(&wg, i.StartTime, SLA100, SLA100ch)
	//Logic Calculation Goes Here

	fmt.Println(i.StartTime.Weekday())

	res = models.SLAResponse{
		SLA50:  <-SLA50ch,
		SLA75:  <-SLA75ch,
		SLA100: <-SLA100ch,
	}

	c.JSON(200, res)
}

// Expecting The Second Part Is Always 00
func CalculateDate(wg *sync.WaitGroup, time time.Time, sla float64, ch chan string) {
	H, M := float64ToHoursMinutes(sla)

	result := addMinutes(time, M)
	// fmt.Println("add minute", result, M)
	result = addHours(result, H)
	// fmt.Println("add hour", result, H)

	ch <- result.String()
}

func float64ToHoursMinutes(floatValue float64) (int, int) {
	// Extract the integer part as hours
	hours := int(math.Floor(floatValue))

	// Calculate the remaining minutes from the fractional part
	fractionalPart := floatValue - float64(hours)
	minutes := int(math.Round(fractionalPart * 60)) // 60 minutes in an hour

	return hours, minutes
}

func addMinutes(inputTime time.Time, minutesToAdd int) time.Time {
	// fmt.Println("minutesToAdd: ", minutesToAdd)
	// Define working hours and break time
	workingEndTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 17, 59, 59, 0, inputTime.Location())
	breakStartTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 11, 59, 59, 0, inputTime.Location())
	breakEndTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 12, 59, 59, 0, inputTime.Location())

	// Start with the input time
	resultTime := inputTime

	for i := 0; i < minutesToAdd; i++ {
		// Add one hour to the result time
		resultTime = resultTime.Add(time.Minute)
		// fmt.Println(resultTime.String())
		// Check if the result time exceeds the working hours or falls within the break time; if so, adjust accordingly

		if resultTime.After(workingEndTime) {
			resultTime = time.Date(resultTime.Year(), resultTime.Month(), resultTime.Day()+1, 9, 1, 0, 0, inputTime.Location()) // Check if the result day is a weekend (Saturday or Sunday); if so, move to Monday
			for resultTime.Weekday() == time.Saturday || resultTime.Weekday() == time.Sunday {
				resultTime = resultTime.Add(24 * time.Hour)
				//Update the day
				workingEndTime, breakStartTime, breakEndTime = workingEndTime.Add(24*time.Hour), breakStartTime.Add(24*time.Hour), breakEndTime.Add(24*time.Hour)
				continue
			}

			//Update the day
			workingEndTime, breakStartTime, breakEndTime = workingEndTime.Add(24*time.Hour), breakStartTime.Add(24*time.Hour), breakEndTime.Add(24*time.Hour)
			continue
		}
		if resultTime.After(breakStartTime) && resultTime.Before(breakEndTime) {
			// Move to the next working day and reset to the start of working hours
			resultTime = time.Date(resultTime.Year(), resultTime.Month(), resultTime.Day(), 13, 1, 0, 0, inputTime.Location())
		}
	}

	return resultTime
}

func addHours(inputTime time.Time, hoursToAdd int) time.Time {
	fmt.Println("hoursToAdd: ", hoursToAdd)
	// Define working hours and break time
	workingEndTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 17, 59, 59, 0, inputTime.Location())
	breakStartTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 11, 59, 59, 0, inputTime.Location())
	breakEndTime := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 12, 59, 59, 0, inputTime.Location())

	// Start with the input time
	resultTime := inputTime

	// Iterate while adding hours
	for i := 0; i < hoursToAdd; i++ {
		// Add one hour to the result time
		resultTime = resultTime.Add(time.Hour)
		fmt.Println("resultTime: ", resultTime)
		// Check if the result time exceeds the working hours or falls within the break time; if so, adjust accordingly

		if resultTime.After(workingEndTime) {
			resultTime = time.Date(resultTime.Year(), resultTime.Month(), resultTime.Day()+1, 9, 0, 0, 0, inputTime.Location()) // Check if the result day is a weekend (Saturday or Sunday); if so, move to Monday
			for resultTime.Weekday() == time.Saturday || resultTime.Weekday() == time.Sunday {
				resultTime = resultTime.Add(24 * time.Hour)
				//Update the day
				workingEndTime, breakStartTime, breakEndTime = workingEndTime.Add(24*time.Hour), breakStartTime.Add(24*time.Hour), breakEndTime.Add(24*time.Hour)
				continue
			}

			//Update the day
			workingEndTime, breakStartTime, breakEndTime = workingEndTime.Add(24*time.Hour), breakStartTime.Add(24*time.Hour), breakEndTime.Add(24*time.Hour)
			continue
		}
		if resultTime.After(breakStartTime) && resultTime.Before(breakEndTime) {
			// Move to the next working day and reset to the start of working hours
			resultTime = time.Date(resultTime.Year(), resultTime.Month(), resultTime.Day(), 13, resultTime.Minute(), 0, 0, inputTime.Location())
		}
	}

	return resultTime
}
