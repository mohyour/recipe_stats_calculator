package app

// Recipe fixture struct
type RecipeData struct {
	Postcode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
}

// RecipeCount to hold unique recipe and count
type RecipeCount struct {
	Recipe string `json:"recipe"`
	Count  int    `json:"count"`
}

// BusiestPostcode to hold busy postcode and count
type BusiestPostcode struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

// PostcodeDeliveryTime holds info about delivery time to postcode
type PostcodeDeliveryTime struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}

// Results holds struct for expected output
type Result struct {
	NumOfRows               int                  `json:",omitempty"`
	UniqueRecipeCount       int                  `json:"unique_recipe_count"`
	CountPerRecipe          []RecipeCount        `json:"count_per_recipe"`
	BusiestPostcode         BusiestPostcode      `json:"busiest_postcode"`
	CountPerPostCodeAndTime PostcodeDeliveryTime `json:"count_per_postcode_and_time"`
	MatchByName             []string             `json:"match_by_name"`
}

// searchOptions for fixture search values
type SearchOptions struct {
	RecipeSearch string
	Postcode     string
	TimeWindow   string
}