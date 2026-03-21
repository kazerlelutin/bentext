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

// knownBentoKeys are legacy line prefixes (Préfixe|valeur) still accepted when reading old files.
var knownBentoKeys = map[string]struct{}{
	"Transport": {}, "Réchauffage": {}, "Froid": {}, "Manger": {},
	"Fuites": {}, "Odeur": {}, "Veille": {}, "Tenue": {}, "Notes": {},
}

func isLegacyBentoSection(section string) bool {
	for _, line := range nonEmptyLines(section) {
		key, _, ok := strings.Cut(line, "|")
		if !ok {
			continue
		}
		if _, known := knownBentoKeys[strings.TrimSpace(key)]; known {
			return true
		}
	}
	return false
}

// valueOnlyBentoShape: 4–9 lines, aucun '|' (les libellés sont uniquement dans le JSON).
func valueOnlyBentoShape(section string) bool {
	lines := nonEmptyLines(section)
	n := len(lines)
	if n < 4 || n > 9 {
		return false
	}
	for _, l := range lines {
		if strings.ContainsRune(l, '|') {
			return false
		}
	}
	return true
}

func isBentoBlock(section string) bool {
	if isLegacyBentoSection(section) {
		return true
	}
	return valueOnlyBentoShape(section)
}

// likelyTagsBlock distingue un bloc de tags (mots courts, peu de mots par ligne) des conseils.
func likelyTagsBlock(section string) bool {
	for _, l := range nonEmptyLines(section) {
		if strings.ContainsAny(l, ".!?") {
			return false
		}
		if len(l) > 48 {
			return false
		}
		if len(strings.Fields(l)) >= 4 {
			return false
		}
	}
	return true
}

func assignNotesTagsBento(r *Recipe, tail []string) {
	n := len(tail)
	if n == 0 {
		return
	}

	// Un seul bloc après les étapes : toujours des tags, sauf ancien format Préfixe|valeur seul (sans section tags).
	if n == 1 {
		if isLegacyBentoSection(tail[0]) {
			if b := parseBentoBlock(tail[0]); !b.isEmpty() {
				r.Bento = b
			}
			return
		}
		r.Tags = linesFromBlock(tail[0])
		return
	}

	last := tail[n-1]
	bentoHere := isBentoBlock(last)

	if n >= 3 && bentoHere {
		if b := parseBentoBlock(last); !b.isEmpty() {
			r.Bento = b
		}
		r.Tags = linesFromBlock(tail[n-2])
		r.Notes = joinSectionLines(tail[:n-2])
		return
	}

	if n == 2 && bentoHere {
		if likelyTagsBlock(tail[0]) {
			r.Tags = linesFromBlock(tail[0])
		} else {
			r.Notes = joinSectionLines(tail[:1])
		}
		if b := parseBentoBlock(last); !b.isEmpty() {
			r.Bento = b
		}
		return
	}

	// Pas de bento : dernier bloc = tags, le reste = notes.
	r.Tags = linesFromBlock(last)
	r.Notes = joinSectionLines(tail[:n-1])
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
	lines := nonEmptyLines(section)
	b := &Bento{}
	if len(lines) == 0 {
		return b
	}
	if isLegacyBentoSection(section) {
		for _, line := range lines {
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
	// Valeurs seules, ordre fixe : transport, réchauffage, froid, manger, puis optionnels Fuites… Notes.
	if len(lines) >= 1 {
		b.Transport = lines[0]
	}
	if len(lines) >= 2 {
		b.Reheat = lines[1]
	}
	if len(lines) >= 3 {
		b.Cold = lines[2]
	}
	if len(lines) >= 4 {
		b.Eating = lines[3]
	}
	if len(lines) >= 5 {
		b.Leaks = lines[4]
	}
	if len(lines) >= 6 {
		b.Smell = lines[5]
	}
	if len(lines) >= 7 {
		b.PrepAhead = lines[6]
	}
	if len(lines) >= 8 {
		b.Holding = lines[7]
	}
	if len(lines) >= 9 {
		b.ExtraNotes = lines[8]
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
