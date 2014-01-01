package bbbike

import (
	"math"
	"strconv"
)

const RADIUS = 6368500.0

type Point struct {
	Lat float64
	Lng float64
}

func MakeNewPoint(lat string, lng string) (newPoint Point) {

	xPath, err := strconv.ParseFloat(lat, 64)
	yPath, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		panic(err)
	}

	newPoint.Lat = xPath
	newPoint.Lng = yPath

	return newPoint
}

func PointLatitudeLongitudeAsString(point Point) (lat string, lng string) {

	lat = strconv.FormatFloat(point.Lat, 'f', 6, 64)
	lng = strconv.FormatFloat(point.Lng, 'f', 6, 64)

	return lat, lng

}

func DistanceFromPointToPoint(firstPoint Point, secondPoint Point) (meters int) {

	firstLatitude := firstPoint.Lat * math.Pi / 180
	firstLongitude := firstPoint.Lng * math.Pi / 180

	secondLatitude := secondPoint.Lat * math.Pi / 180
	secondLongitude := secondPoint.Lng * math.Pi / 180

	return int(math.Acos(math.Sin(firstLatitude)*math.Sin(secondLatitude)+math.Cos(firstLatitude)*math.Cos(secondLatitude)*math.Cos((secondLongitude-firstLongitude))) * RADIUS)

}
