package ingredients

import (
	"strings"
	"testing"
)

func TestLoadPipeAliases(t *testing.T) {
	const src = `---
2 9
margarine | beurre végétal
---
`
	s, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	want := Sprite{Row: 2, Col: 9}
	for _, key := range []string{"margarine", "beurre végétal"} {
		x, y, ok := s.Lookup(key)
		if !ok || x != want.X() || y != want.Y() {
			t.Fatalf("Lookup(%q) = (%d,%d,%v); want sprite %d,%d true", key, x, y, ok, want.X(), want.Y())
		}
	}
}

func TestLookupWordPrefixFallback(t *testing.T) {
	const src = `---
16 1
beurre
---
`
	s, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		query string
		want  Sprite
	}{
		{"beurre", Sprite{Row: 16, Col: 1}},
		{"beurre doux", Sprite{Row: 16, Col: 1}},
		{"BEURRE DEMI SEL", Sprite{Row: 16, Col: 1}},
	}
	for _, tc := range tests {
		sp, ok := s.LookupSprite(tc.query)
		if !ok || sp != tc.want {
			t.Fatalf("LookupSprite(%q) = (%+v,%v); want %+v true", tc.query, sp, ok, tc.want)
		}
	}
}

func TestLookupLongestWordPrefix(t *testing.T) {
	const src = `---
6 0
citron
---
6 1
citron vert
---
`
	s, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	got, ok := s.LookupSprite("citron vert")
	if !ok || got != (Sprite{Row: 6, Col: 1}) {
		t.Fatalf("citron vert: got %+v %v want row 6 col 1", got, ok)
	}
	got, ok = s.LookupSprite("citron jaune")
	if !ok || got != (Sprite{Row: 6, Col: 0}) {
		t.Fatalf("citron jaune: got %+v %v want citron sprite", got, ok)
	}
}

func TestLookupNoFalsePrefix(t *testing.T) {
	const src = `---
1 1
mai
---
`
	s, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	_, ok := s.LookupSprite("mais")
	if ok {
		t.Fatal("mai should not match mais")
	}
}

func TestLookupLevenshteinTypo(t *testing.T) {
	const src = `---
16 1
beurre
---
`
	s, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	want := Sprite{Row: 16, Col: 1}
	sp, ok := s.LookupSprite("beurr") // one missing char → distance 1
	if !ok || sp != want {
		t.Fatalf("LookupSprite(beurr): got %+v %v want %+v true", sp, ok, want)
	}
	sp, ok = s.LookupSprite("beure") // swapped / missing → distance ≤ 2
	if !ok || sp != want {
		t.Fatalf("LookupSprite(beure): got %+v %v want %+v true", sp, ok, want)
	}
}

func TestLookupLevenshteinTooFarRejected(t *testing.T) {
	const src = `---
16 1
beurre
---
`
	s, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.LookupSprite("zzbeurrezzz"); ok {
		t.Fatal("expected no match beyond maxLevenshtein")
	}
}

// If two aliases share the same minimal Levenshtein distance, the longer alias wins (then lexical).
func TestLookupLevenshteinTiePreferLongerAlias(t *testing.T) {
	const src = `---
1 9
abcd
---
1 10
abcde
---
`
	s, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	sp, ok := s.LookupSprite("abcef") // edit dist 2 to both tokens
	if !ok || sp != (Sprite{Row: 1, Col: 10}) {
		t.Fatalf("expected longer-alias sprite {1 10}, got %+v ok=%v", sp, ok)
	}
}

func TestLookupLevenshteinShortQuerySkipped(t *testing.T) {
	const src = `---
18 1
salt
---
`
	s, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.LookupSprite("slt"); ok {
		t.Fatal("query too short rune-wise must not use distance fallback")
	}
}
