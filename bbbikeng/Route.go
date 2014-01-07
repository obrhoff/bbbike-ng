package bbbikeng

import (
	"fmt"
	"log"
)

type Route struct {
}


func CalculateRoute(startPoint Point, endPoint Point) (route Route) {

	log.Println("Start Latitude:", startPoint)
	log.Println("End Latitude", endPoint)

	startStreet := SearchForNearestStreetFromPoint(startPoint)
	endStreet := SearchForNearestStreetFromPoint(endPoint)

	startStreetPoint := IntersectionFromPointToStreet(startStreet, startPoint)
	endStreetPoint := IntersectionFromPointToStreet(endStreet, endPoint)

	fmt.Println("Correction Start Point:", startStreetPoint)
	fmt.Println("Correction End Point:", endStreetPoint)

	return route

}


func Test(){

	street := GetStreetFromId(148)

	firstInter := street.Intersections[1]
	secondInter := street.Intersections[7]

	fmt.Println("Line From:", firstInter.Street.Name)
	fmt.Println("Line To:", secondInter.Street.Name)

	test, distance := PathFromIntersectionToIntersection(firstInter.Coordinate, secondInter.Coordinate, street)

	fmt.Println("Test:", test)
	fmt.Println("Distance", distance)

	cyclepath := GetCyclepathFromStreet(street)
	greenways := GetGreenwaysFromStreet(street)
	quality := GetQualityFromStreet(street)

	fmt.Println("Cyclepath:", cyclepath)
	fmt.Println("Greenways:", greenways)
	fmt.Println("Quality:", quality)



}


