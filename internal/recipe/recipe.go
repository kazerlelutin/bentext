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

// Bento holds key|value lines from the optional section after tags (transport, reheat, cold chain, how to eat, etc.).
// JSON field names are English; values stay as in the source file (per language).
type Bento struct {
	Transport  string `json:"transport,omitempty"`
	Reheat     string `json:"reheat,omitempty"`      // Réchauffage
	Cold       string `json:"cold,omitempty"`        // Froid
	Eating     string `json:"eating,omitempty"`      // Manger
	Leaks      string `json:"leaks,omitempty"`       // Fuites
	Smell      string `json:"smell,omitempty"`       // Odeur
	PrepAhead  string `json:"prep_ahead,omitempty"`  // Veille
	Holding    string `json:"holding,omitempty"`     // Tenue
	ExtraNotes string `json:"extra_notes,omitempty"` // Notes (préfixe Notes| dans la section bento)
}

type Recipe struct {
	Slug        string       `json:"slug"`
	Lang        string       `json:"lang"`
	Identity    Identity     `json:"identity"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Notes       []string     `json:"notes"`
	Tags        []string     `json:"tags"`
	Bento       *Bento       `json:"bento,omitempty"`
	Image       *RecipeImage `json:"image,omitempty"`
	// Bentext is the raw .bentext file contents; set only when the API client asks for it (omitted when empty).
	Bentext string `json:"bentext,omitempty"`
}
