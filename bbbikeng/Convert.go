/**
 * User: Dennis Oberhoff
 * To change this template use File | Settings | File Templates.
 */
package bbbikeng

import (
	"encoding/json"
	"log"
	"math"
)

const X0 = -780761.760862528
const X1 = 67978.2421158527
const X2 = -2285.59137120724
const Y0 = -5844741.03397902
const Y1 = 1214.24447469596
const Y2 = 111217.945663725

func ConvertStandardToWGS84(x float64, y float64) (xLat float64, yLat float64) {

	yLat = ((x-X0)*Y2 - ((y - Y0) * X2)) / (X1*Y2 - Y1*X2)
	xLat = ((x-X0)*Y1 - (y-Y0)*X1) / (X2*Y1 - X1*Y2)
	return xLat, yLat

}

func ConvertLatinToUTF8(iso8859_1_buf []byte) string {

	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)

}

func ConvertGeoJSONtoPoint(jsonInput string) (point Point) {

	var coordinates GeoJSONPoint

	/*

	err := json.Unmarshal([]byte(jsonInput), &coordinates)
	if err != nil {
		log.Fatal(err)
	} */

	return MakeNewPoint(coordinates.Coordinates[1], coordinates.Coordinates[0])

}

func ConvertGeoJSONtoPath(jsonInput string) (path []Point) {

	var f interface{}
	err := json.Unmarshal([]byte(jsonInput), &f)
	if err != nil {
		log.Fatal("JSON Unmarshal error:", err)
	}

	m := f.(map[string]interface{})
	dataType := m["type"]

	if dataType == "LineString" {
		var coordinates GeoJSON
		err := json.Unmarshal([]byte(jsonInput), &coordinates)
		if err != nil {
			log.Fatal("JSON Unmarshal error:", err)
		}
		for _, coord := range coordinates.Coordinates {
			path = append(path, MakeNewPoint(coord[1], coord[0]))
		}
	} else if dataType == "Point" {

		var coordinates GeoJSONPoint
		err := json.Unmarshal([]byte(jsonInput), &coordinates)
		if err != nil {
			log.Fatal("JSON Unmarshal error:", err)
		}

		point := MakeNewPoint(coordinates.Coordinates[1], coordinates.Coordinates[0])
		path = append(path, point)


	}
	return path
}

func ConvertPathToGeoJSON(path []Point)(jsonOutput string) {

	var jsonData []byte
	var err error

	if len(path) == 1 {
		var newJson GeoJSONPoint
		newJson.Type = "Point"
		newJson.Coordinates[1] = path[0].Lat
		newJson.Coordinates[0] = path[0].Lng
		jsonData, err = json.Marshal(newJson)

	} else {

		var newJson GeoJSON
		newJson.Type = "LineString"
		for _, point := range path {
			var newCoordinates [2]float64
			newCoordinates[1] = point.Lat
			newCoordinates[0] = point.Lng
			newJson.Coordinates = append(newJson.Coordinates, newCoordinates)
		}
		jsonData, err = json.Marshal(newJson)

	}

	if err != nil {
		log.Fatal("Failed to Convert Path to GeoJSON: %s", err.Error())
	}

	return string(jsonData)
}


func geoJsonInsert(geoJson string) (statement string) {

	return ("ST_TRANSFORM(ST_SetSRID(ST_GeomFromGeoJSON('"+ geoJson + "'), '4326'),4326)")

}

func Round(val float64, prec int) float64 {

	var rounder float64
	intermed := val * math.Pow(10, float64(prec))

	if val <= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / math.Pow(10, float64(prec))
}
