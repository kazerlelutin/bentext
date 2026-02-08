package recipe

import (
	"strconv"
	"strings"
)

const defaultUnit = "piece"

// Parse decodes text content in bentxt format and returns a Recipe.
// The format uses "---" to separate: identity, ingredients, steps, (notes), tags.
func Parse(content string, id int, lang string) *Recipe {
	lines := splitAndTrim(content)
	sections := strings.Split(content, "---")
	if len(sections) < 3 {
		return nil
	}

	r := &Recipe{
		ID:          id,
		Lang:        lang,
		Identity:    Identity{},
		Ingredients: []Ingredient{},
		Steps:       []string{},
		Notes:       []string{},
		Tags:        []string{},
	}

	currentSection := 0
	for _, line := range lines {
		if line == "---" {
			currentSection++
			continue
		}

		switch currentSection {
		case 0:
			parseIdentity(r, line)
		case 1:
			parseIngredientLine(r, line)
		case 2:
			r.Steps = append(r.Steps, line)
		case 3:
			if len(sections) == 5 {
				r.Notes = append(r.Notes, line)
			} else {
				r.Tags = append(r.Tags, line)
			}
		case 4:
			r.Tags = append(r.Tags, line)
		}
	}

	return r
}

func splitAndTrim(content string) []string {
	var out []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func parseIdentity(r *Recipe, line string) {
	if r.Identity.Name == "" {
		r.Identity.Name = line
		return
	}
	if n, err := strconv.Atoi(line); err == nil {
		r.Identity.Servings = n
		return
	}
	r.Identity.Description = line
}

func parseIngredientLine(r *Recipe, line string) {
	parts := strings.Split(line, "~")
	var alternatives []Alternative
	for i, subLine := range parts {
		name, qty, unit, note := parseIngredientPart(subLine)
		if i == 0 {
			r.Ingredients = append(r.Ingredients, Ingredient{
				Name:         name,
				Quantity:     qty,
				Unit:         unit,
				Note:         note,
				Alternatives: []Alternative{},
			})
		} else {
			alternatives = append(alternatives, Alternative{
				Name:     name,
				Quantity: qty,
				Unit:     unit,
				Note:     note,
			})
		}
	}
	if len(alternatives) > 0 && len(r.Ingredients) > 0 {
		r.Ingredients[len(r.Ingredients)-1].Alternatives = alternatives
	}
}

func parseIngredientPart(s string) (name string, quantity float64, unit, note string) {
	unit = defaultUnit
	quantity = 1
	segments := strings.Split(s, "|")
	if len(segments) >= 1 {
		name = strings.TrimSpace(segments[0])
	}
	if len(segments) >= 2 {
		if q, err := strconv.ParseFloat(strings.TrimSpace(segments[1]), 64); err == nil {
			quantity = q
		}
	}
	if len(segments) >= 3 {
		unit = strings.TrimSpace(segments[2])
		if unit == "" {
			unit = defaultUnit
		}
	}
	if len(segments) >= 4 {
		note = strings.TrimSpace(segments[3])
	}
	return name, quantity, unit, note
}
