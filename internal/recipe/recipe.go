package recipe

type Identity struct {
	Name        string `json:"name"`
	Servings    int    `json:"servings"`
	Description string `json:"description"`
}

type Ingredient struct {
	Name         string        `json:"name"`
	Quantity     float64       `json:"quantity"`
	Unit         string        `json:"unit"`
	Note         string        `json:"note,omitempty"`
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Note     string  `json:"note,omitempty"`
}

type Recipe struct {
	ID          int          `json:"id"`
	Lang        string       `json:"lang"`
	Identity    Identity     `json:"identity"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Notes       []string     `json:"notes"`
	Tags        []string     `json:"tags"`
}
