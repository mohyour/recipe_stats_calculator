# Recipe Stats Calculator

CLI application that automatically generates JSON file with recipe data.

For details about usage, run `go run main.go -help`

Result is printed to stdout and written to a file name provided. Defaults to (`result.json`) if no file name is provided. The time taken and number of rows processed is also printed to stdout.


Given
-----

Json fixtures file with recipe data.
_Notes on given data_

1. Property value `"delivery"` always has the following format: "{weekday} {h}AM - {h}PM", i.e. "Monday 9AM - 5PM"
2. The number of distinct postcodes is lower than `1M`, one postcode is not longer than `10` chars.
3. The number of distinct recipe names is lower than `2K`, one recipe name is not longer than `100` chars.

Completed features
------

1. Counts the number of unique recipe names.
2. Counts the number of occurrences for each unique recipe name (alphabetically ordered by recipe name).
3. Finds the postcode with most delivered recipes.
4. Count the number of deliveries to postcode `10120` that lie within the delivery time between `10AM` and `3PM`, examples _(`12AM` denotes midnight)_:
    - `NO` - `9AM - 2PM`
    - `YES` - `10AM - 2PM`
5. List the recipe names (alphabetically ordered) that contain in their name one of the following words:
    - Potato
    - Veggie
    - Mushroom


Generated output
---------------

Generates a JSON file of the following format:

```json5
{
    "unique_recipe_count": 15,
    "count_per_recipe": [
        {
            "recipe": "Mediterranean Baked Veggies",
            "count": 1
        },
        {
            "recipe": "Speedy Steak Fajitas",
            "count": 1
        },
        {
            "recipe": "Tex-Mex Tilapia",
            "count": 3
        }
    ],
    "busiest_postcode": {
        "postcode": "10120",
        "delivery_count": 1000
    },
    "count_per_postcode_and_time": {
        "postcode": "10120",
        "from": "11AM",
        "to": "3PM",
        "delivery_count": 500
    },
    "match_by_name": [
        "Mediterranean Baked Veggies", "Speedy Steak Fajitas", "Tex-Mex Tilapia"
    ]
}
```

## Running the application:

Run the application using
    
`go run main.go`

Smaller data is provided in `test.json`

`go run main.go -file=test.json`

The application runs using the following values as default:
* file ->  json file - `calculation_fixtures.json`. Should be placed in the root folder
* postcode -> `10120`
* search words -> `Potato, Veggies, Mushrooms`
* output -> `result.json`
* time-window -> `9AM - 4PM`

### Running with custom values - Usage
To add custom values, use the appropriate flag followed by the value:

For file:

    go run main.go -file=test.json

For postcode:

    go run main.go -postcode=1234

For time window

    go run main.go -time-window=1am-3pm

For output

    go run main.go -output=output.json

For search with recipe name. Use the `-search` flag followed by search word separated by space (the words should be the last input)

    go run main.go -search veggies chicken

These flags can be combined for a more robust search and fixture result

    go run main.go -time-window=1am-3pm -postcode=1234

    go run main.go -file=test.json -time-window=1am-3pm

    go run main.go -file=test.json -time-window=10am-4pm -output=output.json

Although the search flag can be used with any other flag, the recipe search words has to come last. Example shown below

    go run main.go -file=test.json -time-window=1am-3pm -search pork meat

    go run main.go -search -file=test.json -time-window=1am-3pm pork meat

    go run main.go -file=test.json -search -time-window=1am-3pm -output=output.json

The commands below is going to be read as invalid and/or produce wrong results:

    go run main.go -search pork meat -file=test.json -time-window=1am-3pm

    go run main.go -file=test.json -search pork meat -time-window=1am-3pm

Where no flags are indicated, they search values are set to their default values

### Running in docker:
- Build the image 

    `docker build <image-name> .`

- Run the image using:

    `docker run <image-name>`

    This is with assumption that the `hf_test_calculation_fixtures.json` file was downloaded, placed in the root directory and built with the image.

The flag option can be applied as mentioned earlier with same rule.

    docker run <image-name> -file=test.json -search onions fajitas chicken

This is with assumption that the `test.json` file was built with the image.


## Todo: Improvements
- File loading and processing can use goroutine for faster operation time
- Better code organization, including using methods in some place
- Better error handling
- Code refactoring to make it more efficient and easily testable