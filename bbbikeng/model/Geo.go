package model

import (
	"math"
)

const RADIUS = 6368500.0

type Point struct {
	Lat float64
	Lng float64
}

func DistanceFromPointToPoint(firstPoint Point, secondPoint Point) (meters int) {

	firstLatitude := firstPoint.Lat * math.Pi / 180
	firstLongitude := firstPoint.Lng * math.Pi / 180

	secondLatitude := secondPoint.Lat * math.Pi / 180
	secondLongitude := secondPoint.Lng * math.Pi / 180

	return int(math.Acos(math.Sin(firstLatitude)*math.Sin(secondLatitude)+math.Cos(firstLatitude)*math.Cos(secondLatitude)*math.Cos((secondLongitude-firstLongitude))) * RADIUS)

}
