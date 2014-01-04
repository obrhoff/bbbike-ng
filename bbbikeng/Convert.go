/**
 * User: Dennis Oberhoff
 * To change this template use File | Settings | File Templates.
 */
package bbbikeng

import (
	"encoding/json"
	"log"
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

func ConvertGeoJSONtoPath(json string) (path []Point) {

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

func ConvertPathToGeoJSON(path []Point)(json string) {

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
