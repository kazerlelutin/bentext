package recipe

import "strings"

// MergeCoverEating combine l’ancienne paire Couvert + Manger en une seule ligne ustensiles (localisée).
// Utilisé pour la lecture des anciens fichiers et pour la migration de contenu.
func MergeCoverEating(cover, eating string) string {
	c := strings.TrimSpace(cover)
	e := strings.TrimSpace(eating)
	if c == "" {
		return e
	}
	if e == "" {
		return c
	}
	// Déjà une combinaison explicite sur la ligne manger
	if strings.Contains(e, "~") {
		return e
	}
	if isCoverNonToken(c) {
		return e
	}
	if isCoverYesToken(c) {
		return e
	}
	if isCoverOptionalToken(c) {
		return mergeOptionalCover(e)
	}
	return e
}

func isCoverNonToken(s string) bool {
	lower := strings.ToLower(s)
	switch lower {
	case "non", "no":
		return true
	}
	switch s {
	case "不要", "不需要", "불필요":
		return true
	}
	return false
}

func isCoverYesToken(s string) bool {
	lower := strings.ToLower(s)
	switch lower {
	case "oui", "yes":
		return true
	}
	switch s {
	case "必要", "需要", "필요":
		return true
	}
	return false
}

func isCoverOptionalToken(s string) bool {
	lower := strings.ToLower(s)
	switch lower {
	case "optionnel", "optional":
		return true
	}
	switch s {
	case "任意", "可选", "선택":
		return true
	}
	return false
}

// mergeOptionalCover : cover=Optionnel + eating court → main ~ couverts (par langue détectée sur eating).
func mergeOptionalCover(eating string) string {
	e := strings.TrimSpace(eating)
	lower := strings.ToLower(e)

	// FR
	switch e {
	case "À la main":
		return "À la main ~ Couverts"
	case "Couverts":
		return "Couverts"
	case "Baguettes":
		return "Baguettes ~ Couverts"
	}
	// EN
	switch lower {
	case "by hand":
		return "By hand ~ cutlery"
	case "cutlery":
		return "Cutlery"
	case "chopsticks":
		return "Chopsticks ~ cutlery"
	}
	// JA
	switch e {
	case "手づかみ":
		return "手づかみ ~ カトラリー"
	case "カトラリー":
		return "カトラリー"
	case "箸":
		return "箸 ~ カトラリー"
	}
	// ZH
	switch e {
	case "手抓":
		return "手抓 ~ 餐具"
	case "餐具":
		return "餐具"
	case "筷子":
		return "筷子 ~ 餐具"
	}
	// KO
	switch e {
	case "손으로":
		return "손으로 ~ 수저"
	case "수저":
		return "수저"
	case "젓가락":
		return "젓가락 ~ 수저"
	}
	return eating + " ~ Couverts"
}

// MigrateBentoValueLines fusionne les lignes legacy cover+eating (indices 3–4) en une seule ligne utensils.
// Retourne ok=true si une modification a été faite.
func MigrateBentoValueLines(lines []string) ([]string, bool) {
	n := len(lines)
	if n < 5 || n > 10 {
		return lines, false
	}
	l4 := strings.TrimSpace(lines[3])
	lower := strings.ToLower(l4)
	if !isShortCoverToken(l4, lower) {
		return lines, false
	}
	merged := MergeCoverEating(lines[3], lines[4])
	out := make([]string, 0, n-1)
	out = append(out, lines[0], lines[1], lines[2], merged)
	out = append(out, lines[5:]...)
	return out, true
}
