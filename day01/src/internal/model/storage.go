package model

type Storage struct {
	Cake []Cake `json:"cake" xml:"cake"`
}

type Cake struct {
	Name        string       `json:"name" xml:"name"`
	Time        string       `json:"time" xml:"stovetime"`
	Ingredients []Ingredient `json:"ingredients" xml:"ingredients>item"`
}

type Ingredient struct {
	Ingredient_name  string `json:"ingredient_name" xml:"itemname"`
	Ingredient_count string `json:"ingredient_count" xml:"itemcount"`
	Ingredient_unit  string `json:"ingredient_unit" xml:"itemunit"`
}
