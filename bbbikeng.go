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
		handler.SetRoutes(
		rest.Route{"GET", "/search?:", Search},
		rest.Route{"GET", "/route?:", Route},
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
		results := bbbikeng.SearchForStreetName(search[0])
		fmt.Println("Results:", results)
		w.WriteJson(&results)
	}
}

func Route(w *rest.ResponseWriter, req *rest.Request) {

	parameters := req.URL.Query()

	start, okStart := parameters["start"]
	end, okEnd := parameters["end"]

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
	log.Printf("Start Routing from: %f,%f to %f,%f", startPoint.Lat, startPoint.Lng, endPoint.Lat, endPoint.Lng)
	route := bbbikeng.GetAStarRoute(startPoint, endPoint)

	w.WriteJson(route.GetGeojson())

}
