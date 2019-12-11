package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

// Resort models a WDW resort
type Resort struct {
	Name      string
	RoomTypes []RoomType
}

// RoomType models a room size and view combination
type RoomType struct {
	Name        string
	Description string
	ViewType    string

	PointChart []PointBlock
}

// PointBlock models the points needed to stay in a RoomType over a range of dates
type PointBlock struct {
	StartDate     time.Time
	EndDate       time.Time
	WeekdayPoints int
	WeekendPoints int
}

// ErrPointsNotAvailable is reported when points are not available for the RoomType and date
var ErrPointsNotAvailable = errors.New("Points not available")

var resorts []Resort
var est *time.Location

func init() {
	var err error
	est, err = time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}

	resorts = []Resort{
		{
			Name: "The Villas at Grand Floridian",
			RoomTypes: []RoomType{
				{
					Name:        "Deluxe Studio",
					Description: "(Sleeps up to 5)",
					ViewType:    "Standard",
					PointChart: []PointBlock{
						// 2019
						{ // Adventure Season 1
							StartDate:     time.Date(2019, time.January, 1, 0, 0, 0, 0, est),
							EndDate:       time.Date(2019, time.January, 31, 23, 59, 59, 0, est),
							WeekendPoints: 20,
							WeekdayPoints: 17,
						},
						{ // Adventure Season 2
							StartDate:     time.Date(2019, time.September, 1, 0, 0, 0, 0, est),
							EndDate:       time.Date(2019, time.September, 30, 23, 59, 59, 0, est),
							WeekendPoints: 20,
							WeekdayPoints: 17,
						},
						{ // Adventure Season 3
							StartDate:     time.Date(2019, time.December, 1, 0, 0, 0, 0, est),
							EndDate:       time.Date(2019, time.December, 14, 23, 59, 59, 0, est),
							WeekendPoints: 20,
							WeekdayPoints: 17,
						},
						{ // Choice Season 3
							StartDate:     time.Date(2019, time.December, 15, 0, 0, 0, 0, est),
							EndDate:       time.Date(2019, time.December, 23, 23, 59, 59, 0, est),
							WeekendPoints: 20,
							WeekdayPoints: 17,
						},
						{ // Premier Season 2
							StartDate:     time.Date(2019, time.December, 24, 0, 0, 0, 0, est),
							EndDate:       time.Date(2019, time.December, 31, 23, 59, 59, 0, est),
							WeekendPoints: 31,
							WeekdayPoints: 36,
						},

						// 2020
						{ // Adventure Season 1
							StartDate:     time.Date(2020, time.January, 1, 0, 0, 0, 0, est),
							EndDate:       time.Date(2020, time.January, 31, 23, 59, 59, 0, est),
							WeekendPoints: 20,
							WeekdayPoints: 17,
						},
					},
				},
			},
		},
	}
}

func main() {
	goalPoints := 120
	// minDays := 4
	startDate := time.Now().AddDate(0, 0, 1)

	for _, resort := range resorts {
		for _, roomType := range resort.RoomTypes {
			curDate := startDate
			tripStart := curDate
			runningPoints := 0

			for {
				nextPoints, err := getNextPoints(roomType, curDate)
				if err != nil {
					if !errors.Is(err, ErrPointsNotAvailable) {
						fmt.Printf("Error getting points: %s", err.Error())
					}
					break
				}

				if runningPoints+nextPoints > goalPoints {
					// report found stay
					fmt.Printf("%s - %s - %s\t%s - %s \t%d nights\t%d\n",
						resort.Name,
						roomType.Name,
						roomType.ViewType,
						tripStart.Format("2006-01-02"),
						curDate.Format("2006-01-02"),
						int(math.Ceil(curDate.Sub(tripStart).Hours()/24.0)),
						runningPoints,
					)
					runningPoints = 0
					tripStart = curDate
					continue
				}

				curDate = curDate.AddDate(0, 0, 1)
				runningPoints += nextPoints
			}
		}
	}
}

func getNextPoints(roomType RoomType, curDate time.Time) (int, error) {
	for _, pointBlock := range roomType.PointChart {
		if curDate.After(pointBlock.StartDate) && curDate.Before(pointBlock.EndDate) {
			dayOfWeek := curDate.Weekday()
			if dayOfWeek == time.Friday || dayOfWeek == time.Saturday {
				return pointBlock.WeekendPoints, nil
			}
			return pointBlock.WeekdayPoints, nil
		}
	}

	return -1, ErrPointsNotAvailable
}
