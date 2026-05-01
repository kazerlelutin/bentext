package ingredients

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	spriteSize = 32
	// maxLevenshtein is used only after exact match and leading-word fallback fail.
	maxLevenshtein = 2
	// minRuneLenForDistance avoids bogus matches on very short names (e.g. mais vs mai).
	minRuneLenForDistance = 4
)

// Sprite holds grid position (row, col) for a 32x32 sprite.
type Sprite struct {
	Row, Col int
}

// X returns the x coordinate in pixels.
func (s Sprite) X() int { return s.Col * spriteSize }

// Y returns the y coordinate in pixels.
func (s Sprite) Y() int { return s.Row * spriteSize }

// Store holds the alias -> sprite index and optional list of sprites with aliases.
type Store struct {
	byAlias map[string]Sprite
}

// Load reads a bentext file and returns a Store. Format: blocks separated by "---",
// first line of each block is "row col", following lines are aliases. A line may list
// several aliases separated by "|" (whitespace around "|" is trimmed).
// The caller opens the file and controls the path; this function only consumes the reader.
func Load(r io.Reader) (*Store, error) {
	byAlias := make(map[string]Sprite)
	scanner := bufio.NewScanner(r)
	var block []string

	flushBlock := func() {
		if len(block) < 2 {
			block = nil
			return
		}
		parts := strings.SplitN(strings.TrimSpace(block[0]), " ", 2)
		if len(parts) != 2 {
			block = nil
			return
		}
		row, errRow := strconv.Atoi(strings.TrimSpace(parts[0]))
		col, errCol := strconv.Atoi(strings.TrimSpace(parts[1]))
		if errRow != nil || errCol != nil {
			block = nil
			return
		}
		sprite := Sprite{Row: row, Col: col}
		for i := 1; i < len(block); i++ {
			line := strings.TrimSpace(block[i])
			if line == "" {
				continue
			}
			for _, seg := range strings.Split(line, "|") {
				alias := strings.TrimSpace(seg)
				if alias == "" {
					continue
				}
				key := normalize(alias)
				byAlias[key] = sprite
			}
		}
		block = nil
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			flushBlock()
			continue
		}
		block = append(block, line)
	}
	flushBlock()

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read ingredient-sprites: %w", err)
	}
	return &Store{byAlias: byAlias}, nil
}

// Lookup returns sprite coordinates in pixels when the ingredient matches an alias, or falls back to
// the longest leading word sequence registered (e.g. "beurre doux" → sprite of "beurre"), then to
// the closest alias by Levenshtein distance if ≤ maxLevenshtein (see minRuneLenForDistance).
func (s *Store) Lookup(name string) (x, y int, ok bool) {
	sp, ok := s.lookupSprite(name)
	if !ok {
		return 0, 0, false
	}
	return sp.X(), sp.Y(), true
}

// LookupSprite returns the sprite when the ingredient matches, falls back via leading words, or
// matches within maxLevenshtein edits (with a minimum rune length on query and candidate).
func (s *Store) LookupSprite(name string) (sprite Sprite, ok bool) {
	return s.lookupSprite(name)
}

func (s *Store) lookupSprite(name string) (sprite Sprite, ok bool) {
	if s == nil || s.byAlias == nil {
		return Sprite{}, false
	}
	key := normalize(name)
	if sp, ok := s.byAlias[key]; ok {
		return sp, true
	}
	parts := strings.Fields(key)
	for n := len(parts); n >= 1; n-- {
		cand := strings.Join(parts[:n], " ")
		if sp, ok := s.byAlias[cand]; ok {
			return sp, true
		}
	}
	return lookupByLevenshtein(s.byAlias, key)
}

// lookupByLevenshtein picks the registered alias closest to query (≤ maxLevenshtein),
// requiring min length on both sides in runes. Tie-break: lower distance, then longer alias, then lexical.
func lookupByLevenshtein(byAlias map[string]Sprite, query string) (Sprite, bool) {
	q := []rune(query)
	if len(q) < minRuneLenForDistance {
		return Sprite{}, false
	}
	bestDist := maxLevenshtein + 1
	var best Sprite
	bestAlias := ""
	for alias, sp := range byAlias {
		a := []rune(alias)
		if min(len(a), len(q)) < minRuneLenForDistance {
			continue
		}
		d := levenshtein(a, q)
		if d > maxLevenshtein {
			continue
		}
		if d < bestDist || (d == bestDist && len(a) > len(bestAlias)) || (d == bestDist && len(a) == len(bestAlias) && alias < bestAlias) {
			bestDist = d
			best = sp
			bestAlias = alias
		}
	}
	if bestDist > maxLevenshtein {
		return Sprite{}, false
	}
	return best, true
}

func levenshtein(a, b []rune) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		ai := a[i-1]
		for j := 1; j <= len(b); j++ {
			cost := 0
			if ai != b[j-1] {
				cost = 1
			}
			cur[j] = min3(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// AllByAlias returns a map of normalized alias -> pixel coordinates (X, Y) for the sprite sheet.
func (s *Store) AllByAlias() map[string]struct{ X, Y int } {
	if s == nil || s.byAlias == nil {
		return nil
	}
	out := make(map[string]struct{ X, Y int }, len(s.byAlias))
	for alias, sp := range s.byAlias {
		out[alias] = struct{ X, Y int }{X: sp.X(), Y: sp.Y()}
	}
	return out
}

// normalize returns a key for matching: trim, lowercase.
func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
