package ingredients

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

const dataFile = "ingredient-sprites.bentext"

var (
	mu          sync.Mutex
	store       *Store
	loadErr     error
	cachedPath  string
	cachedMTime time.Time
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

func emptyStore() *Store {
	return &Store{byAlias: make(map[string]Sprite)}
}

// ensureLoaded loads or reloads the store when ingredient-sprites.bentext changes (mtime).
func ensureLoaded() {
	path := dataPath()
	info, statErr := os.Stat(path)

	mu.Lock()
	defer mu.Unlock()

	if statErr != nil {
		store = emptyStore()
		loadErr = statErr
		cachedPath = ""
		cachedMTime = time.Time{}
		return
	}

	if store != nil && path == cachedPath && info.ModTime().Equal(cachedMTime) {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		store = emptyStore()
		loadErr = err
		cachedPath = ""
		cachedMTime = time.Time{}
		return
	}
	defer f.Close()

	st, err := Load(f)
	loadErr = err
	if st == nil {
		st = emptyStore()
	}
	store = st
	cachedPath = path
	cachedMTime = info.ModTime()
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
