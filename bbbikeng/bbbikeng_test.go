package main

import (
	"./bbbikeng"
	"fmt"
)

func testSearchForNearestStreetFromPoint() {

	db = bbbike.ConnectToDatabase()
	defer db.Close()

	testPlacePoint1 := bbbike.MakeNewPoint("52.551080", "13.373370")
	testPlaceResult1 := bbbike.SearchForNearestStreetFromPoint(testPlacePoint1, db)
	fmt.Println("Results", testPlaceResult1)

}

func testSearchForStreetName() {

	db = bbbike.ConnectToDatabase()
	defer db.Close()

	bbbike.SearchForStreetName("urbanstr", db)

}
