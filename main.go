package main

import (
	"errors"
	"fmt"
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
}

// ErrPointsNotAvailable is reported when points are not available for the RoomType and date
var ErrPointsNotAvailable = errors.New("Points not available")

var resorts = []Resort{
	{
		Name: "The Villas at Grand Floridian",
		RoomTypes: []RoomType{
			{
				Name:        "Deluxe Studio",
				Description: "(Sleeps up to 5)",
				ViewType:    "Standard",
			},
		},
	},
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
				// fmt.Println(resort, roomType, curDate, nextPoints, runningPoints)
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
	if curDate.After(time.Now().AddDate(0, 1, 0)) {
		return -1, ErrPointsNotAvailable
	}

	return 19, nil
}
