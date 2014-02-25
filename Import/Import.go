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

func readLines(path string, fileName string) (newData []bbbikeng.AttributeInterface, error error) {

	file, err := os.Open(path + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nameRegex := regexp.MustCompile(nameRegex)
	typeRegex := regexp.MustCompile(typeRegex)
	coordsRegex := regexp.MustCompile(` ([-+\d]?\d+),([-+]?\d+)`)

	for scanner.Scan() {

		var newGeneric bbbikeng.AttributeInterface
		switch fileName {
			case "strassen": {
				newGeneric = new(bbbikeng.WayAttribute)
			}
			case "plaetze": {
				newGeneric = new(bbbikeng.WayAttribute)
			}
			case "qualitaet_s": {
				newGeneric = new(bbbikeng.QualityAttribute)
			}
			case "radwege_exact": {
				newGeneric = new(bbbikeng.CyclepathAttribute)
			}
			case "green": {
				newGeneric = new(bbbikeng.GreenwayAttribute)
			}
			case "nolighting": {
				newGeneric = new(bbbikeng.UnlitAttribute)
			}
			case "ampeln": {
				newGeneric = new(bbbikeng.TrafficLightAttribute)
			}
			case "handicap_s": {
				newGeneric = new(bbbikeng.HandicapAttribute)
			}

		}


		infoLine := scanner.Text()
		infoLineConverted := bbbikeng.ConvertLatinToUTF8([]byte(infoLine))

		name := nameRegex.FindString(infoLineConverted)
		streetType := typeRegex.FindString(infoLineConverted)
		coords := coordsRegex.FindAllString(infoLineConverted, -1)

		if len(coords) > 0 {

			if name == "" {
				name = untitled
			}

			newGeneric.SetName(strings.TrimSpace(name))
			newGeneric.SetType(strings.TrimSpace(streetType))
			var tempPath []bbbikeng.Point

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
				tempPath = append(tempPath, point)
			}
			newGeneric.SetGeometry(tempPath)
			newData = append(newData, newGeneric)
		}

	}

	return newData, scanner.Err()
}

func ParseData(path string) {

	log.Println("Parsing Pathdata.")
	var data []bbbikeng.AttributeInterface

	streets, fileErr := readLines(path, "strassen")
	places, fileErr := readLines(path, "plaetze")
	qualitys, fileErr := readLines(path, "qualitaet_s")
	cyclepaths, fileErr := readLines(path, "radwege_exact")
	green, fileErr := readLines(path, "green")
	lights, fileErr := readLines(path, "ampeln");
	unlit, fileErr := readLines(path, "nolighting");
	handicap, fileErr := readLines(path, "handicap_s");
	for _, newData := range streets {
		data = append(data, newData)
	}
	for _, newData := range places {
		data = append(data, newData)
	}
	for _, newData := range qualitys {
		data = append(data, newData)
	}
	for _, newData := range cyclepaths {
		data = append(data, newData)
	}
	for _, newData := range green {
		data = append(data, newData)
	}
	for _, newData := range lights {
		data = append(data, newData)
	}
	for _, newData := range unlit {
		data = append(data, newData)
	}
	for _, newData := range handicap {
		data = append(data, newData)
	}

	if fileErr != nil {
		log.Fatalf("Failed reading Strassen File: %s", fileErr)
	}

	for _, processedData := range data {
		log.Println("Processed:", processedData.Name(), processedData.Type(), processedData.Geometry())
		switch processedData.(type)  {
			case *bbbikeng.WayAttribute: {
				if len(processedData.Geometry()) > 1 {
					bbbikeng.InsertStreetToDatabase(processedData.(*bbbikeng.WayAttribute))
				} else {
					bbbikeng.InsertPlaceToDatabase(processedData.(*bbbikeng.WayAttribute))
				}
			}
			case *bbbikeng.CyclepathAttribute: {
				bbbikeng.InsertCyclePathToDatabase(processedData.(*bbbikeng.CyclepathAttribute))
			}
			case *bbbikeng.GreenwayAttribute: {
				bbbikeng.InsertGreenToDatabase(processedData.(*bbbikeng.GreenwayAttribute))
			}
			case *bbbikeng.QualityAttribute: {
				bbbikeng.InsertQualityToDatabase(processedData.(*bbbikeng.QualityAttribute))
			}
			case *bbbikeng.UnlitAttribute: {
				bbbikeng.InsertUnlitToDatabase(processedData.(*bbbikeng.UnlitAttribute))
			}
			case *bbbikeng.TrafficLightAttribute: {
				bbbikeng.InsertStreetLightToDatabase(processedData.(*bbbikeng.TrafficLightAttribute))
			}
			case *bbbikeng.HandicapAttribute: {
				bbbikeng.InsertHandicapToDatabase(processedData.(*bbbikeng.HandicapAttribute))
			}

		}
	}

}
