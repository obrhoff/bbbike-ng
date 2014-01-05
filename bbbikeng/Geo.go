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

type GeoJSONPoint struct {
	Type        string
	Coordinates [2]float64
}

func (f *Point) SetLat(lat float64){
	f.Lat = Round(lat, 6)
}

func (f *Point) SetLng(lng float64){
	f.Lng = Round(lng, 6)
}

func (f *Point) Coordinates()(lat float64, lng float64) {
	return f.Lat, f.Lng
}

func (f *Point) Compare(comparePoint Point) (equal bool) {
	thresholdLat := math.Abs(f.Lat) - math.Abs(comparePoint.Lat)
	thresholdLng := math.Abs(f.Lng) - math.Abs(comparePoint.Lng)
	return (thresholdLat <= 0.0000001 && thresholdLng <= 0.0000001)
}

func MakeNewPoint(lat float64, lng float64) (newPoint Point) {

	newPoint.SetLat(lat)
	newPoint.SetLng(lng)
	return newPoint

}

func MakeNewPointFromString(lat string, lng string) (newPoint Point) {

	xPath, err := strconv.ParseFloat(lng, 64)
	yPath, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		panic(err)
	}
	return MakeNewPoint(xPath, yPath)
}

func (f *Point) LatitudeLongitudeAsString() (lat string, lng string) {

	lat = strconv.FormatFloat(f.Lat, 'f', 6, 64)
	lng = strconv.FormatFloat(f.Lng, 'f', 6, 64)

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
		if isBoundedBox(firstPoint, secondPoint, entranceIntersection) {
			beginningIndex = i + 1
		} else if isBoundedBox(firstPoint, secondPoint, exitIntersection) {
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

func BearingBetweenPoints(firstSegment Point, secondSegment Point) (angle float64) {

	_, dLng := pointDifference(firstSegment, secondSegment)

	y := math.Sin(dLng) * math.Cos(secondSegment.Lat)
	x := math.Cos(firstSegment.Lat) * math.Sin(secondSegment.Lat) - math.Sin(firstSegment.Lat) * math.Cos(secondSegment.Lat) * math.Cos(dLng)

	return radiansToDegrees(math.Atan2(y, x))
}

func IntersectionFromLines(firstLine [2]Point, secondLine [2]Point) (intersection Point, valid bool) {

	firstLineFirstPoint := firstLine[0]
	firstLineSecondPoint := firstLine[1]
	secondLineFirstPoint := secondLine[0]
	secondLineSecondPoint := secondLine[1]

	A1 := firstLineSecondPoint.Lat - firstLineFirstPoint.Lat
	B1 := firstLineFirstPoint.Lng - firstLineSecondPoint.Lng
	C1 := A1 * firstLineFirstPoint.Lng + B1 * firstLineFirstPoint.Lat

	A2 := secondLineSecondPoint.Lat - secondLineFirstPoint.Lat
	B2 := secondLineFirstPoint.Lng - secondLineSecondPoint.Lng
	C2 := A2 * secondLineFirstPoint.Lng + B2 * secondLineFirstPoint.Lat

	determinate := A1 * B2 - A2 * B1

	if determinate != 0 {
		intersection.SetLat((A1 * C2 - A2 * C1) / determinate)
		intersection.SetLng((B2 * C1 - B1 * C2) / determinate)
	}

	return intersection, valid

}

func IntersectionFromPaths(firstPath []Point, secondPath []Point) (intersection Point){

	for innerIndex := 0; innerIndex < len(firstPath)-1; innerIndex++ {

		firstPathFirstPoint := firstPath[innerIndex]
		firstPathSecondPoint := firstPath[innerIndex+1]

		for outerIndex := 0; outerIndex < len(secondPath)-1; outerIndex++ {

			secondPathFirstPoint := secondPath[outerIndex]
			secondPathSecondPoint := secondPath[outerIndex+1]

			var firstLine [2]Point
			var secondLine [2]Point

			firstLine[0] = firstPathFirstPoint
			firstLine[1] = firstPathSecondPoint

			secondLine[0] = secondPathFirstPoint
			secondLine[1] = secondPathSecondPoint

			foundIntersection, _ := IntersectionFromLines(firstLine, secondLine)
			if isBoundedBox(firstPathFirstPoint, firstPathSecondPoint, foundIntersection) && isBoundedBox(secondPathFirstPoint, secondPathSecondPoint, foundIntersection) {
				if foundIntersection.Compare(firstPathFirstPoint) {
					return firstPathFirstPoint
				} else if foundIntersection.Compare(firstPathSecondPoint) {
					return firstPathSecondPoint
				} else if foundIntersection.Compare(secondPathFirstPoint) {
					return secondPathFirstPoint
				} else if foundIntersection.Compare(secondPathSecondPoint) {
					return secondPathSecondPoint
				} else {
					return foundIntersection
				}
			}
		}
	}

	return intersection
}

func DistanceFromLinePoint(points []Point) (distance int) {

	for i := 0; i < len(points)-1; i++ {
		firstPoint := points[i]
		secondPoint := points[i+1]
		distance += DistanceFromPointToPoint(firstPoint, secondPoint)
	}

	return distance

}

func DistanceFromPointToPath(point Point, path []Point) (distance int) {

	distance = -1
	for i := 0; i < len(path)-1; i++ {

		firstPoint := path[i]
		secondPoint := path[i+1]
		magnitude := magnitude(secondPoint, firstPoint)

		U := (((point.Lat - firstPoint.Lat) * (secondPoint.Lat - firstPoint.Lat)) * ((point.Lng - firstPoint.Lng) * (secondPoint.Lng - firstPoint.Lng))) / math.Pow(magnitude, 2)
		if U > 0.0 || U < 1.0 {
			var newIntersection Point
			newIntersection.Lat = firstPoint.Lat + U*(secondPoint.Lat-firstPoint.Lat)
			newIntersection.Lng = firstPoint.Lng + U*(secondPoint.Lng-firstPoint.Lng)

			lastDistance := DistanceFromPointToPoint(point, newIntersection)
			if lastDistance > distance || lastDistance <= 0 {
				distance = lastDistance
			}
		}

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

			newIntersection := MakeNewPoint(firstPoint.Lat + U*(secondPoint.Lat-firstPoint.Lat),firstPoint.Lng + U*(secondPoint.Lng-firstPoint.Lng))

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

func isOnLine(point Point, line [2]Point) (onLine bool) {

	endPoint1 := line[0]
	endPoint2 := line[1]

	/*    return (((double)checkPoint.Y - endPoint1.Y)) / ((double)(checkPoint.X - endPoint1.X))
        == ((double)(endPoint2.Y - endPoint1.Y)) / ((double)(endPoint2.X - endPoint1.X));*/


	return ((point.Lat - endPoint1.Lat)) / ((point.Lng - endPoint1.Lng)) == ((point.Lat - endPoint2.Lat)) / ((point.Lng - endPoint2.Lng));
}

func isBoundedBox(firstPoint Point, secondPoint Point, checkPoint Point) (isBounded bool) {

	var upLatitude float64
	var downLatitude float64
	var upLongitude float64
	var downLongitude float64

	if firstPoint.Lat >= secondPoint.Lat {
		upLatitude = firstPoint.Lat
		downLatitude = secondPoint.Lat
	} else {
		upLatitude = secondPoint.Lat
		downLatitude = firstPoint.Lat
	}
	if firstPoint.Lng >= secondPoint.Lng {
		upLongitude = firstPoint.Lng
		downLongitude = secondPoint.Lng
	} else {
		upLongitude = secondPoint.Lng
		downLongitude = firstPoint.Lng
	}

	return (upLatitude - checkPoint.Lat <= upLatitude - downLatitude &&  upLongitude - checkPoint.Lng <= upLongitude - downLongitude)

}

func pointDifference(firstPoint Point, secondPoint Point) (dLat float64, dLon float64) {


	dLat = secondPoint.Lat - firstPoint.Lat
	dLon = secondPoint.Lat - firstPoint.Lat

	return dLat, dLon

}

func degreeToRadians(degree float64) (radians float64){

	return (degree * math.Pi / 180)
}

func radiansToDegrees(radians float64) ( degrees float64) {

	return (radians * 180 / math.Pi)

}
