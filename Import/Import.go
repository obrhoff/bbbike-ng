package Import

import (
	"../bbbikeng"
	"bufio"
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

const untitled = "untitled path"

//const coordinateRegex = "[-+]?[0-9]+,[-+]?[0-9]+"
const coordinateRegex = " ([-+\\d]?\\d+),([-+]?\\d+)"
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
	coordsRegex := regexp.MustCompile(` ([-+\d]?\d+),([-+]?\d+)`)

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

				xPath, err := strconv.ParseFloat(strings.Replace(splittedCoords[1], " ", "", -1), 64)
				yPath, err := strconv.ParseFloat(strings.Replace(splittedCoords[0], " ", "", -1), 64)
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

func ParseData(path string) {

	log.Println("Parsing Pathdata.")
//	citys, fileErr := readLines(path, "Berlin")

	streets, fileErr := readLines(path, "strassen")
	places, fileErr := readLines(path, "plaetze")
	for _, place := range places {
		streets = append(streets, place)
	}

	qualitys, fileErr := readLines(path, "qualitaet_s")
	cyclepaths, fileErr := readLines(path, "radwege_exact")
	greens, fileErr := readLines(path, "green")
	lights, fileErr := readLines(path, "ampeln");
	unlits, fileErr := readLines(path, "nolighting");

	if fileErr != nil {
		log.Fatalf("Failed reading Strassen File: %s", fileErr)
	}


	/*
	for _, city := range citys {
		var newCity bbbikeng.City
		newCity.Name = city.Name
		newCity.Geometry = city.Path
		bbbikeng.InsertCityToDatabase(newCity)
	} */


	for _, street := range streets {
		var newStreet bbbikeng.Street
		newStreet.Name = street.Name
		newStreet.Type = street.Type
		newStreet.Path = street.Path
		// some data are incomplete and produce only a point and not a multiline. points are inserted into place table instead of way
		if len(newStreet.Path) > 1 {
			bbbikeng.InsertStreetToDatabase(newStreet)
		} else {
			bbbikeng.InsertPlaceToDatabase(newStreet)
		}
	}

	for i, cyclepath := range cyclepaths {
		var newCyclepath bbbikeng.Street
		newCyclepath.ID = i
		newCyclepath.Name = cyclepath.Name
		newCyclepath.Type = cyclepath.Type
		newCyclepath.Path = cyclepath.Path
		if len(newCyclepath.Path) > 1 {
			bbbikeng.InsertCyclePathToDatabase(newCyclepath)
		}
	}

	for i, green := range greens {
		var newGreen bbbikeng.Street
		newGreen.ID = i
		newGreen.Name = green.Name
		newGreen.Type = green.Type
		newGreen.Path = green.Path
		if len(newGreen.Path) > 1 {
			bbbikeng.InsertGreenToDatabase(newGreen)
		}
	}

	for i, quality := range qualitys {
		var newQuality bbbikeng.Street
		newQuality.ID = i
		newQuality.Name = quality.Name
		newQuality.Type = quality.Type
		newQuality.Path = quality.Path
		if len(newQuality.Path) > 1 {
			bbbikeng.InsertQualityToDatabase(newQuality)
		}
	}

	for _, lights := range lights {
		var newLight bbbikeng.Street
		newLight.Type = lights.Type
		newLight.Path = lights.Path
		bbbikeng.InsertStreetLightToDatabase(newLight);
	}

	for _, unlit := range unlits {
		var newUnlit bbbikeng.Street
		newUnlit.Path = unlit.Path
		bbbikeng.InsertUnlitToDatabase(newUnlit);
	}

}
