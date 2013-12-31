/**
 * Created by IntelliJ IDEA.
 * User: DocterD
 * Date: 28/12/13
 * Time: 11:19
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"./bbbikeng/helper"
	"./bbbikeng/model"
	"./misc"
)

var dataPathFlag = flag.String("path", "", "bbbike data path")
var db *sql.DB

const untitled = "untitled path"

const coordinateRegex = "[0-9]+,[0-9]+"
const nameRegex = "^(.*)(\t)"
const typeRegex = "\t+(.*?)\\s+"

func readLines(path string, fileName string) ([]model.Street, error) {

	file, err := os.Open(path + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var streets []model.Street
	scanner := bufio.NewScanner(file)

	nameRegex := regexp.MustCompile(nameRegex)
	typeRegex := regexp.MustCompile(typeRegex)
	coordsRegex := regexp.MustCompile(coordinateRegex)

	for scanner.Scan() {

		var newStreet model.Street

		infoLine := scanner.Text()
		infoLineConverted := helper.ConvertLatinToUTF8([]byte(infoLine))

		name := nameRegex.FindString(infoLineConverted)
		streetType := typeRegex.FindString(infoLineConverted)
		coords := coordsRegex.FindAllString(infoLineConverted, -1)

		if name == "" {
			name = untitled
		}

		newStreet.Name = name
		newStreet.StreetType = streetType

		for _, coord := range coords {
			splittedCoords := strings.Split(coord, ",")
			xPath, err := strconv.ParseFloat(splittedCoords[1], 64)
			yPath, err := strconv.ParseFloat(splittedCoords[0], 64)
			if err != nil {
				panic(err)
			}
			var point model.Point

			lat, lng := helper.ConvertStandardToWGS84(yPath, xPath)
			point.Lat = lat
			point.Lng = lng
			newStreet.Path = append(newStreet.Path, point)
		}
		streets = append(streets, newStreet)

	}
	return streets, scanner.Err()

}

func parseData(path string) {

	fmt.Println("Parsing Pathdata.")
	streets, fileErr := readLines(path, "strassen")
	cyclepath, fileErr := readLines(path, "radwege")

	if fileErr != nil {
		log.Fatalf("Failed reading Strassen File: %s", fileErr)
	}

	db = util.ConnectToDatabase()
	defer db.Close()

	InsertStreetIntoDatabase(streets)
	InsertStreetIntoDatabase(cyclepath)

}

func InsertStreetIntoDatabase(streets []model.Street) {

	for _, street := range streets {

		var points string
		for _, pathPart := range street.Path {
			latPath := strconv.FormatFloat(pathPart.Lat, 'f', 6, 64)
			lngPath := strconv.FormatFloat(pathPart.Lng, 'f', 6, 64)
			point := ("(" + latPath + "," + lngPath + ")")
			if points != "" {
				points = (points + "," + point)
			} else {
				points = point
			}
		}

		if points != "" {
			fmt.Println("Inserting:", street.Name, "(", street.StreetType, " ) - (", points, ")")
			_, err := db.Exec("INSERT INTO public.streets(name, type, streetpath) VALUES ($1, $2, path($3))", street.Name, street.StreetType, points)
			if err != nil {
				log.Fatalf("Database Error - : %s", err.Error())
			}
		}
	}

}

func main() {

	flag.Parse()
	println("Data dir is: ", *dataPathFlag)
	parseData(*dataPathFlag)

}
