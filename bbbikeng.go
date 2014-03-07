package main

import (
	"./bbbikeng"
	"./Import"
	"fmt"
	"flag"
	"os"
	"github.com/ant0ine/go-json-rest"
	"log"
	"net/http"
	"strings"
	"strconv"
)

var startFlag = flag.Bool("run",true, "start bbbikeng")
var dataImportFlag = flag.Bool("import-data", false, "import bbbike path")
var dataImportPathFlag = flag.String("import-path", "./bbbike/data", "bbbike data path")

func main() {

	flag.Parse()
	bbbikeng.ConnectToDatabase()
	defer bbbikeng.Connection.Close()

	if *startFlag && !*dataImportFlag {
		StartBBBikeServer()
	} else if *dataImportFlag{
		StartParsingBBBikeData(*dataImportPathFlag)
	} else {
		fmt.Printf("--import-data --import-path=/bbbike/data\n")
		os.Exit(1)
	}

}

func StartBBBikeServer() {

	log.Println("Starting server...")

	handler := rest.ResourceHandler{}
		handler.EnableGzip = true
		handler.DisableJsonIndent = true
		handler.EnableRelaxedContentType = false
		handler.EnableStatusService = true
		handler.SetRoutes(
		rest.Route{"GET", "/route?:", Route},
		rest.Route{"GET", "/search?:", Search},
	)

	http.ListenAndServe(":8080", &handler)

}

func StartParsingBBBikeData(path string) {

	Import.ParseData(path)

}

func Search(w *rest.ResponseWriter, req *rest.Request) {

	parameters := req.URL.Query()
	search, okSearch := parameters["name"]
	if okSearch && len(search) > 0 {
		//	results := bbbikeng.SearchForStreetName(search[0])
		//fmt.Println("Results:", results)
		//w.WriteJson(&results)
	}
}


func Route(w *rest.ResponseWriter, req *rest.Request) {

	parameters := req.URL.Query()
	start, okStart := parameters["start"]
	end, okEnd := parameters["end"]
	format, okFormat := parameters["format"]

	quality, okQuality := parameters["quality"]
	types, okTypes := parameters["types"]
	greenways, okGreen := parameters["green"]
	unlit, okUnlit := parameters["unlit"]
	trafficLight, okTrafficLight := parameters["lights"]
	speed, okSpeed := parameters["speed"]
	ferries, okFerries := parameters["ferries"]
	performance, okPerformance := parameters["performance"]

	if !okStart || !okEnd {
		return
	}

	splittedStart := strings.Split(start[0], ",")
	splittedEnd := strings.Split(end[0], ",")

	var err error
	var startLat, startLng, endLat, endLng float64

	startLat, err = strconv.ParseFloat(splittedStart[0], 64)
	startLng, err = strconv.ParseFloat(splittedStart[1], 64)
	endLat, err = strconv.ParseFloat(splittedEnd[0], 64)
	endLng, err = strconv.ParseFloat(splittedEnd[1], 64)

	if err != nil {
		return
	}

	startPoint := bbbikeng.MakeNewPoint(startLat, startLng)
	endPoint := bbbikeng.MakeNewPoint(endLat, endLng)
	var route bbbikeng.Route

	preferences := bbbikeng.Preferences{Speed: 20.0,
										Quality: "Q2",
										Types: "N0",
										Greenways: "GR0",
										AvoidUnlit: false,
										AvoidLight: false,
										IncludeFerries: true}

	if okQuality {
		preferences.SetPreferedQuality(quality[0])
	}

	if okTypes {
		preferences.SetPreferedTypes(types[0])
	}

	if okGreen {
		preferences.SetPreferedGreen(greenways[0])
	}

	if okSpeed {
		PreferedSpeed, speedParseError := strconv.ParseInt(speed[0], 0, 64)
		if (speedParseError == nil) {
			preferences.SetPreferedSpeed(PreferedSpeed)
		}
	}

	if okUnlit {
		AvoidUnlit, unlitParseError := strconv.ParseBool(unlit[0])
		if (unlitParseError == nil) {
			preferences.SetAvoidUnlit(AvoidUnlit)
		}
	}


	if okTrafficLight {
		AvoidTrafficLight, trafficLightParseError := strconv.ParseBool(trafficLight[0])
		if (trafficLightParseError == nil) {
			preferences.SetAvoidTrafficLight(AvoidTrafficLight)
		}
	}

	if okFerries {
		IncludeFerries, ferriesParseError := strconv.ParseBool(ferries[0])
		if (ferriesParseError == nil) {
			preferences.SetIncludeFerries(IncludeFerries)
		}

	}

	route.Preferences = preferences
	log.Printf("Start Routing from: %f,%f to %f,%f", startPoint.Lat, startPoint.Lng, endPoint.Lat, endPoint.Lng)
	log.Printf("Preferences:", route.Preferences)
	if okPerformance {
		useBidirectional, _ := strconv.ParseBool(performance[0])
		if useBidirectional {
			route.StartBiRouting(startPoint, endPoint)
		} else {
			route.StartRouting(startPoint, endPoint)
		}
	} else {
		route.StartRouting(startPoint, endPoint)
	}


	if !okFormat {
		w.WriteJson(route.GetGeojson())
	} else {
		switch format[0] {
			case "geojson":
				w.WriteJson(route.GetGeojson())
			case "bbybike":
				w.WriteJson(route.GetBBJson())
			default:
				w.WriteJson(route.GetGeojson())
		}
	}
}
