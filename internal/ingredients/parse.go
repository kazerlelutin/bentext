package ingredients

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const spriteSize = 32

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
// first line of each block is "row col", following lines are one alias per line.
// Returns nil if the file cannot be read (caller may use empty Store).
func Load(path string) (*Store, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byAlias := make(map[string]Sprite)
	scanner := bufio.NewScanner(f)
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
			alias := strings.TrimSpace(block[i])
			if alias == "" {
				continue
			}
			key := normalize(alias)
			byAlias[key] = sprite
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

// Lookup returns sprite coordinates in pixels and true if the ingredient name matches an alias.
func (s *Store) Lookup(name string) (x, y int, ok bool) {
	if s == nil || s.byAlias == nil {
		return 0, 0, false
	}
	key := normalize(name)
	sp, ok := s.byAlias[key]
	if !ok {
		return 0, 0, false
	}
	return sp.X(), sp.Y(), true
}

// LookupSprite returns the sprite and the canonical alias for the name, or false.
// The canonical alias is the first one registered for that sprite (arbitrary); for API we may use the requested name.
func (s *Store) LookupSprite(name string) (sprite Sprite, ok bool) {
	if s == nil || s.byAlias == nil {
		return Sprite{}, false
	}
	key := normalize(name)
	sp, ok := s.byAlias[key]
	return sp, ok
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
