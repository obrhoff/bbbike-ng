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
	"encoding/json"
	"os"
	"io/ioutil"
)

const user = "root"
const password = "root"
const host = "127.0.0.1"
const port = "5433"
const database = "bbbikeng"

var Connection *sql.DB

// connectionParameter := fmt.Sprint("user=", user, " password=", password, " host=", host, " port=", port, " dbname=", database)

type Config struct {
	User string
	Password string
	Host string
	Port string
	Database string
}

func InsertStreetToDatabase(street Street) {

	//log.Println("Inserting Streetpath:", street)
	var err error
	//points := preparePointsForDatabase(street.Path)
	points := geoJsonInsert(ConvertPathToGeoJSON(street.Path))
	fmt.Println("Inserting:", street)
	fixedName := strings.Replace(street.Name, "'", "''", -1)

	query := fmt.Sprintf("INSERT INTO streetpaths(streetpathid, name, type, path) VALUES (%s, '%s', '%s', %s)", strconv.Itoa(street.ID), fixedName, street.Type, points)
	fmt.Println("Insert:", query)
	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Street Into Database: %s", err.Error())
	}


}

func InsertCyclePathToDatabase(cyclepath Street) {

	var err error
	//points := preparePointsForDatabase(cyclepath.Path)
	points := geoJsonInsert(ConvertPathToGeoJSON(cyclepath.Path))
	query := fmt.Sprintf("INSERT INTO cyclepaths(cyclepathid, type, path) VALUES (%s, '%s', %s)", strconv.Itoa(cyclepath.ID), cyclepath.Type, points)
	fmt.Println("Insert:", query)

	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Cyclepath Into Database: %s", err.Error())
	}

}

func InsertGreenToDatabase(green Street) {

	var err error
	points := geoJsonInsert(ConvertPathToGeoJSON(green.Path))
	query := fmt.Sprintf("INSERT INTO greenways(greenwaypathid, type, path) VALUES (%s, '%s', %s)", strconv.Itoa(green.ID), green.Type, points)
	fmt.Println("Insert:", query)

	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Cyclepath Into Database: %s", err.Error())
	}

}

func InsertQualityToDatabase(quality Street) {

	var err error
	//points := preparePointsForDatabase(cyclepath.Path)
	points := geoJsonInsert(ConvertPathToGeoJSON(quality.Path))
	query := fmt.Sprintf("INSERT INTO qualitys(qualitypathid, type, path) VALUES (%s, '%s', %s)", strconv.Itoa(quality.ID), quality.Type, points)
	fmt.Println("Insert:", query)

	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting Qualitys Into Database: %s", err.Error())
	}

}


func InsertCityToDatabase(city City) {

	var err error
	//points := preparePointsForDatabase(cyclepath.Path)
	points := geoJsonInsert(ConvertPathToGeoJSON(city.Geometry))
	query := fmt.Sprintf("INSERT INTO city(citypathid, name, bounds) VALUES (%s, '%s', %s)", strconv.Itoa(city.ID), city.Name, points)
	fmt.Println("Insert:", query)

	_, err = Connection.Exec(query)

	if err != nil {
		log.Fatal("Error inserting City Into Database: %s", err.Error())
	}

}


func GetStreetFromId(id int) (street Street) {

	var geometrys string
	err := Connection.QueryRow("select streetpathid, name, type, ST_AsGeoJSON(path) from streetpaths where streetpathid = $1", id).Scan(&street.ID, &street.Name, &street.Type, &geometrys)
	if err != nil {
		log.Fatal("Error on getting Street from ID %s", err.Error())
	}

	street.SetPathFromGeoJSON(geometrys)
	street.SetIntersections()
	return street

}

func GetCyclepathFromId(id int) (cyclepath Street) {

	var geometrys string
	err := Connection.QueryRow("select cyclepathid, type, ST_AsGeoJSON(path) from cyclepaths where cyclepathid = $1", id).Scan(&cyclepath.ID, &cyclepath.Type, &geometrys)
	if err != nil {
		log.Fatal("Error on getting Cyclepath from ID: %s", err.Error())
	}
	cyclepath.SetPathFromGeoJSON(geometrys)
	cyclepath.SetIntersections()
	return cyclepath

}

// returns crossing streets and intersections
func GetStreetIntersections(street *Street) (intersections []Intersection) {

	rows, err := Connection.Query("select s2.streetpathid, s2.name, s2.type, ST_AsGeoJSON(ST_Intersection(s1.path, s2.path)), ST_AsGeoJSON(s2.path) from streetpaths s1, streetpaths s2 where s1.streetpathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.ID)
	if err != nil {
		log.Fatal("Error on getting Intersections: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {

		var newIntersection Intersection

		var geometrys string
		var intersectionCoordinate string

		err := rows.Scan(&newIntersection.Street.ID, &newIntersection.Street.Name, &newIntersection.Street.Type, &intersectionCoordinate,&geometrys)
		if err != nil {
			log.Fatal(err)
		}

		if street.ID != newIntersection.Street.ID {
			newIntersection.Street.SetPathFromGeoJSON(geometrys)
			newIntersection.SetCoordinationFromGeoJSON(intersectionCoordinate)
			intersections = append(intersections, newIntersection)
		}
	}

	return intersections

}

func GetCyclepathFromStreet(street Street) (cyclepaths []Cyclepath) {

	rows, err := Connection.Query("select s2.type, ST_AsGeoJSON(s2.path) from streetpaths s1, cyclepaths s2 where s1.streetpathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.ID)
	if err != nil {
		log.Fatal("Error on getting Cyclepaths: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var newCyclepath Cyclepath
		var geometrys string
		err := rows.Scan(&newCyclepath.Type, &geometrys)
		if err != nil {
			log.Fatal(err)
		}
		newCyclepath.Path = ConvertGeoJSONtoPath(geometrys)

		if len(newCyclepath.Path) > 0 {
			cyclepaths = append(cyclepaths, newCyclepath)
		}

	}

	fmt.Println("Cyclepaths:", cyclepaths)

	return cyclepaths
}

func GetQualityFromStreet(street Street) (qualitys []Quality) {

	rows, err := Connection.Query("select s2.type, ST_AsGeoJSON(s2.path) from streetpaths s1, qualitys s2 where s1.streetpathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.ID)
	if err != nil {
		log.Fatal("Error on getting Qualitys: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var newQuality Quality
		var geometrys string
		err := rows.Scan(&newQuality.Type, &geometrys)
		if err != nil {
			log.Fatal(err)
		}
		newQuality.Path = ConvertGeoJSONtoPath(geometrys)
		qualitys = append(qualitys, newQuality)
	}

	fmt.Println("Quality:", qualitys)

	return qualitys
}

func GetGreenwaysFromStreet(street Street) (greenways []Greenway) {

	rows, err := Connection.Query("select s2.type, ST_AsGeoJSON(s2.path) from streetpaths s1, greenways s2 where s1.streetpathid = $1 AND (ST_Crosses(s2.path, s1.path) OR ST_Intersects(s2.path, s1.path))", street.ID)
	if err != nil {
		log.Fatal("Error on getting Greenways: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var newGreenway Greenway
		var geometrys string
		err := rows.Scan(&newGreenway.Type, &geometrys)
		if err != nil {
			log.Fatal(err)
		}
		newGreenway.Path = ConvertGeoJSONtoPath(geometrys)
		greenways = append(greenways, newGreenway)
	}

	fmt.Println("Greenway:", greenways)

	return greenways
}

func SearchForStreetName(name string) (streets []Street) {

	log.Println("Searching for Streetname:", name)

	rows, err := Connection.Query("select streetpathid, name, type, ST_AsGeoJSON(path) from streetpaths where name ilike $1 LIMIT 10", ("%" + name + "%"))
	if err != nil {
		log.Fatal("Error on opening database connection: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {

		var newStreet Street
		var geometrys string
		err := rows.Scan(&newStreet.ID, &newStreet.Name, &newStreet.Type, &geometrys)
		if err != nil {
			log.Fatal(err)
		}

		newStreet.SetPathFromGeoJSON(geometrys)
		newStreet.SetIntersections()

		streets = append(streets, newStreet)

	}


	return streets
}

func SearchForNearestStreetFromPoint(point Point) (street Street) {

	//SELECT * FROM streetpath ORDER BY ST_Distance(path, ST_GeomFromText('POINT(13.373370 52.551080)', 4326)) LIMIT 1;
	var geometrys string
	latPath, lngPath := point.LatitudeLongitudeAsString()

	makePoint := ("ST_Distance(path, ST_GeomFromText('POINT(" + lngPath + " " + latPath + ")', 4326))")
	query := fmt.Sprintf("SELECT streetpathid, name, type, ST_AsGeoJSON(path)  FROM streetpaths ORDER BY %s LIMIT 1", makePoint)
	err := Connection.QueryRow(query).Scan(&street.ID, &street.Name, &street.Type, &geometrys)

	if err != nil {
		log.Fatal(err)
	}

	street.SetPathFromGeoJSON(geometrys)
	street.SetIntersections()

	return street

}

func SearchForNearestCyclepathFromPoint(point Point) (cyclepath Street) {

	var geometrys string
	latPath, lngPath := point.LatitudeLongitudeAsString()

	makePoint := ("ST_Distance(path, ST_GeomFromText('POINT(" + lngPath + " " + latPath + ")', 4326))")
	query := fmt.Sprintf("SELECT cycleID, type, ST_AsGeoJSON(path) FROM cyclepaths ORDER BY %s LIMIT 1", makePoint)

	err := Connection.QueryRow(query).Scan(&cyclepath.ID, &cyclepath.Type, &geometrys)

	if err != nil {
		log.Fatal(err)
	}

	cyclepath.SetPathFromGeoJSON(geometrys)
	cyclepath.SetIntersections()

	return cyclepath

}


func ConnectToDatabase() {

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var newConfig Config
	json.Unmarshal(file, &newConfig)
	connectionParameter := fmt.Sprint("user=", newConfig.User, " password=", newConfig.Password, " host=", newConfig.Host, " port=", newConfig.Port, " dbname=", newConfig.Database)
	log.Println("Connecting to Database:", connectionParameter)

	var err error
	Connection, err = sql.Open("postgres", connectionParameter)
	err = Connection.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible

	if err != nil {
		log.Panic("Error on opening database connection: %s", err.Error())
	}


}
