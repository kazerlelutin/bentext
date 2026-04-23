package recipe

import "testing"

func TestMergeCoverEating(t *testing.T) {
	t.Parallel()
	if got := MergeCoverEating("Optionnel", "À la main"); got != "À la main ~ Couverts" {
		t.Errorf("Optionnel+À la main: got %q", got)
	}
	if got := MergeCoverEating("Non", "Baguettes"); got != "Baguettes" {
		t.Errorf("Non+Baguettes: got %q", got)
	}
	if got := MergeCoverEating("Optionnel", "Baguettes ~ Couverts"); got != "Baguettes ~ Couverts" {
		t.Errorf("optional+already combined: got %q", got)
	}
}

func TestParseValueOnlyBentoNewFormat(t *testing.T) {
	t.Parallel()
	raw := `R
1
D
---
X|1|g
---
S
---
t
---
Facile
Non
Au frais
À la main ~ Couverts
Faible
Discrète
Moyen
Note
Extra
`
	r := Parse(raw, "t", "fr")
	if r == nil || r.Bento == nil {
		t.Fatal("expected bento")
	}
	b := r.Bento
	if b.Transport != "Facile" || b.Reheat != "Non" || b.Cold != "Au frais" {
		t.Fatalf("core: %+v", b)
	}
	if b.Utensils != "À la main ~ Couverts" {
		t.Fatalf("utensils: %q", b.Utensils)
	}
	if b.Stains != "Faible" || b.Smell != "Discrète" || b.PrepTime != "Moyen" || b.Holding != "Note" || b.ExtraNotes != "Extra" {
		t.Fatalf("optional: %+v", b)
	}
}

func TestParseLegacyCoverEatingLines(t *testing.T) {
	t.Parallel()
	// Ancien fichier 10 lignes (cover + eating séparés)
	lines := []string{
		"Facile", "Non", "Non",
		"Optionnel", "À la main",
		"Faible", "Discrète", "Moyen", "H", "X",
	}
	b := &Bento{}
	parseValueOnlyBentoLines(lines, b)
	if b.Utensils != "À la main ~ Couverts" {
		t.Fatalf("got %q", b.Utensils)
	}
	if b.Stains != "Faible" {
		t.Fatalf("stains: %q", b.Stains)
	}
}

func TestMigrateBentoValueLines(t *testing.T) {
	t.Parallel()
	in := []string{"a", "b", "c", "Optionnel", "À la main", "z"}
	out, ok := MigrateBentoValueLines(in)
	if !ok {
		t.Fatal("expected ok")
	}
	if len(out) != 5 || out[3] != "À la main ~ Couverts" || out[4] != "z" {
		t.Fatalf("%v", out)
	}
}
