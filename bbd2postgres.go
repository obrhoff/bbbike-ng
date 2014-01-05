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
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Generic struct {
	ID		int
	Name	string
	Type 	string
	Path	[]bbbikeng.Point
}

// go run bbd2postgres.go --path=/Users/DocterD/Development/bbbikeng/bbbike/data

var dataPathFlag = flag.String("path", "", "bbbike data path")

const untitled = "untitled path"

const coordinateRegex = "[0-9]+,[0-9]+"
const nameRegex = "^(.*)(\t)"
const typeRegex = "\t+(.*?)\\s+"

func readLines(path string, fileName string) ([]Generic, error) {

	file, err := os.Open(path + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nameRegex := regexp.MustCompile(nameRegex)
	typeRegex := regexp.MustCompile(typeRegex)
	coordsRegex := regexp.MustCompile(coordinateRegex)

	var newGenerics []Generic

	for scanner.Scan() {

		var newGeneric Generic

		infoLine := scanner.Text()
		infoLineConverted := bbbikeng.ConvertLatinToUTF8([]byte(infoLine))

		name := nameRegex.FindString(infoLineConverted)
		streetType := typeRegex.FindString(infoLineConverted)
		coords := coordsRegex.FindAllString(infoLineConverted, -1)

		if len(coords) > 0 {

			if name == "" {
				name = untitled
			}

			newGeneric.Name = strings.TrimSpace(name)
			newGeneric.Type = strings.TrimSpace(streetType)

			for _, coord := range coords {
				splittedCoords := strings.Split(coord, ",")

				xPath, err := strconv.ParseFloat(splittedCoords[1], 64)
				yPath, err := strconv.ParseFloat(splittedCoords[0], 64)
				if err != nil {
					panic(err)
				}

				var point bbbikeng.Point
				lat, lng := bbbikeng.ConvertStandardToWGS84(yPath, xPath)
				point.Lat = lat
				point.Lng = lng
				newGeneric.Path = append(newGeneric.Path, point)

			}

			newGenerics = append(newGenerics, newGeneric)
		}

	}

	return newGenerics, scanner.Err()
}

func parseData(path string) {

	fmt.Println("Parsing Pathdata.")
	//citys, fileErr := readLines(path, "Berlin")
	streets, fileErr := readLines(path, "strassen")
	cyclepaths, fileErr := readLines(path, "radwege")

	greens, fileErr := readLines(path, "green")
	qualitys, fileErr := readLines(path, "qualitaet_s")

	if fileErr != nil {
		log.Fatalf("Failed reading Strassen File: %s", fileErr)
	}

	bbbikeng.ConnectToDatabase()
	defer bbbikeng.Connection.Close()

	/*
	for i, city := range citys {
		var newCity bbbikeng.City
		newCity.CityID = i
		newCity.Name = city.Name
		newCity.Border = city.Path
		bbbikeng.InsertCityToDatabase(newCity)
	} */


	for i, cyclepath := range cyclepaths {
		var newCyclepath bbbikeng.Street
		newCyclepath.PathID = i
		newCyclepath.Name = cyclepath.Name
		newCyclepath.StreetType = cyclepath.Type
		newCyclepath.Path = cyclepath.Path
		bbbikeng.InsertCyclePathToDatabase(newCyclepath)
	}

	for i, street := range streets {
		var newStreet bbbikeng.Street
		newStreet.PathID = i
		newStreet.Name = street.Name
		newStreet.StreetType = street.Type
		newStreet.Path = street.Path
		bbbikeng.InsertStreetToDatabase(newStreet)
	}

	for i, green := range greens {
		var newGreen bbbikeng.Street
		newGreen.PathID = i
		newGreen.Name = green.Name
		newGreen.StreetType = green.Type
		newGreen.Path = green.Path
		bbbikeng.InsertGreenToDatabase(newGreen)
	}

	for i, quality := range qualitys {
		var newQuality bbbikeng.Street
		newQuality.PathID = i
		newQuality.Name = quality.Name
		newQuality.StreetType = quality.Type
		newQuality.Path = quality.Path
		bbbikeng.InsertQualityToDatabase(newQuality)
	}

}

func main() {

	flag.Parse()
	println("Data dir is: ", *dataPathFlag)
	parseData(*dataPathFlag)

}
