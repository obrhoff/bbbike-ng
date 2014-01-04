/**
 * User: DocterD
 * Date: 28/12/13
 * Time: 11:19
 */

package bbbikeng

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
)

const user = "root"
const password = "root"
const host = "127.0.0.1"
const port = "5433"
const database = "bbbikeng"

var Connection *sql.DB

func InsertStreetToDatabase(street Street) {

	//log.Println("Inserting Streetpath:", street)
	var err error

	if len(street.Path) >= 2 {
		points := preparePointsForDatabase(street.Path)
		fmt.Println("Inserting:", street)
		fixedName := strings.Replace(street.Name, "'", "''", -1)
		query := fmt.Sprintf("INSERT INTO streetpath(pathid, name, type, path) VALUES (%s, '%s', '%s', %s)", strconv.Itoa(street.PathID), fixedName, street.StreetType, points)
		_, err = Connection.Exec(query)
	}

	if err != nil {
		log.Fatal("Error inserting Street Into Database: %s", err.Error())
	}

}

func InsertCyclePathToDatabase(cyclepath Street) {

	//log.Println("Inserting Cyclepath:", cyclepath)
	var err error
	if len(cyclepath.Path) >= 2 {
		points := preparePointsForDatabase(cyclepath.Path)
		query := fmt.Sprintf("INSERT INTO cyclepath(pathid, type, path) VALUES (%s, '%s', %s)", strconv.Itoa(cyclepath.PathID), cyclepath.StreetType, points)
		_, err = Connection.Exec(query)
	}

	if err != nil {
		log.Fatal("Error inserting Cyclepath Into Database: %s", err.Error())
	}
}

func GetStreetFromId(id int) (street Street) {

	var geometrys string
	err := Connection.QueryRow("select pathid, name, type, ST_AsGeoJSON(path)  from streetpath where pathid = $1", id).Scan(&street.PathID, &street.Name, &street.StreetType, &geometrys)
	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}


	street.Path = ConvertGeoJSONtoPath(geometrys)
	return street

}

func GetCyclepathFromId(id int) (cyclepath Street) {

	var geometrys string
	err := Connection.QueryRow("select pathid, type, ST_AsGeoJSON(path) from cyclepath where pathid = $1", id).Scan(&cyclepath.PathID, &cyclepath.StreetType, &geometrys)
	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}
	cyclepath.Path = ConvertGeoJSONtoPath(geometrys)
	return cyclepath

}

// returns crossing streets and intersections
func GetStreetIntersections(street Street) (intersections []Street) {

	// select s2.* from streetpath s1, streetpath s2 where s1.pathid=148 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path));

	rows, err := Connection.Query("select s2.pathid, s2.name, s2.type, ST_AsGeoJSON(s2.path) from streetpath s1, streetpath s2 where s1.pathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.PathID)
	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var newStreet Street
		var geometrys string
		err := rows.Scan(&newStreet.PathID, &newStreet.Name, &newStreet.StreetType, &geometrys)
		if err != nil {
			log.Fatal(err)
		}
		newStreet.Path = ConvertGeoJSONtoPath(geometrys)
		intersections = append(intersections, newStreet)
	}

	return intersections

}

func SearchForStreetName(name string) (streets []Street) {

	log.Println("Searching for Streetname:", name)

	rows, err := Connection.Query("select pathid, name, type, ST_AsGeoJSON(path) from streetpath where name ilike $1", ("%" + name + "%"))
	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var newStreet Street
		var geometrys string

		err := rows.Scan(&newStreet.PathID, &newStreet.Name, &newStreet.StreetType, &geometrys)
		if err != nil {
			log.Fatal(err)
		}

		newStreet.Path = ConvertGeoJSONtoPath(geometrys)
		streets = append(streets, newStreet)
	}

	return streets
}

func SearchForNearestStreetFromPoint(point Point) (street Street) {

	//SELECT * FROM streetpath ORDER BY ST_Distance(path, ST_GeomFromText('POINT(13.373370 52.551080)', 4326)) LIMIT 1;
	var geometrys string
	latPath, lngPath := PointLatitudeLongitudeAsString(point)

	makePoint := ("ST_Distance(path, ST_GeomFromText('POINT(" + lngPath + " " + latPath + ")', 4326))")
	query := fmt.Sprintf("SELECT pathid, name, type, ST_AsGeoJSON(path)  FROM streetpath ORDER BY %s LIMIT 1", makePoint)
	err := Connection.QueryRow(query).Scan(&street.PathID, &street.Name, &street.StreetType, &geometrys)
	street.Path = ConvertGeoJSONtoPath(geometrys)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Nearest Street:", street)

	return street

}

func SearchForNearestCyclepathFromPoint(point Point) (cyclepath Street) {

	var geometrys string
	latPath, lngPath := PointLatitudeLongitudeAsString(point)
	makePoint := ("ST_Distance(path, ST_GeomFromText('POINT(" + lngPath + " " + latPath + ")', 4326))")
	query := fmt.Sprintf("SELECT pathid, type, ST_AsGeoJSON(path) FROM cyclepath ORDER BY %s LIMIT 1", makePoint)

	err := Connection.QueryRow(query).Scan(&cyclepath.PathID, &cyclepath.StreetType, &geometrys)
	cyclepath.Path = ConvertGeoJSONtoPath(geometrys)

	if err != nil {
		log.Fatal(err)
	}

	return cyclepath

}



func ConnectToDatabase() {

	connectionParameter := fmt.Sprint("user=", user, " password=", password, " host=", host, " port=", port, " dbname=", database)
	fmt.Println("Connecting to Database:", host)


	var err error
	Connection, err = sql.Open("postgres", connectionParameter)
	err = Connection.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible

	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}


}
