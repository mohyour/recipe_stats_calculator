package app

import (
	"encoding/json"
	"fmt"
)

// ProcessFileSync process input fixture files one row after the other
func ProcessFixtureFile(file, output string) ([]byte, int, error) {
	result := Result{}
	var numOfRows int
	r, err := loadFile(file)
	if err != nil {
		ErrorLogger.Printf("ProcessFixtureFile(): Unable to decode file input: %v", err)
		return nil, 0, err
	}
	InfoLogger.Println("ProcessFixtureFile(): File successfully loaded")
	decoder := json.NewDecoder(r)

	// Start of object as first token
	t, err := decoder.Token()
	if err != nil {
		return nil, 0, err
	}
	if t != json.Delim('[') {
		ErrorLogger.Printf("ProcessFixtureFile(): Error parsing json file, %v", err)
		return nil, 0, fmt.Errorf("expected {, got %v", t)
	}

	// Continue for while there are more token in stream
	for decoder.More() {
		// Decode token value.
		var recipe RecipeData
		if err := decoder.Decode(&recipe); err != nil {
			return nil, 0, err
		}

		// process recipes.
		err := ProcessRecipeLine(recipe)
		if err != nil {
			ErrorLogger.Printf("ProcessFixtureFile(): Error processing recipe line: %v", err)
			return nil, 0, err
		}
		// increase processed rows
		numOfRows++
	}

	// build result into struct
	response, err := formatOutputData(result)
	if err != nil {
		ErrorLogger.Printf("ProcessFixtureFile(): Error formatting output: %v", err)
		return nil, 0, err
	}

	// write struct to file
	jsonData, err := writeOutputToFile(response, output)
	if err != nil {
		ErrorLogger.Printf("formatOutputData(): Error writing output to file %s : %v", output, err)
		return nil, 0, err
	}

	return jsonData, numOfRows, nil
}
