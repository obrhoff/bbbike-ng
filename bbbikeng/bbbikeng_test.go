package bbbikeng

import (
	"database/sql"
	"fmt"
	"testing"
)

var db *sql.DB

func TestSearchForNearestStreetFromPoint() *testing.T {

	db = ConnectToDatabase()
	defer db.Close()

	testPlacePoint1 := MakeNewPoint("52.551080", "13.373370")
	testPlaceResult1 := SearchForNearestStreetFromPoint(testPlacePoint1, db)
	fmt.Println("Results", testPlaceResult1)

}

func TestSearchForStreetName() *testing.T {

	db = ConnectToDatabase()
	defer db.Close()

	SearchForStreetName("urbanstr", db)

}
