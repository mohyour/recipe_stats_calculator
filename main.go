package main

import (
	"flag"
	"fmt"
	"recipe-stats/app"
	"time"
)

func main() {
	// Load command flags with default values
	fileName := flag.String("file", "calculation_fixtures.json", "a string")
	search := flag.Bool("search", false, "Provide custom search inputs")
	postcode := flag.String("postcode", "10120", "Custom postcode to search")
	timeWindow := flag.String("time-window", "9AM - 4PM", "Custom time window to search")
	output := flag.String("output", "result.json", "output file")
	// Parse flag commands
	flag.Parse()

	var recipes []string
	if *search {
		recipes = flag.Args()
	}
	app.InfoLogger.Printf("File: %s", *fileName)
	app.InfoLogger.Printf("Postcode: %s", *postcode)
	app.InfoLogger.Printf("Time window: %s", *timeWindow)
	app.InfoLogger.Printf("Recipe to search: %s", recipes)
	// update search options
	app.SetFixturesSearchOptions(recipes, *postcode, *timeWindow)
	
	// Process recipe fixtures
	start := time.Now()
	app.InfoLogger.Println("Start processing recipe fixtures...")
	result, rows, err := app.ProcessFixtureFile(*fileName, *output)
	if err != nil {
		app.ErrorLogger.Printf("Unable to process fixture file input : %v", err)
		return
	}
	end := time.Now()
	app.InfoLogger.Printf("Process complete: %v rows processed in %v\n", rows, end.Sub(start))
	app.InfoLogger.Println("Displaying Json Output...")
	// fmt.Printf prints to stdout
	// Alternatively, I can use -> fmt.Fprintf(os.Stdout, "Output: %v \n", result)
	fmt.Println(string(result))
}
