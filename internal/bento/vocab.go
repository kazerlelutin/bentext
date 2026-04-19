// Package bento holds canonical vocabulary rows for lunch-box metadata (documentation / presets).
// API responses use localized strings from recipe files; these tables align FR reference with other langs.
package bento

// VocabRow is one semantic row (same meaning across languages).
type VocabRow struct {
	Key string // stable id for docs, e.g. "transport_easy"
	FR  string
	EN  string
	JA  string
	ZH  string
	KO  string
}

// TransportEase rows (examples; recipes may add short precision after the scale).
var TransportEase = []VocabRow{
	{Key: "transport_easy", FR: "Facile", EN: "Easy", JA: "簡単", ZH: "容易", KO: "쉬움"},
	{Key: "transport_medium", FR: "Moyen", EN: "Medium", JA: "やや注意", ZH: "中等", KO: "보통"},
	{Key: "transport_delicate", FR: "Délicat", EN: "Delicate", JA: "壊れやすい", ZH: "易损", KO: "다소 까다로움"},
}

// Reheat : échelle courte (Optionnel, Non, …) ; plusieurs modes avec « ~ » comme les alternatives d’ingrédients (ex. four ~ micro-ondes).
var Reheat = []VocabRow{
	{Key: "reheat_no", FR: "Non", EN: "No", JA: "不要", ZH: "不需要", KO: "불필요"},
	{Key: "reheat_optional", FR: "Optionnel", EN: "Optional", JA: "任意", ZH: "可选", KO: "선택"},
	{Key: "reheat_optional_oven_micro", FR: "Optionnel four ~ micro-ondes", EN: "Optional oven ~ microwave", JA: "任意 オーブン ~ 電子レンジ", ZH: "可选 烤箱 ~ 微波炉", KO: "선택 오븐 ~ 전자레인지"},
	{Key: "reheat_optional_oven_pan", FR: "Optionnel four ~ poêle", EN: "Optional oven ~ pan", JA: "任意 オーブン ~ フライパン", ZH: "可选 烤箱 ~ 平底锅", KO: "선택 오븐 ~ 프라이팬"},
	{Key: "reheat_recommended_steam_micro", FR: "Recommandé vapeur ~ micro-ondes", EN: "Recommended steam ~ microwave", JA: "推奨 蒸し ~ 電子レンジ", ZH: "建议 蒸 ~ 微波", KO: "권장 찜 ~ 전자레인지"},
}

// Cold : conservation / chaîne du froid ; plusieurs modes avec « ~ » comme le réchauffage (ex. frais ~ ambiant, délai ~ froid).
var Cold = []VocabRow{
	{Key: "cold_no", FR: "Non", EN: "No", JA: "不要", ZH: "否", KO: "아니오"},
	{Key: "cold_keep_short", FR: "Au frais", EN: "Keep cold", JA: "要冷蔵", ZH: "需冷藏", KO: "냉장"},
	{Key: "cold_chilled_short", FR: "Au frais", EN: "Chilled", JA: "冷蔵", ZH: "冷藏", KO: "냉장"},
	{Key: "cold_delay", FR: "Frais si délai", EN: "Chilled if delay", JA: "遅延 ~ 冷蔵", ZH: "久置 ~ 冷藏", KO: "지연 ~ 냉장"},
	{Key: "cold_fresh_ambient", FR: "Frais ~ ambiant", EN: "Fresh ~ room temp", JA: "冷蔵 ~ 常温", ZH: "冷藏 ~ 常温", KO: "냉장 ~ 실온"},
}

// Utensils : une seule ligne fichier / champ JSON `utensils` — besoin de couverts et modalité (main, baguettes, couverts, combinaisons avec ` ~ `).
var Utensils = []VocabRow{
	{Key: "utensils_hand_cutlery", FR: "À la main ~ Couverts", EN: "By hand ~ cutlery", JA: "手づかみ ~ カトラリー", ZH: "手抓 ~ 餐具", KO: "손으로 ~ 수저"},
	{Key: "utensils_hand_chopsticks", FR: "À la main ~ Baguettes", EN: "Hand ~ chopsticks", JA: "手づかみ ~ 箸", ZH: "手抓 ~ 筷子", KO: "손 ~ 젓가락"},
	{Key: "utensils_chopsticks_cutlery", FR: "Baguettes ~ Couverts", EN: "Chopsticks ~ cutlery", JA: "箸 ~ カトラリー", ZH: "筷子 ~ 餐具", KO: "젓가락 ~ 수저"},
	{Key: "utensils_cutlery", FR: "Couverts", EN: "Cutlery", JA: "カトラリー", ZH: "餐具", KO: "수저"},
	{Key: "utensils_chopsticks", FR: "Baguettes", EN: "Chopsticks", JA: "箸", ZH: "筷子", KO: "젓가락"},
	{Key: "utensils_hand", FR: "À la main", EN: "By hand", JA: "手づかみ", ZH: "手抓", KO: "손으로"},
}

// Cover et Eating : ancienne décomposition (référence) ; le format fichier actuel n’a qu’une ligne `utensils`.
var Cover = []VocabRow{
	{Key: "cover_no", FR: "Non", EN: "No", JA: "不要", ZH: "不需要", KO: "불필요"},
	{Key: "cover_optional", FR: "Optionnel", EN: "Optional", JA: "任意", ZH: "可选", KO: "선택"},
	{Key: "cover_yes", FR: "Oui", EN: "Yes", JA: "必要", ZH: "需要", KO: "필요"},
}

var Eating = []VocabRow{
	{Key: "eating_hand", FR: "À la main", EN: "By hand", JA: "手づかみ", ZH: "手抓", KO: "손으로"},
	{Key: "eating_cutlery", FR: "Couverts", EN: "Cutlery", JA: "カトラリー", ZH: "餐具", KO: "수저"},
	{Key: "eating_chopsticks", FR: "Baguettes", EN: "Chopsticks", JA: "箸", ZH: "筷子", KO: "젓가락"},
	{Key: "eating_chopsticks_cutlery", FR: "Baguettes ~ Couverts", EN: "Chopsticks ~ cutlery", JA: "箸 ~ カトラリー", ZH: "筷子 ~ 餐具", KO: "젓가락 ~ 수저"},
	{Key: "eating_hand_chopsticks", FR: "À la main ~ Baguettes", EN: "Hand ~ chopsticks", JA: "手づかみ ~ 箸", ZH: "手抓 ~ 筷子", KO: "손 ~ 젓가락"},
}

// Stains risk.
var Stains = []VocabRow{
	{Key: "stains_none", FR: "Non", EN: "None", JA: "なし", ZH: "无", KO: "없음"},
	{Key: "stains_low", FR: "Faible", EN: "Low", JA: "低い", ZH: "低", KO: "낮음"},
	{Key: "stains_medium", FR: "Moyen", EN: "Medium", JA: "中程度", ZH: "中", KO: "보통"},
	{Key: "stains_high", FR: "Élevé", EN: "High", JA: "高い", ZH: "高", KO: "높음"},
}

// Smell in closed box.
var Smell = []VocabRow{
	{Key: "smell_discrete", FR: "Discrète", EN: "Mild", JA: "控えめ", ZH: "清淡", KO: "약함"},
	{Key: "smell_marked", FR: "Marquée", EN: "Strong", JA: "強い", ZH: "明显", KO: "강함"},
}

// Prep time scale (replaces make-ahead yes/no).
var PrepTime = []VocabRow{
	{Key: "prep_quick", FR: "Rapide", EN: "Quick", JA: "短時間", ZH: "快", KO: "짧음"},
	{Key: "prep_medium", FR: "Moyen", EN: "Medium", JA: "中程度", ZH: "中", KO: "보통"},
	{Key: "prep_long", FR: "Long", EN: "Long", JA: "長め", ZH: "长", KO: "긴 편"},
}
