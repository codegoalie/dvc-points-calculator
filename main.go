package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
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

var dbBucketName = []byte("2021")

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

	db, err := bolt.Open("../dvc-points.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbBucketName)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			resort := Resort{}
			err := json.Unmarshal(v, &resort)
			if err != nil {
				err = fmt.Errorf("failed to parse resort %s: %w", v, err)
				log.Fatal(err)
			}
			resorts = append(resorts, resort)
		}

		return nil
	})
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
