#!/usr/bin/env python3
"""Ajoute la section Bento (Préfixe|valeur) aux recettes .bentext si absente."""

from __future__ import annotations

import re
from pathlib import Path

RECIPE_ROOT = Path(__file__).resolve().parent.parent / "recipes"
LANGS = ("fr", "en", "ja", "zh", "ko")

COVER_DEFAULT = {"fr": "Optionnel", "en": "Optional", "ja": "任意", "zh": "可选", "ko": "선택"}

# Par langue : (Transport, Réchauffage, Froid, Manger) — la ligne « couvert » est insérée automatiquement avant Manger.
PRESETS: dict[str, dict[str, tuple[str, str, str, str]]] = {
    "baked_sweet": {
        "fr": (
            "Facile",
            "Optionnel four ~ micro-ondes",
            "Non",
            "À la main",
        ),
        "en": (
            "Easy",
            "Optional oven ~ microwave",
            "No",
            "By hand",
        ),
        "ja": (
            "簡単",
            "任意 オーブン ~ 電子レンジ",
            "不要",
            "手づかみ",
        ),
        "zh": ("容易", "可选 烤箱 ~ 微波炉", "否", "手抓"),
        "ko": (
            "쉬움",
            "선택 오븐 ~ 전자레인지",
            "아니오",
            "손으로",
        ),
    },
    "fried": {
        "fr": (
            "Délicat",
            "Non",
            "Non, consommer de préférence le jour même",
            "À la main",
        ),
        "en": (
            "Delicate",
            "No",
            "No; best eaten the same day",
            "By hand",
        ),
        "ja": (
            "壊れやすい",
            "不要",
            "できれば当日中に",
            "手づかみ",
        ),
        "zh": ("易损", "不需要", "尽量当日食用", "手抓"),
        "ko": ("다소 까다로움", "불필요", "가능하면 당일 섭취", "손으로"),
    },
    "custard": {
        "fr": (
            "Moyen",
            "Optionnel",
            "Au frais",
            "Couverts",
        ),
        "en": (
            "Medium",
            "Optional",
            "Keep cold",
            "Cutlery",
        ),
        "ja": (
            "やや注意",
            "任意",
            "要冷蔵",
            "カトラリー",
        ),
        "zh": ("中等", "可选", "需冷藏", "餐具"),
        "ko": (
            "보통",
            "선택",
            "냉장",
            "수저",
        ),
    },
    "tarte_meringue": {
        "fr": ("Délicat", "Non", "Au frais", "Couverts"),
        "en": ("Delicate", "No", "Chilled", "Cutlery"),
        "ja": ("壊れやすい", "不要", "冷蔵", "カトラリー"),
        "zh": ("易损", "不需要", "冷藏", "餐具"),
        "ko": ("다소 까다로움", "불필요", "냉장", "수저"),
    },
    "choux": {
        "fr": (
            "Délicat",
            "Non",
            "Au frais",
            "Couverts ~ main",
        ),
        "en": (
            "Delicate",
            "No",
            "Keep cold",
            "Cutlery ~ hand",
        ),
        "ja": (
            "壊れやすい",
            "不要",
            "要冷蔵",
            "カトラリー ~ 手づかみ",
        ),
        "zh": ("易损", "不需要", "需冷藏", "餐具 ~ 手抓"),
        "ko": (
            "다소 까다로움",
            "불필요",
            "냉장",
            "수저 ~ 손",
        ),
    },
    "dacquoise": {
        "fr": ("Délicat", "Non", "Au frais", "Couverts"),
        "en": ("Delicate", "No", "Chilled", "Cutlery"),
        "ja": ("壊れやすい", "不要", "冷蔵", "カトラリー"),
        "zh": ("易损", "不需要", "冷藏", "餐具"),
        "ko": ("다소 까다로움", "불필요", "냉장", "수저"),
    },
    "biscuit_roule": {
        "fr": ("Délicat", "Non", "Au frais", "Couverts"),
        "en": ("Delicate", "No", "Chilled", "Cutlery"),
        "ja": ("壊れやすい", "不要", "冷蔵", "カトラリー"),
        "zh": ("易损", "不需要", "冷藏", "餐具"),
        "ko": ("다소 까다로움", "불필요", "냉장", "수저"),
    },
    "viennoiserie": {
        "fr": (
            "Facile",
            "Optionnel",
            "Non",
            "À la main",
        ),
        "en": ("Easy", "Optional", "No", "By hand"),
        "ja": ("簡単", "任意", "不要", "手づかみ"),
        "zh": ("容易", "可选", "否", "手抓"),
        "ko": ("쉬움", "선택", "아니오", "손으로"),
    },
    "empanada": {
        "fr": ("Facile", "Optionnel", "Non", "À la main"),
        "en": ("Easy", "Optional", "No", "By hand"),
        "ja": ("簡単", "任意", "不要", "手づかみ"),
        "zh": ("容易", "可选", "否", "手抓"),
        "ko": ("쉬움", "선택", "아니오", "손으로"),
    },
    "onigiri_gimbap": {
        "fr": (
            "Facile",
            "Optionnel",
            "Frais ~ ambiant",
            "À la main ~ Baguettes",
        ),
        "en": (
            "Easy",
            "Optional",
            "Fresh ~ room temp",
            "Hand ~ chopsticks",
        ),
        "ja": (
            "簡単",
            "任意",
            "冷蔵 ~ 常温",
            "手づかみ ~ 箸",
        ),
        "zh": (
            "容易",
            "可选",
            "冷藏 ~ 常温",
            "手抓 ~ 筷子",
        ),
        "ko": (
            "쉬움",
            "선택",
            "냉장 ~ 실온",
            "손 ~ 젓가락",
        ),
    },
    "mandu": {
        "fr": (
            "Moyen",
            "Recommandé vapeur ~ micro-ondes",
            "Frais si délai",
            "Baguettes ~ Couverts",
        ),
        "en": (
            "Medium",
            "Recommended steam ~ microwave",
            "Chilled if delay",
            "Chopsticks ~ cutlery",
        ),
        "ja": (
            "やや注意",
            "推奨 蒸し ~ 電子レンジ",
            "遅延 ~ 冷蔵",
            "箸 ~ カトラリー",
        ),
        "zh": ("中等", "建议 蒸 ~ 微波", "久置 ~ 冷藏", "筷子 ~ 餐具"),
        "ko": (
            "보통",
            "권장 찜 ~ 전자레인지",
            "지연 ~ 냉장",
            "젓가락 ~ 수저",
        ),
    },
    "guacamole": {
        "fr": (
            "Moyen",
            "Non",
            "Au frais",
            "Couverts (tortillas à la main)",
        ),
        "en": (
            "Medium",
            "No",
            "Keep cold",
            "Cutlery (tortillas by hand)",
        ),
        "ja": (
            "やや注意",
            "不要",
            "要冷蔵",
            "カトラリー（トルティーヤは手づかみ）",
        ),
        "zh": ("中等", "不需要", "需冷藏", "餐具（玉米饼可手抓）"),
        "ko": (
            "보통",
            "불필요",
            "냉장",
            "수저(또띠아는 손으로)",
        ),
    },
    "omelette": {
        "fr": (
            "Moyen",
            "Recommandé tiède",
            "Frais si délai",
            "Couverts",
        ),
        "en": (
            "Medium",
            "Recommended warm",
            "Chilled if delay",
            "Cutlery",
        ),
        "ja": (
            "やや注意",
            "温かいうちが望ましい",
            "遅延 ~ 冷蔵",
            "カトラリー",
        ),
        "zh": ("中等", "建议温热食用", "久置 ~ 冷藏", "餐具"),
        "ko": (
            "보통",
            "따뜻할 때 권장",
            "지연 ~ 냉장",
            "수저",
        ),
    },
    "baba": {
        "fr": ("Moyen", "Non", "Au frais", "Couverts"),
        "en": ("Medium", "No", "Chilled", "Cutlery"),
        "ja": ("やや注意", "不要", "冷蔵", "カトラリー"),
        "zh": ("中等", "不需要", "冷藏", "餐具"),
        "ko": ("보통", "불필요", "냉장", "수저"),
    },
    "blinis": {
        "fr": (
            "Moyen",
            "Optionnel",
            "Frais si garniture",
            "Couverts ~ main",
        ),
        "en": (
            "Medium",
            "Optional",
            "Chilled if fresh toppings",
            "Cutlery ~ hand",
        ),
        "ja": (
            "やや注意",
            "任意",
            "生のトッピングなら冷蔵",
            "カトラリー ~ 手づかみ",
        ),
        "zh": ("中等", "可选", "新鲜配料需冷藏", "餐具 ~ 手抓"),
        "ko": (
            "보통",
            "선택",
            "신선한 토핑이면 냉장",
            "수저 ~ 손",
        ),
    },
    "naan": {
        "fr": ("Facile", "Optionnel four ~ poêle", "Non", "À la main"),
        "en": ("Easy", "Optional oven ~ pan", "No", "By hand"),
        "ja": ("簡単", "任意 オーブン ~ フライパン", "不要", "手づかみ"),
        "zh": ("容易", "可选 烤箱 ~ 平底锅", "否", "手抓"),
        "ko": ("쉬움", "선택 오븐 ~ 프라이팬", "아니오", "손으로"),
    },
    "dan_bing": {
        "fr": ("Facile", "Optionnel", "Non", "À la main"),
        "en": ("Easy", "Optional", "No", "By hand"),
        "ja": ("簡単", "任意", "不要", "手づかみ"),
        "zh": ("容易", "可选", "否", "手抓"),
        "ko": ("쉬움", "선택", "아니오", "손으로"),
    },
}


def choose_preset(slug: str) -> str:
    s = slug.lower()
    if "guacamole" in s:
        return "guacamole"
    if "beignet" in s:
        return "fried"
    if "dan-bing" in s:
        return "dan_bing"
    if "omelette" in s:
        return "omelette"
    if "onigiri" in s or "gimbap" in s:
        return "onigiri_gimbap"
    if "mandu" in s:
        return "mandu"
    if "empanada" in s or "empanadas" in s:
        return "empanada"
    if "choux" in s:
        return "choux"
    if s == "tarte-au-citron-meringuee":
        return "tarte_meringue"
    if s in (
        "flan",
        "creme-caramel",
        "creme-brulee-a-la-vanille",
        "cremeux",
    ) or "flan-patissier" in s:
        return "custard"
    if "baba-au-rhum" in s:
        return "baba"
    if "dacquoise" in s:
        return "dacquoise"
    if "biscuit-roule" in s:
        return "biscuit_roule"
    if "cheese-naan" in s:
        return "naan"
    if "blinis" in s:
        return "blinis"
    if "brioche" in s or "croissants" == s:
        return "viennoiserie"
    return "baked_sweet"


def split_sections(text: str) -> list[str]:
    text = text.replace("\r\n", "\n").replace("\r", "\n")
    parts = re.split(r"^---\s*$", text, flags=re.MULTILINE)
    return [p.strip("\n") for p in parts]


def has_bento_block(parts: list[str]) -> bool:
    if not parts:
        return False
    last = parts[-1].strip()
    lines = [ln.strip() for ln in last.split("\n") if ln.strip()]
    if not lines:
        return False
    if lines[0].startswith("Transport|"):
        return True
    if 4 <= len(lines) <= 10 and all("|" not in ln for ln in lines):
        return True
    return False


def format_bento(lang: str, preset: str) -> str:
    t = PRESETS[preset][lang]
    if len(t) == 4:
        c = COVER_DEFAULT[lang]
        t = (t[0], t[1], t[2], c, t[3])
    return "\n".join(t)


def process_file(path: Path, lang: str, preset: str) -> bool:
    raw = path.read_text(encoding="utf-8")
    parts = split_sections(raw)
    if has_bento_block(parts):
        return False
    if not raw.endswith("\n"):
        raw = raw + "\n"
    bento = format_bento(lang, preset)
    new_content = raw.rstrip() + "\n---\n" + bento + "\n"
    path.write_text(new_content, encoding="utf-8", newline="\n")
    return True


def main() -> None:
    by_slug: dict[str, dict[str, Path]] = {}
    for p in RECIPE_ROOT.glob("*.bentext"):
        m = re.match(r"^(.+)\.(fr|en|ja|zh|ko)\.bentext$", p.name)
        if not m:
            continue
        slug, lang = m.group(1), m.group(2)
        by_slug.setdefault(slug, {})[lang] = p

    updated = 0
    skipped = 0
    for slug in sorted(by_slug):
        preset = choose_preset(slug)
        for lang in LANGS:
            path = by_slug[slug].get(lang)
            if path is None:
                continue
            if process_file(path, lang, preset):
                updated += 1
            else:
                skipped += 1

    print(f"Fichiers mis à jour : {updated}, déjà OK ou ignorés : {skipped}")


if __name__ == "__main__":
    main()
