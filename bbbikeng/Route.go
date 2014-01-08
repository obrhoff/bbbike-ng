package bbbikeng

import (
	"log"
)

type Route struct {

	time int
	distance int
	path []Point

}

func CalculateRoute(startPoint Point, endPoint Point) (route Route) {

	startStreet := SearchForNearestStreetFromPoint(startPoint)
	endStreet := SearchForNearestStreetFromPoint(endPoint)

	startPoint = IntersectionFromPointToStreet(startStreet, startPoint)
	endPoint = IntersectionFromPointToStreet(endStreet, endPoint)

	log.Printf("Start: %f,%f", startPoint.Lat, startPoint.Lng)
	log.Printf("End: %f,%f", endPoint.Lat, endPoint.Lng)

	return route

}

func Test(){

	endPoint:= MakeNewPoint(52.55108,13.37337)
	startPoint := MakeNewPoint(52.483943,13.356135)

	startStreet := SearchForNearestStreetFromPoint(startPoint)
	endStreet := SearchForNearestStreetFromPoint(endPoint)

	startPoint = IntersectionFromPointToStreet(startStreet, startPoint)
	endPoint = IntersectionFromPointToStreet(endStreet, endPoint)

	log.Printf("Start: %f,%f", startPoint.Lat, startPoint.Lng)
	log.Printf("End: %f,%f", endPoint.Lat, endPoint.Lng)

	log.Println("Intersection first", startStreet.Intersections)

	bla := PathFromPointToIntersections(startPoint, startStreet)

	for _, bl := range bla {
		log.Println("Score:", bl.distance)
		log.Println("Path:", bl.path)

	}
	log.Println("Bla", bla)
}

