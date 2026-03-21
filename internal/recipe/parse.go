package recipe

import (
	"strconv"
	"strings"
)

const defaultUnit = "piece"

// Parse decodes text content in bentxt format and returns a Recipe.
// Sections (séparées par ---) : identité, ingrédients, étapes, [conseils], tags, [bento].
func Parse(content string, slug string, lang string) *Recipe {
	parts := splitIntoSections(content)
	if len(parts) < 3 {
		return nil
	}

	r := &Recipe{
		Slug:        slug,
		Lang:        lang,
		Identity:    Identity{},
		Ingredients: []Ingredient{},
		Steps:       []string{},
		Notes:       []string{},
		Tags:        []string{},
	}

	for _, line := range nonEmptyLines(parts[0]) {
		parseIdentity(r, line)
	}
	for _, line := range nonEmptyLines(parts[1]) {
		parseIngredientLine(r, line)
	}
	for _, line := range nonEmptyLines(parts[2]) {
		r.Steps = append(r.Steps, line)
	}

	if len(parts) > 3 {
		assignNotesTagsBento(r, parts[3:])
	}

	return r
}

func splitIntoSections(content string) []string {
	raw := strings.ReplaceAll(strings.ReplaceAll(content, "\r\n", "\n"), "\r", "\n")
	parts := strings.Split(raw, "---")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	for len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func nonEmptyLines(block string) []string {
	var out []string
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func isBentoSection(section string) bool {
	for _, line := range nonEmptyLines(section) {
		return strings.HasPrefix(line, "Transport|")
	}
	return false
}

func assignNotesTagsBento(r *Recipe, tail []string) {
	n := len(tail)
	if n == 0 {
		return
	}

	if n >= 2 && isBentoSection(tail[n-1]) {
		b := parseBentoBlock(tail[n-1])
		if !b.isEmpty() {
			r.Bento = b
		}
		r.Tags = linesFromBlock(tail[n-2])
		if n >= 3 {
			r.Notes = joinSectionLines(tail[:n-2])
		}
		return
	}

	if n >= 2 {
		r.Tags = linesFromBlock(tail[n-1])
		r.Notes = joinSectionLines(tail[:n-1])
		return
	}

	if isBentoSection(tail[0]) {
		b := parseBentoBlock(tail[0])
		if !b.isEmpty() {
			r.Bento = b
		}
		return
	}
	r.Tags = linesFromBlock(tail[0])
}

func linesFromBlock(block string) []string {
	return nonEmptyLines(block)
}

func joinSectionLines(sections []string) []string {
	var out []string
	for _, s := range sections {
		out = append(out, linesFromBlock(s)...)
	}
	return out
}

func parseBentoBlock(section string) *Bento {
	b := &Bento{}
	for _, line := range nonEmptyLines(section) {
		key, val, ok := strings.Cut(line, "|")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		switch key {
		case "Transport":
			b.Transport = val
		case "Réchauffage":
			b.Reheat = val
		case "Froid":
			b.Cold = val
		case "Manger":
			b.Eating = val
		case "Fuites":
			b.Leaks = val
		case "Odeur":
			b.Smell = val
		case "Veille":
			b.PrepAhead = val
		case "Tenue":
			b.Holding = val
		case "Notes":
			b.ExtraNotes = val
		}
	}
	return b
}

func (b *Bento) isEmpty() bool {
	if b == nil {
		return true
	}
	return b.Transport == "" && b.Reheat == "" && b.Cold == "" && b.Eating == "" &&
		b.Leaks == "" && b.Smell == "" && b.PrepAhead == "" && b.Holding == "" && b.ExtraNotes == ""
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
