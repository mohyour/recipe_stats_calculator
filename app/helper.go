package app

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var uniqueRecipes = map[string]int{}
var uniquePostCode = map[string]int{}
var postcodeDeliveryCount int
var matchByName = []string{}
var nameMatchMap = map[string]int{}

// Default search words if no search word(s) is provided
var defaultSearchWords = []string{"Potato", "Veggie", "Mushroom"}

// Fixture search options
var fixtureSearchOptions = SearchOptions{
	RecipeSearch: strings.Join(defaultSearchWords[:], "|"),
}

// SetFixturesSearchOptions updates fixtures search options with postcode, time window
// and recipes to be searched for.
func SetFixturesSearchOptions(searchWords []string, postcode, timeWindow string) {
	InfoLogger.Printf("SetFixturesSearchOptions(): Updating fixture search options")
	if len(searchWords) > 0 {
		fixtureSearchOptions.RecipeSearch = strings.Join(searchWords[:], "|")
	}
	fixtureSearchOptions.Postcode = postcode
	fixtureSearchOptions.TimeWindow = timeWindow
	InfoLogger.Printf("SetFixturesSearchOptions(): fixture search options - %v", fixtureSearchOptions)
}

// loadFile opens file path to fixtures
func loadFile(filepath string) (*os.File, error) {
	InfoLogger.Printf("loadFile(): Loading file - %s", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		ErrorLogger.Printf("loadFile(): Error loading file %s: %v", filepath, err)
		return nil, err
	}
	return file, nil
}

// processRecipeLine process each line of fixture and applies appropriate calculations
func ProcessRecipeLine(recipe RecipeData) error {
	InfoLogger.Printf("processRecipeLine(): Processing %v", recipe)
	calculateUniqueRecipes(recipe.Recipe)
	calculateUniquePostcodes(recipe.Postcode)
	err := calculateDeliveryCount(recipe.Delivery, recipe.Postcode)
	if err != nil {
		ErrorLogger.Printf("processRecipeLine(): Error converting processing recipe line: %v", err)
		return err
	}
	matchRecipeName(recipe.Recipe)
	return nil
}

// calculateUniqueRecipes calculates unique recipes and keeps counts of their occurrences
func calculateUniqueRecipes(recipe string) {
	if recp, ok := uniqueRecipes[recipe]; ok {
		uniqueRecipes[recipe] = recp + 1
	} else {
		uniqueRecipes[recipe] = 1
	}
}

// calculateUniquePostcodes calculates and keeps count of unique postcodes in fixtures
func calculateUniquePostcodes(postcode string) {
	if pc, ok := uniquePostCode[postcode]; ok {
		uniquePostCode[postcode] = pc + 1
	} else {
		uniquePostCode[postcode] = 1
	}
}

// calculateDeliveryCount calculates delivery time window for postcode
// Defaults to 10120 fro postcode and 9AM - 4PM for delivery
func calculateDeliveryCount(delivery string, postcode string) error {
	// process delivery start time and end time for each recipe delivery
	deliverySlice := strings.Split(delivery, " ")
	am := regexp.MustCompile(`am|AM|pm|PM`)
	lowerTimeLimit := am.Split(deliverySlice[1], 2)[0]
	upperTimeLimit := am.Split(deliverySlice[3], 2)[0]
	deliveryStartTime, err := strconv.Atoi(lowerTimeLimit)
	if err != nil {
		ErrorLogger.Printf("calculateDeliveryCount(): Error converting to integer: %v", err)
		return err
	}
	deliveryEndTime, err := strconv.Atoi(upperTimeLimit)
	if err != nil {
		ErrorLogger.Printf("calculateDeliveryCount(): Error converting to integer: %v", err)
		return err
	}

	// process delivery start time and end time for each recipe search option
	searchSlice := strings.Split(fixtureSearchOptions.TimeWindow, "-")
	searchLowerTimeLimit := strings.TrimSpace(am.Split(searchSlice[0], 2)[0])
	searchUpperTimeLimit := strings.TrimSpace(am.Split(searchSlice[1], 2)[0])
	searchStartTime, err := strconv.Atoi(searchLowerTimeLimit)
	if err != nil {
		ErrorLogger.Printf("calculateDeliveryCount(): Error converting to integer: %v", err)
		return err
	}
	searchEndTime, err := strconv.Atoi(searchUpperTimeLimit)
	if err != nil {
		ErrorLogger.Printf("calculateDeliveryCount(): Error converting to integer: %v", err)
		return err
	}
	if postcode == fixtureSearchOptions.Postcode && (deliveryStartTime >= searchStartTime && deliveryEndTime <= searchEndTime) {
		postcodeDeliveryCount += 1
	}
	return nil
}

// matchRecipeNAme matches recipe names with provided recipe search words or default words
func matchRecipeName(recipe string) {
	var re = regexp.MustCompile(strings.ToLower(fixtureSearchOptions.RecipeSearch))
	if re.MatchString(strings.ToLower(recipe)) {
		if pc, ok := nameMatchMap[recipe]; ok {
			nameMatchMap[recipe] = pc + 1
		} else {
			nameMatchMap[recipe] = 1
		}
	}
}

// formatOutputData formats calculation result into a struct
func formatOutputData(res Result) (Result, error) {
	InfoLogger.Println("formatOutputData(): Formatting output data...")
	res.UniqueRecipeCount = len(uniqueRecipes)
	for key, value := range uniqueRecipes {
		rc := RecipeCount{
			Recipe: key,
			Count:  value,
		}
		res.CountPerRecipe = append(res.CountPerRecipe, rc)
	}

	sort.Slice(res.CountPerRecipe[:], func(i, j int) bool {
		return res.CountPerRecipe[i].Recipe < res.CountPerRecipe[j].Recipe
	})

	var postCodes []BusiestPostcode
	for key, value := range uniquePostCode {
		rc := BusiestPostcode{
			Postcode:      key,
			DeliveryCount: value,
		}
		postCodes = append(postCodes, rc)
	}
	sort.Slice(postCodes[:], func(i, j int) bool {
		return postCodes[i].DeliveryCount > postCodes[j].DeliveryCount
	})

	res.BusiestPostcode = postCodes[0]

	timeWindow := strings.Split(fixtureSearchOptions.TimeWindow, "-")

	res.CountPerPostCodeAndTime = PostcodeDeliveryTime{
		Postcode:      fixtureSearchOptions.Postcode,
		From:          strings.ToUpper(timeWindow[0]),
		To:            strings.ToUpper(timeWindow[1]),
		DeliveryCount: postcodeDeliveryCount,
	}

	for key := range nameMatchMap {
		matchByName = append(matchByName, key)
	}

	res.MatchByName = matchByName

	return res, nil
}

// writeOutputToFile writes calculated result struct to file, defaults to 'result.json' if no output is provided
func writeOutputToFile(result Result, outputFileName string) ([]byte, error) {
	InfoLogger.Printf("writeOutputToFile(): Writing output to %s...", outputFileName)
	jsonData, err := json.MarshalIndent(&result, "", " ")
	if err != nil {
		ErrorLogger.Printf("writeOutputToFile(): Unable to marshal Json : %v", err)
		return nil, err
	}

	// Write to file
	err = os.WriteFile(outputFileName, jsonData, 0644)
	if err != nil {
		ErrorLogger.Printf("writeOutputToFile(): Unable to write result to %s : %v", outputFileName, err)
		return nil, err
	}
	return jsonData, nil
}
