package bbbikeng

import (
	"database/sql"
	"fmt"
	"log"
)

type Route struct {
}

var found bool

func CalculateRoute(startPoint Point, endPoint Point, db *sql.DB) (route Route) {

	log.Println("Start Latitude:", startPoint)
	log.Println("End Latitude", endPoint)

	startStreet := SearchForNearestStreetFromPoint(startPoint)
	endStreet := SearchForNearestStreetFromPoint(endPoint)

	startStreetPoint := IntersectionFromPointToStreet(startStreet, startPoint)
	endStreetPoint := IntersectionFromPointToStreet(endStreet, endPoint)

	fmt.Println("Correction Start Point:", startStreetPoint)
	fmt.Println("Correction End Point:", endStreetPoint)

	startIntersections := GetStreetIntersections(startStreet)
	endIntersections := GetStreetIntersections(endStreet)

	fmt.Println("Start Intersections:", startIntersections)
	fmt.Println("End Intersections", endIntersections)

	return route

}

func starSearch(startStreet Street, lastStreet Street, endStreet Street, result chan []Street, db *sql.DB) {

	if lastStreet.PathID == endStreet.PathID {
		return
	}

	nextIntersections := GetStreetIntersections(startStreet)
	for _, nextStreet := range nextIntersections {
		if nextStreet.PathID != startStreet.PathID && nextStreet.PathID != lastStreet.PathID {
			log.Println("From Start:", nextStreet.Name, "To:", lastStreet.Name)
			starSearch(nextStreet, startStreet, endStreet, result, db)
		}
	}

}
