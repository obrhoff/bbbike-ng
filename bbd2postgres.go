/**
 * Created by IntelliJ IDEA.
 * User: DocterD
 * Date: 28/12/13
 * Time: 11:19
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"./bbbikeng"
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var dataPathFlag = flag.String("path", "", "bbbike data path")
var db *sql.DB

const untitled = "untitled path"

const coordinateRegex = "[0-9]+,[0-9]+"
const nameRegex = "^(.*)(\t)"
const typeRegex = "\t+(.*?)\\s+"

func readLines(path string, fileName string) ([]bbbike.Street, error) {

	file, err := os.Open(path + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var streets []bbbike.Street
	scanner := bufio.NewScanner(file)

	nameRegex := regexp.MustCompile(nameRegex)
	typeRegex := regexp.MustCompile(typeRegex)
	coordsRegex := regexp.MustCompile(coordinateRegex)

	for scanner.Scan() {

		var newStreet bbbike.Street

		infoLine := scanner.Text()
		infoLineConverted := bbbike.ConvertLatinToUTF8([]byte(infoLine))

		name := nameRegex.FindString(infoLineConverted)
		streetType := typeRegex.FindString(infoLineConverted)
		coords := coordsRegex.FindAllString(infoLineConverted, -1)

		if len(coords) > 0 {

			if name == "" {
				name = untitled
			}

			newStreet.Name = strings.TrimSpace(name)
			newStreet.StreetType = strings.TrimSpace(streetType)

			for i, coord := range coords {
				splittedCoords := strings.Split(coord, ",")

				xPath, err := strconv.ParseFloat(splittedCoords[1], 64)
				yPath, err := strconv.ParseFloat(splittedCoords[0], 64)
				if err != nil {
					panic(err)
				}
				var point bbbike.Point

				lat, lng := bbbike.ConvertStandardToWGS84(yPath, xPath)
				point.Lat = lat
				point.Lng = lng
				newStreet.Path = append(newStreet.Path, point)

			}

			fmt.Println("New GeoJSON:", newGeoJSON)

			streets = append(streets, newStreet)
		}

	}

	return streets, scanner.Err()
}

func parseData(path string) {

	fmt.Println("Parsing Pathdata.")
	streets, fileErr := readLines(path, "strassen")
	cyclepaths, fileErr := readLines(path, "radwege")

	if fileErr != nil {
		log.Fatalf("Failed reading Strassen File: %s", fileErr)
	}

	db = bbbike.ConnectToDatabase()
	defer db.Close()

	for i, cyclepath := range cyclepaths {
		cyclepath.PathID = i
		bbbike.InsertCyclePathToDatabase(cyclepath, db)
	}

	for i, street := range streets {
		street.PathID = i
		bbbike.InsertStreetToDatabase(street, db)
	}

}

func main() {

	flag.Parse()
	println("Data dir is: ", *dataPathFlag)
	parseData(*dataPathFlag)

}
