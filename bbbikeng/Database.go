/**
 * User: DocterD
 * Date: 28/12/13
 * Time: 11:19
 */

package bbbikeng

import (
	"database/sql"
	"encoding/json"
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

var connection *sql.DB

func InsertStreetToDatabase(street Street) {

	//log.Println("Inserting Streetpath:", street)
	var err error

	if len(street.Path) >= 2 {
		points := preparePointsForDatabase(street.Path)
		fmt.Println("Inserting:", street)
		fixedName := strings.Replace(street.Name, "'", "''", -1)
		query := fmt.Sprintf("INSERT INTO streetpath(pathid, name, type, path) VALUES (%s, '%s', '%s', %s)", strconv.Itoa(street.PathID), fixedName, street.StreetType, points)
		_, err = connection.Exec(query)
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
		_, err = connection.Exec(query)
	}

	if err != nil {
		log.Fatal("Error inserting Cyclepath Into Database: %s", err.Error())
	}
}

func GetStreetFromId(id int) (street Street) {

	var geometrys string
	err := connection.QueryRow("select pathid, name, type, ST_AsGeoJSON(path)  from streetpath where pathid = $1", id).Scan(&street.PathID, &street.Name, &street.StreetType, &geometrys)
	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}
	street.Path = ConvertStreetPathToObject(geometrys)
	return street

}

func GetCyclepathFromId(id int) (cyclepath Street) {

	var geometrys string
	err := connection.QueryRow("select pathid, type, ST_AsGeoJSON(path) from cyclepath where pathid = $1", id).Scan(&cyclepath.PathID, &cyclepath.StreetType, &geometrys)
	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}
	cyclepath.Path = ConvertStreetPathToObject(geometrys)
	return cyclepath

}

// returns crossing streets and intersections
func GetStreetIntersections(street Street) (intersections []Street) {

	// select s2.* from streetpath s1, streetpath s2 where s1.pathid=148 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path));

	rows, err := connection.Query("select s2.pathid, s2.name, s2.type, ST_AsGeoJSON(s2.path) from streetpath s1, streetpath s2 where s1.pathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.PathID)
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
		newStreet.Path = ConvertStreetPathToObject(geometrys)
		intersections = append(intersections, newStreet)
	}

	return intersections

}

func SearchForStreetName(name string) (streets []Street) {

	log.Println("Searching for Streetname:", name)

	rows, err := connection.Query("select pathid, name, type, ST_AsGeoJSON(path) from streetpath where name ilike $1", ("%" + name + "%"))
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

		newStreet.Path = ConvertStreetPathToObject(geometrys)
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
	err := connection.QueryRow(query).Scan(&street.PathID, &street.Name, &street.StreetType, &geometrys)
	street.Path = ConvertStreetPathToObject(geometrys)
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

	err := connection.QueryRow(query).Scan(&cyclepath.PathID, &cyclepath.StreetType, &geometrys)
	cyclepath.Path = ConvertStreetPathToObject(geometrys)

	if err != nil {
		log.Fatal(err)
	}

	return cyclepath

}

/*
func SearchForNearestStreetIntersectionFromPoint(point Point, street Street) (intersection Point) {

} */

func ConvertStreetPathToObject(json string) (path []Point) {

	var coordinates GeoJSON
	err := json.Unmarshal([]byte(json), &coordinates)
	if err != nil {
		log.Fatal(err)
	}

	for _, coord := range coordinates.Coordinates {
		var newPoint Point
		newPoint.Lat = coord[1]
		newPoint.Lng = coord[0]
		path = append(path, newPoint)
	}

	return path

}


func ConvertPointsToGeoJSON(path []Point)(json string) {

	var newJson GeoJSON
	newJson.Type = "LineString"
	for _, point := range path {
		var newCoordinates [2]float64
		newCoordinates[1] = point.Lat
		newCoordinates[0] = point.Lng
		newJson.Coordinates = append(newJson.Coordinates, newCoordinates)
	}

	jsonData, err := json.Marshal(newJson)
	if err != nil {
		log.Fatal("Failed to Convert Path to GeoJSON: %s", err.Error())
	}

	return string(jsonData)
}

func preparePointsForDatabase(points []Point) (preparedPoints string) {

	for i, point := range points {
		latPath := strconv.FormatFloat(point.Lat, 'f', 6, 64)
		lngPath := strconv.FormatFloat(point.Lng, 'f', 6, 64)
		//(-71.060316 48.432044, -71.060316 48.432044)
		newPoint := (lngPath + " " + latPath)
		if i > 0 {
			preparedPoints = (preparedPoints + ",")
		}
		preparedPoints = (preparedPoints + " " + newPoint)

	}

	return ("ST_GeomFromText('LINESTRING(" + preparedPoints + ")', 4326)")
}

func ConnectToDatabase() {

	connectionParameter := fmt.Sprint("user=", user, " password=", password, " host=", host, " port=", port, " dbname=", database)
	fmt.Println("Connecting to Database:", host)

	connection, err := sql.Open("postgres", connectionParameter)
	err = connection.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible

	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}

}
