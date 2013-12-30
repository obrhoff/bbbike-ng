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
	"./misc"
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
		infoLineConverted := toUtf8([]byte(infoLine))

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
			xPath, err := strconv.ParseFloat(splittedCoords[0], 64)
			yPath, err := strconv.ParseFloat(splittedCoords[1], 64)
			if err != nil {
				panic(err)
			}
			var point bbbike.Point
			lat, lng := convertToWGS84(yPath, xPath)
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

func convertToWGS84(x float64, y float64) (xLat float64, yLat float64) {
	return x, y
}

func InsertStreetIntoDatabase(streets []bbbike.Street) {

	for _, street := range streets {

		var points string
		for _, pathPart := range street.Path {
			latPath := strconv.FormatFloat(pathPart.Lat, 'f', 1, 64)
			lngPath := strconv.FormatFloat(pathPart.Lng, 'f', 1, 64)
			point := ("(" + latPath + "," + lngPath + ")")
			if points != "" {
				points = (points + "," + point)
			} else {
				points = point
			}
		}

		if points != "" {
			fmt.Println("Inserting:", street.Name, "(", street.StreetType, " )")
			_, err := db.Exec("INSERT INTO public.streets(name, type, streetpath) VALUES ($1, $2, path($3))", street.Name, street.StreetType, points)
			if err != nil {
				log.Fatalf("Database Error - : %s", err.Error())
			}
		}
	}

}

func toUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func main() {

	flag.Parse()
	println("Data dir is: ", *dataPathFlag)
	parseData(*dataPathFlag)

}
