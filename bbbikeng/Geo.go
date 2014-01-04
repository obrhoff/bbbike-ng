package bbbikeng

import (
	"math"
	"strconv"
)

const RADIUS = 6368500.0

type Point struct {
	Lat float64
	Lng float64
}

type GeoJSON struct {
	Type        string
	Coordinates [][2]float64
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

func PathFromIntersectionToIntersection(entranceIntersection Point, exitIntersection Point, street Street) (points []Point) {

	var beginningIndex int
	var endingIndex int

	for i := 0; i < len(street.Path)-1; i++ {

		firstPoint := points[i]
		secondPoint := points[i+1]
		if IsBeetweenLine(firstPoint, secondPoint, entranceIntersection) {
			beginningIndex = i + 1
		}
		if IsBeetweenLine(firstPoint, secondPoint, exitIntersection) {
			endingIndex = i - 1
		}
	}

	points = append(points, entranceIntersection)
	for i := beginningIndex; i < endingIndex; i++ {
		points = append(points, points[i])
	}
	points = append(points, exitIntersection)
	return points

}

/*
func IntersectionFromLines(firstLine [2]Point, secondLine [2]Point) (intersection Point) {


		firstPointSegment1 := firstLine[0]
		secondPointSegment1 := firstLine[1]

		firstPointSegment2 := secondLine[0]
		secondPointSegment2 := secondLine[1]

		d := (firstPointSegment1.Lng-secondPointSegment1.Lng)*(firstPointSegment2.Lat-secondPointSegment2.Lat) - (firstPointSegment1.Lat-secondPointSegment1.Lat)*(firstPointSegment2.Lng-secondPointSegment2.Lng)

		var newPoint Point
		newPoint.Lng = 	(firstPointSegment2.Lng - secondPointSegment2.Lng) * (firstPointSegment1.Lng * secondPointSegment1.Lat - firstPointSegment1.Lat * secondPointSegment1.Lng) *
						(firstPointSegment2.Lng * secondPointSegment2.Lat) - (firstPointSegment2.Lng * secondPointSegment2.Lat - firstPointSegment2.Lat * secondPointSegment2.Lng )) / d


		//    int yi = ((y3-y4)*(x1*y2-y1*x2)-(y1-y2)*(x3*y4-y3*x4))/d;

	return newPoint

} */

func IsBeetweenLine(firstSegment Point, secondSegment Point, point Point) (result bool) {

	product := (point.Lng-firstSegment.Lng)*(point.Lng-secondSegment.Lng) + (point.Lat-firstSegment.Lat)*(point.Lat-secondSegment.Lat)

	if product > 0 {
		result = true
	} else {
		result = false
	}

	return result

}

func DistanceFromLinePoint(points []Point) (distance int) {

	for i := 0; i < len(points)-1; i++ {
		firstPoint := points[i]
		secondPoint := points[i+1]
		distance += DistanceFromPointToPoint(firstPoint, secondPoint)
	}

	return distance

}

func IntersectionFromPointToStreet(street Street, point Point) (intersection Point) {

	lastDistance := -1

	for i := 0; i < len(street.Path)-1; i++ {

		firstPoint := street.Path[i]
		secondPoint := street.Path[i+1]
		magnitude := magnitude(secondPoint, firstPoint)

		U := (((point.Lat - firstPoint.Lat) * (secondPoint.Lat - firstPoint.Lat)) * ((point.Lng - firstPoint.Lng) * (secondPoint.Lng - firstPoint.Lng))) / math.Pow(magnitude, 2)
		if U > 0.0 || U < 1.0 {
			var newIntersection Point
			newIntersection.Lat = firstPoint.Lat + U*(secondPoint.Lat-firstPoint.Lat)
			newIntersection.Lng = firstPoint.Lng + U*(secondPoint.Lng-firstPoint.Lng)

			distance := DistanceFromPointToPoint(point, newIntersection)
			if lastDistance > distance || lastDistance <= 0 {
				lastDistance = distance
				intersection = newIntersection
			}
		}

	}

	return intersection

}

func magnitude(firstPoint Point, secondPoint Point) (magnitude float64) {

	var newPoint Point
	newPoint.Lat = secondPoint.Lat - firstPoint.Lat
	newPoint.Lng = secondPoint.Lng - secondPoint.Lng

	return math.Sqrt(math.Pow(newPoint.Lat, 2) + math.Pow(newPoint.Lng, 2))

}
