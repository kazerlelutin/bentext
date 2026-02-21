package ingredients

import (
	"os"
	"path/filepath"
	"sync"
)

const dataFile = "ingredient-sprites.bentext"

var (
	once   sync.Once
	store  *Store
	loadErr error
)

// dataPath returns the path to ingredient-sprites.bentext (cwd then next to executable).
func dataPath() string {
	cwd := filepath.Join(".", dataFile)
	if _, err := os.Stat(cwd); err == nil {
		return cwd
	}
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exe)
		if _, err := os.Stat(filepath.Join(dir, dataFile)); err == nil {
			return filepath.Join(dir, dataFile)
		}
	}
	return cwd
}

// ensureLoaded loads the store once.
func ensureLoaded() {
	once.Do(func() {
		path := dataPath()
		store, loadErr = Load(path)
		if loadErr != nil && store == nil {
			store = &Store{byAlias: make(map[string]Sprite)}
		}
	})
}

// Lookup returns sprite coordinates in pixels for the given ingredient name (any language/alias).
// Uses lazy-loaded store. Returns (0, 0, false) if not found or file missing.
func Lookup(name string) (x, y int, ok bool) {
	ensureLoaded()
	return store.Lookup(name)
}

// LookupSprite returns the sprite for the name, or false.
func LookupSprite(name string) (sprite Sprite, ok bool) {
	ensureLoaded()
	return store.LookupSprite(name)
}

// AllByAlias returns the full alias -> {X,Y} map for the sprite sheet endpoint.
func AllByAlias() map[string]struct{ X, Y int } {
	ensureLoaded()
	return store.AllByAlias()
}
