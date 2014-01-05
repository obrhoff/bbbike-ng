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


