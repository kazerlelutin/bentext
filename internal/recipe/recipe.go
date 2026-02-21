package recipe

type Identity struct {
	Name        string `json:"name"`
	Servings    int    `json:"servings"`
	Description string `json:"description"`
}

// SpriteCoords holds 32x32 sprite position in the ingredient sprite sheet.
type SpriteCoords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Ingredient struct {
	Name         string        `json:"name"`
	Quantity     float64       `json:"quantity"`
	Unit         string        `json:"unit"`
	Note         string        `json:"note,omitempty"`
	Alternatives []Alternative `json:"alternatives"`
	Icon         *SpriteCoords `json:"icon,omitempty"`
}

type Alternative struct {
	Name     string         `json:"name"`
	Quantity float64        `json:"quantity"`
	Unit     string         `json:"unit"`
	Note     string         `json:"note,omitempty"`
	Icon     *SpriteCoords  `json:"icon,omitempty"`
}

// RecipeImage holds the public URL and dimensions of a recipe image.
type RecipeImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Recipe struct {
	Slug        string       `json:"slug"`
	Lang        string       `json:"lang"`
	Identity    Identity     `json:"identity"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Notes       []string     `json:"notes"`
	Tags        []string     `json:"tags"`
	Image       *RecipeImage `json:"image,omitempty"`
}
