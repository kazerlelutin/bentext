#!/usr/bin/env python3
"""Ajoute les lignes Bento optionnelles 6–10 (stains, smell, prep_time, holding, extra_notes) si absentes."""

from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent / "recipes"
LANGS = ("fr", "en", "ja", "zh", "ko")

# (stains, smell, prep_time, holding, extra_notes) par langue
OPTIONAL: dict[str, dict[str, tuple[str, str, str, str, str]]] = {
    "shitake-mandu": {
        "fr": (
            "Moyen",
            "Marquée",
            "Long",
            "Bon tiède ou à la vapeur",
            "Sauce soja en option à part",
        ),
        "en": (
            "Medium",
            "Strong",
            "Long",
            "Good warm or steamed",
            "Soy sauce on the side",
        ),
        "ja": (
            "中程度",
            "強い",
            "長め",
            "温かいか蒸しがよい",
            "醤油は別添え",
        ),
        "zh": ("中", "明显", "长", "温热或蒸食更佳", "酱油另配"),
        "ko": ("보통", "강함", "긴 편", "따뜻하거나 찜", "간장은 따로"),
    },
    "onigiri-kimchi-mozza": {
        "fr": (
            "Faible",
            "Marquée",
            "Moyen",
            "Riz moelleux si consommé vite",
            "Mozzarella peut durcir au froid",
        ),
        "en": (
            "Low",
            "Strong",
            "Medium",
            "Rice stays soft if eaten soon",
            "Mozzarella can firm up when cold",
        ),
        "ja": (
            "低い",
            "強い",
            "中程度",
            "早めならご飯は柔らかめ",
            "冷めるとチーズが硬くなりやすい",
        ),
        "zh": ("低", "明显", "中", "尽快吃米饭仍软", "冷后奶酪易变硬"),
        "ko": ("낮음", "강함", "보통", "빨리 먹으면 밥이 부드러움", "차가우면 치즈가 딱딱해질 수 있음"),
    },
    "gimbap-sweetcrunchy": {
        "fr": (
            "Faible",
            "Discrète",
            "Moyen",
            "Le nori ramollit avec le temps",
            "Prévoir de couper en tranches",
        ),
        "en": (
            "Low",
            "Mild",
            "Medium",
            "Nori softens over time",
            "Plan to slice before eating",
        ),
        "ja": (
            "低い",
            "控えめ",
            "中程度",
            "時間で海苔が軟らかくなる",
            "食べやすい幅に切る",
        ),
        "zh": ("低", "清淡", "中", "海苔会变软", "建议切段再吃"),
        "ko": ("낮음", "약함", "보통", "시간 지나면 김이 무름", "먹기 좋게 썰기"),
    },
    "omelette-roulee": {
        "fr": (
            "Faible",
            "Discrète",
            "Rapide",
            "Mieux le jour même",
            "Fragile au transport",
        ),
        "en": (
            "Low",
            "Mild",
            "Quick",
            "Best the same day",
            "Fragile to carry",
        ),
        "ja": (
            "低い",
            "控えめ",
            "短時間",
            "当日中がよい",
            "持ち運びでは崩れやすい",
        ),
        "zh": ("低", "清淡", "快", "当日最好", "携带易碎"),
        "ko": ("낮음", "약함", "짧음", "당일이 좋음", "휴대 시 잘 부서짐"),
    },
    "cremeux": {
        "fr": (
            "Moyen",
            "Discrète",
            "Long",
            "Surface peut filmer",
            "Boîte hermétique conseillée",
        ),
        "en": (
            "Medium",
            "Mild",
            "Long",
            "Surface may skin over",
            "Airtight container recommended",
        ),
        "ja": (
            "中程度",
            "控えめ",
            "長め",
            "表面が乾きやすい",
            "密閉容器がおすすめ",
        ),
        "zh": ("中", "清淡", "长", "表面易结皮", "建议密封盒"),
        "ko": ("보통", "약함", "긴 편", "표면 막이 생길 수 있음", "밀폐 용기 권장"),
    },
    "orange-cake": {
        "fr": (
            "Faible",
            "Discrète",
            "Long",
            "Bon le lendemain",
            "Repos long de pâte au frais",
        ),
        "en": (
            "Low",
            "Mild",
            "Long",
            "Good the next day",
            "Long chilled batter rest",
        ),
        "ja": (
            "低い",
            "控えめ",
            "長め",
            "翌日も美味しい",
            "生地の長時間冷蔵休息",
        ),
        "zh": ("低", "清淡", "长", "隔夜仍佳", "面糊需长时间冷藏静置"),
        "ko": ("낮음", "약함", "긴 편", "다음 날도 좋음", "반죽 장시간 냉장 휴지"),
    },
    "muffin": {
        "fr": (
            "Faible",
            "Discrète",
            "Moyen",
            "Reste moelleux",
            "Miettes possibles",
        ),
        "en": ("Low", "Mild", "Medium", "Stays moist", "Crumbs likely"),
        "ja": ("低い", "控えめ", "中程度", "しっとり持ち", "食べかすが出やすい"),
        "zh": ("低", "清淡", "中", "口感仍湿润", "可能有碎屑"),
        "ko": ("낮음", "약함", "보통", "촉촉함 유지", "부스러기 있을 수 있음"),
    },
    "muffin-pistache": {
        "fr": (
            "Faible",
            "Discrète",
            "Moyen",
            "Reste moelleux",
            "Miettes possibles",
        ),
        "en": ("Low", "Mild", "Medium", "Stays moist", "Crumbs likely"),
        "ja": ("低い", "控えめ", "中程度", "しっとり持ち", "食べかすが出やすい"),
        "zh": ("低", "清淡", "中", "口感仍湿润", "可能有碎屑"),
        "ko": ("낮음", "약함", "보통", "촉촉함 유지", "부스러기 있을 수 있음"),
    },
    "muffin-pomme-cannelle": {
        "fr": (
            "Faible",
            "Discrète",
            "Moyen",
            "Reste moelleux",
            "Miettes possibles",
        ),
        "en": ("Low", "Mild", "Medium", "Stays moist", "Crumbs likely"),
        "ja": ("低い", "控えめ", "中程度", "しっとり持ち", "食べかすが出やすい"),
        "zh": ("低", "清淡", "中", "口感仍湿润", "可能有碎屑"),
        "ko": ("낮음", "약함", "보통", "촉촉함 유지", "부스러기 있을 수 있음"),
    },
    "gateau-de-savoie": {
        "fr": (
            "Faible",
            "Discrète",
            "Long",
            "Léger, fragile à manipuler",
            "Sucre glace peut s’effriter",
        ),
        "en": (
            "Low",
            "Mild",
            "Long",
            "Light, delicate to handle",
            "Powdered sugar may shed",
        ),
        "ja": (
            "低い",
            "控えめ",
            "長め",
            "軽く崩れやすい",
            "粉糖が落ちやすい",
        ),
        "zh": ("低", "清淡", "长", "轻脆易碎", "糖粉易掉"),
        "ko": ("낮음", "약함", "긴 편", "가볍고 잘 부서짐", "슈가파우더가 떨어질 수 있음"),
    },
    "financier": {
        "fr": (
            "Faible",
            "Discrète",
            "Moyen",
            "Reste compact",
            "Parfum beurre noisette marqué",
        ),
        "en": (
            "Low",
            "Mild",
            "Medium",
            "Stays firm",
            "Pronounced brown-butter aroma",
        ),
        "ja": (
            "低い",
            "控えめ",
            "中程度",
            "形が保ちやすい",
            "焦がしバターの香りが強い",
        ),
        "zh": ("低", "清淡", "中", "形状稳定", "焦黄油香气明显"),
        "ko": ("낮음", "약함", "보통", "형태 유지", "버터 풍미가 진함"),
    },
    "financier-lin": {
        "fr": (
            "Faible",
            "Discrète",
            "Moyen",
            "Reste compact",
            "Parfum beurre noisette marqué",
        ),
        "en": (
            "Low",
            "Mild",
            "Medium",
            "Stays firm",
            "Pronounced brown-butter aroma",
        ),
        "ja": (
            "低い",
            "控えめ",
            "中程度",
            "形が保ちやすい",
            "焦がしバターの香りが強い",
        ),
        "zh": ("低", "清淡", "中", "形状稳定", "焦黄油香气明显"),
        "ko": ("낮음", "약함", "보통", "형태 유지", "버터 풍미가 진함"),
    },
    "cheese-naan": {
        "fr": (
            "Moyen",
            "Marquée",
            "Moyen",
            "Mieux tiède",
            "Fromage filant",
        ),
    },
    "cake-au-citron": {
        "fr": (
            "Faible",
            "Discrète",
            "Moyen",
            "Bon le lendemain",
            "Tranches épaisses pratiques",
        ),
        "en": (
            "Low",
            "Mild",
            "Medium",
            "Good the next day",
            "Thick slices pack well",
        ),
        "ja": (
            "低い",
            "控えめ",
            "中程度",
            "翌日も美味しい",
            "厚切りが持ち運びに便利",
        ),
        "zh": ("低", "清淡", "中", "隔夜仍佳", "厚切便于携带"),
        "ko": ("낮음", "약함", "보통", "다음 날도 좋음", "두껍게 썰면 휴대에 좋음"),
    },
    "banana-cake": {
        "fr": (
            "Faible",
            "Discrète",
            "Moyen",
            "Bon le lendemain",
            "Banane parfum discret",
        ),
        "en": (
            "Low",
            "Mild",
            "Medium",
            "Good the next day",
            "Mild banana aroma",
        ),
        "ja": (
            "低い",
            "控えめ",
            "中程度",
            "翌日も美味しい",
            "バナナの香りは控えめ",
        ),
        "zh": ("低", "清淡", "中", "隔夜仍佳", "香蕉味清淡"),
        "ko": ("낮음", "약함", "보통", "다음 날도 좋음", "바나나 향은 은은함"),
    },
    "carrot-cake": {
        "fr": (
            "Moyen",
            "Discrète",
            "Moyen",
            "Glaçage peut ramollir",
            "Épices et noix parfum discret",
        ),
        "en": (
            "Medium",
            "Mild",
            "Medium",
            "Frosting may soften",
            "Mild spice and nut aroma",
        ),
        "ja": (
            "中程度",
            "控えめ",
            "中程度",
            "フロスティングが軟らかくなることがある",
            "スパイスとナッツの香りは控えめ",
        ),
        "zh": ("中", "清淡", "中", "糖霜可能变软", "香料与坚果味清淡"),
        "ko": ("보통", "약함", "보통", "프로스팅이 무를 수 있음", "향신료·견과 향은 은은함"),
    },
    "empanadas-citron": {
        "fr": (
            "Moyen",
            "Discrète",
            "Moyen",
            "Croûte plus souple après refroidissement",
            "Citron peut humidifier la garniture",
        ),
        "en": (
            "Medium",
            "Mild",
            "Medium",
            "Crust softens when cool",
            "Lemon can moisten filling",
        ),
        "ja": (
            "中程度",
            "控えめ",
            "中程度",
            "冷めると皮がやわらかくなる",
            "レモンで具が湿りやすい",
        ),
        "zh": ("中", "清淡", "中", "冷却后皮变软", "柠檬易让馅变湿"),
        "ko": ("보통", "약함", "보통", "식으면 껍질이 무름", "레몬이 속을 축축하게 할 수 있음"),
    },
    "empanadas-butternut": {
        "fr": (
            "Moyen",
            "Discrète",
            "Moyen",
            "Croûte plus souple après refroidissement",
            "Garniture moelleuse",
        ),
        "en": (
            "Medium",
            "Mild",
            "Medium",
            "Crust softens when cool",
            "Soft filling",
        ),
        "ja": (
            "中程度",
            "控えめ",
            "中程度",
            "冷めると皮がやわらかくなる",
            "具はやわらかめ",
        ),
        "zh": ("中", "清淡", "中", "冷却后皮变软", "馅料柔软"),
        "ko": ("보통", "약함", "보통", "식으면 껍질이 무름", "속은 부드러움"),
    },
    "empanada-intensesalty": {
        "fr": (
            "Moyen",
            "Marquée",
            "Moyen",
            "Croûte plus souple après refroidissement",
            "Saveur très salée",
        ),
        "en": (
            "Medium",
            "Strong",
            "Medium",
            "Crust softens when cool",
            "Very salty filling",
        ),
        "ja": (
            "中程度",
            "強い",
            "中程度",
            "冷めると皮がやわらかくなる",
            "塩味が強い",
        ),
        "zh": ("中", "明显", "中", "冷却后皮变软", "咸味很重"),
        "ko": ("보통", "강함", "보통", "식으면 껍질이 무름", "짠맛이 강함"),
    },
    "dan-bing": {
        "fr": (
            "Faible",
            "Discrète",
            "Rapide",
            "Mieux consommé chaud",
            "Œuf peut encore couler",
        ),
    },
}


def bento_line_count(text: str) -> int:
    parts = re.split(r"^---\s*$", text, flags=re.MULTILINE)
    last = parts[-1].strip()
    return len([ln for ln in last.splitlines() if ln.strip()])


def append_optional(path: Path, slug: str) -> bool:
    lang = path.name.rsplit(".", 2)[1]
    rows = OPTIONAL.get(slug)
    if not rows or lang not in rows:
        return False
    text = path.read_text(encoding="utf-8")
    if bento_line_count(text) != 5:
        return False
    extra = "\n".join(rows[lang])
    # Retirer tout saut de ligne final du fichier puis ajouter bento
    text = text.rstrip() + "\n" + extra + "\n"
    path.write_text(text, encoding="utf-8", newline="\n")
    return True


def main() -> None:
    n = 0
    for fr in sorted(ROOT.glob("*.fr.bentext")):
        slug = fr.name.rsplit(".", 2)[0]
        for lang in LANGS:
            p = ROOT / f"{slug}.{lang}.bentext"
            if not p.exists():
                continue
            if append_optional(p, slug):
                print(f"+ {p.relative_to(ROOT)}")
                n += 1
    print(f"Fichiers mis à jour : {n}")


if __name__ == "__main__":
    main()
