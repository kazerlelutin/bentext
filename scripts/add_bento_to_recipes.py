#!/usr/bin/env python3
"""Ajoute la section Bento (Préfixe|valeur) aux recettes .bentext si absente."""

from __future__ import annotations

import re
from pathlib import Path

RECIPE_ROOT = Path(__file__).resolve().parent.parent / "recipes"
LANGS = ("fr", "en", "ja", "zh", "ko")

# (Transport, Réchauffage, Froid, Manger) par langue
PRESETS: dict[str, dict[str, tuple[str, str, str, str]]] = {
    "baked_sweet": {
        "fr": (
            "Facile",
            "Optionnel (four ou micro-ondes)",
            "Non",
            "À la main ou Couverts",
        ),
        "en": (
            "Easy",
            "Optional (oven or microwave)",
            "No",
            "By hand or cutlery",
        ),
        "ja": (
            "簡単",
            "任意（オーブンまたは電子レンジ）",
            "不要",
            "手づかみまたはカトラリー",
        ),
        "zh": ("容易", "可选（烤箱或微波炉）", "否", "手抓或餐具"),
        "ko": (
            "쉬움",
            "선택(오븐 또는 전자레인지)",
            "아니오",
            "손으로 또는 수저",
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
            "Optionnel micro-ondes",
            "Conserver au frais jusqu'à consommation",
            "Couverts",
        ),
        "en": (
            "Medium",
            "Optional microwave",
            "Keep chilled until eaten",
            "Cutlery",
        ),
        "ja": (
            "やや注意",
            "任意（電子レンジ）",
            "食べるまで冷蔵",
            "カトラリー",
        ),
        "zh": ("中等", "可选微波炉", "食用前冷藏", "餐具"),
        "ko": (
            "보통",
            "선택(전자레인지)",
            "먹을 때까지 냉장",
            "수저",
        ),
    },
    "tarte_meringue": {
        "fr": ("Délicat", "Non", "Conserver au frais", "Couverts"),
        "en": ("Delicate", "No", "Keep chilled", "Cutlery"),
        "ja": ("壊れやすい", "不要", "冷蔵", "カトラリー"),
        "zh": ("易损", "不需要", "冷藏", "餐具"),
        "ko": ("다소 까다로움", "불필요", "냉장", "수저"),
    },
    "choux": {
        "fr": (
            "Délicat",
            "Non",
            "Conserver au frais jusqu'à consommation",
            "Couverts ou à la main",
        ),
        "en": (
            "Delicate",
            "No",
            "Keep chilled until eaten",
            "Cutlery or by hand",
        ),
        "ja": (
            "壊れやすい",
            "不要",
            "食べるまで冷蔵",
            "カトラリーまたは手づかみ",
        ),
        "zh": ("易损", "不需要", "食用前冷藏", "餐具或手抓"),
        "ko": (
            "다소 까다로움",
            "불필요",
            "먹을 때까지 냉장",
            "수저 또는 손으로",
        ),
    },
    "dacquoise": {
        "fr": ("Délicat", "Non", "Conserver au frais", "Couverts"),
        "en": ("Delicate", "No", "Keep chilled", "Cutlery"),
        "ja": ("壊れやすい", "不要", "冷蔵", "カトラリー"),
        "zh": ("易损", "不需要", "冷藏", "餐具"),
        "ko": ("다소 까다로움", "불필요", "냉장", "수저"),
    },
    "biscuit_roule": {
        "fr": ("Délicat", "Non", "Conserver au frais", "Couverts"),
        "en": ("Delicate", "No", "Keep chilled", "Cutlery"),
        "ja": ("壊れやすい", "不要", "冷蔵", "カトラリー"),
        "zh": ("易损", "不需要", "冷藏", "餐具"),
        "ko": ("다소 까다로움", "불필요", "냉장", "수저"),
    },
    "viennoiserie": {
        "fr": (
            "Facile",
            "Optionnel four",
            "Non",
            "À la main",
        ),
        "en": ("Easy", "Optional (oven)", "No", "By hand"),
        "ja": ("簡単", "任意（オーブン）", "不要", "手づかみ"),
        "zh": ("容易", "可选烤箱", "否", "手抓"),
        "ko": ("쉬움", "선택(오븐)", "아니오", "손으로"),
    },
    "empanada": {
        "fr": ("Facile", "Optionnel four", "Non", "À la main"),
        "en": ("Easy", "Optional (oven)", "No", "By hand"),
        "ja": ("簡単", "任意（オーブン）", "不要", "手づかみ"),
        "zh": ("容易", "可选烤箱", "否", "手抓"),
        "ko": ("쉬움", "선택(오븐)", "아니오", "손으로"),
    },
    "onigiri_gimbap": {
        "fr": (
            "Facile",
            "Optionnel micro-ondes",
            "À manger frais ou à température ambiante selon garniture",
            "À la main ou Baguettes",
        ),
        "en": (
            "Easy",
            "Optional microwave",
            "Eat fresh or at room temperature depending on filling",
            "By hand or chopsticks",
        ),
        "ja": (
            "簡単",
            "任意（電子レンジ）",
            "具材に応じて冷蔵のまままたは常温で",
            "手づかみまたは箸",
        ),
        "zh": (
            "容易",
            "可选微波炉",
            "视馅料冷藏或常温食用",
            "手抓或筷子",
        ),
        "ko": (
            "쉬움",
            "선택(전자레인지)",
            "속재료에 따라 냉장 또는 실온",
            "손으로 또는 젓가락",
        ),
    },
    "mandu": {
        "fr": (
            "Moyen",
            "Recommandé vapeur ou micro-ondes",
            "Conserver au frais si délai long",
            "Baguettes ou Couverts",
        ),
        "en": (
            "Medium",
            "Recommended—steam or microwave",
            "Keep chilled if long delay",
            "Chopsticks or cutlery",
        ),
        "ja": (
            "やや注意",
            "推奨（蒸しまたは電子レンジ）",
            "時間が空く場合は冷蔵",
            "箸またはカトラリー",
        ),
        "zh": ("中等", "建议蒸或微波", "间隔较久需冷藏", "筷子或餐具"),
        "ko": (
            "보통",
            "권장(찜 또는 전자레인지)",
            "시간 간격이 길면 냉장",
            "젓가락 또는 수저",
        ),
    },
    "guacamole": {
        "fr": (
            "Moyen",
            "Non",
            "Conserver au frais jusqu'à consommation",
            "Couverts (tortillas à la main)",
        ),
        "en": (
            "Medium",
            "No",
            "Keep chilled until eaten",
            "Cutlery (tortillas by hand)",
        ),
        "ja": (
            "やや注意",
            "不要",
            "食べるまで冷蔵",
            "カトラリー（トルティーヤは手づかみ）",
        ),
        "zh": ("中等", "不需要", "食用前冷藏", "餐具（玉米饼可手抓）"),
        "ko": (
            "보통",
            "불필요",
            "먹을 때까지 냉장",
            "수저(또띠아는 손으로)",
        ),
    },
    "omelette": {
        "fr": (
            "Moyen",
            "Recommandé tiède",
            "Conserver au frais si délai long",
            "Couverts",
        ),
        "en": (
            "Medium",
            "Recommended warm",
            "Keep chilled if long delay",
            "Cutlery",
        ),
        "ja": (
            "やや注意",
            "温かいうちが望ましい",
            "時間が空く場合は冷蔵",
            "カトラリー",
        ),
        "zh": ("中等", "建议温热食用", "间隔较久需冷藏", "餐具"),
        "ko": (
            "보통",
            "따뜻할 때 권장",
            "시간 간격이 길면 냉장",
            "수저",
        ),
    },
    "baba": {
        "fr": ("Moyen", "Non", "Conserver au frais", "Couverts"),
        "en": ("Medium", "No", "Keep chilled", "Cutlery"),
        "ja": ("やや注意", "不要", "冷蔵", "カトラリー"),
        "zh": ("中等", "不需要", "冷藏", "餐具"),
        "ko": ("보통", "불필요", "냉장", "수저"),
    },
    "blinis": {
        "fr": (
            "Moyen",
            "Optionnel micro-ondes",
            "Conserver au frais si garniture fraîche",
            "Couverts ou à la main",
        ),
        "en": (
            "Medium",
            "Optional microwave",
            "Keep chilled if using fresh toppings",
            "Cutlery or by hand",
        ),
        "ja": (
            "やや注意",
            "任意（電子レンジ）",
            "生のトッピングなら冷蔵",
            "カトラリーまたは手づかみ",
        ),
        "zh": ("中等", "可选微波炉", "新鲜配料需冷藏", "餐具或手抓"),
        "ko": (
            "보통",
            "선택(전자레인지)",
            "신선한 토핑이면 냉장",
            "수저 또는 손으로",
        ),
    },
    "naan": {
        "fr": ("Facile", "Optionnel four ou poêle", "Non", "À la main"),
        "en": ("Easy", "Optional (oven or pan)", "No", "By hand"),
        "ja": ("簡単", "任意（オーブンまたはフライパン）", "不要", "手づかみ"),
        "zh": ("容易", "可选烤箱或平底锅", "否", "手抓"),
        "ko": ("쉬움", "선택(오븐 또는 프라이팬)", "아니오", "손으로"),
    },
    "dan_bing": {
        "fr": ("Facile", "Optionnel poêle", "Non", "À la main"),
        "en": ("Easy", "Optional (pan)", "No", "By hand"),
        "ja": ("簡単", "任意（フライパン）", "不要", "手づかみ"),
        "zh": ("容易", "可选平底锅", "否", "手抓"),
        "ko": ("쉬움", "선택(프라이팬)", "아니오", "손으로"),
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
    if 4 <= len(lines) <= 9 and all("|" not in ln for ln in lines):
        return True
    return False


def format_bento(lang: str, preset: str) -> str:
    t = PRESETS[preset][lang]
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
