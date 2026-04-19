#!/usr/bin/env python3
"""Raccourcit les lignes Bento `cold` et forme les alternatives `eating` avec « ~ »."""

from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
DIRS = (ROOT / "recipes", ROOT / "draft_recipes")

# Ordre : chaînes longues d’abord (évite les remplacements partiels).
REPLACEMENTS: dict[str, list[tuple[str, str]]] = {
    "fr": [
        (
            "À manger frais ou à température ambiante selon garniture",
            "Frais ~ ambiant",
        ),
        ("Conserver au frais si délai long", "Frais si délai"),
        ("Conserver au frais jusqu'à consommation", "Au frais"),
        ("Conserver au frais si garniture fraîche", "Frais si garniture"),
        ("Conserver au frais", "Au frais"),
        ("Baguettes ou Couverts", "Baguettes ~ Couverts"),
        ("À la main ou Baguettes", "À la main ~ Baguettes"),
        ("Couverts ou à la main", "Couverts ~ main"),
    ],
    "en": [
        (
            "Eat fresh or at room temperature depending on filling",
            "Fresh ~ room temp",
        ),
        ("Keep chilled if using fresh toppings", "Chilled if fresh toppings"),
        ("Keep chilled if long delay", "Chilled if delay"),
        ("Keep chilled until eaten", "Keep cold"),
        ("Keep chilled", "Chilled"),
        ("Chopsticks or cutlery", "Chopsticks ~ cutlery"),
        ("By hand or chopsticks", "Hand ~ chopsticks"),
        ("Cutlery or by hand", "Cutlery ~ hand"),
    ],
    "ja": [
        ("具材に応じて冷蔵のまままたは常温で", "冷蔵 ~ 常温"),
        ("時間が空く場合は冷蔵", "遅延 ~ 冷蔵"),
        ("食べるまで冷蔵", "要冷蔵"),
        ("箸またはカトラリー", "箸 ~ カトラリー"),
        ("手づかみまたは箸", "手づかみ ~ 箸"),
        ("カトラリーまたは手づかみ", "カトラリー ~ 手づかみ"),
    ],
    "zh": [
        ("视馅料冷藏或常温食用", "冷藏 ~ 常温"),
        ("间隔较久需冷藏", "久置 ~ 冷藏"),
        ("食用前冷藏", "需冷藏"),
        ("筷子或餐具", "筷子 ~ 餐具"),
        ("手抓或筷子", "手抓 ~ 筷子"),
        ("餐具或手抓", "餐具 ~ 手抓"),
    ],
    "ko": [
        ("속재료에 따라 냉장 또는 실온", "냉장 ~ 실온"),
        ("시간 간격이 길면 냉장", "지연 ~ 냉장"),
        ("먹을 때까지 냉장", "냉장"),
        ("젓가락 또는 수저", "젓가락 ~ 수저"),
        ("손으로 또는 젓가락", "손 ~ 젓가락"),
        ("수저 또는 손으로", "수저 ~ 손"),
    ],
}


def patch_file(path: Path) -> bool:
    suffix = path.suffix.removeprefix(".")  # .fr.bentext -> wrong
    parts = path.name.rsplit(".", 2)
    if len(parts) != 3:
        return False
    lang = parts[1]
    pairs = REPLACEMENTS.get(lang)
    if not pairs:
        return False
    text = path.read_text(encoding="utf-8")
    new = text
    for old, repl in pairs:
        new = new.replace(old, repl)
    if new != text:
        path.write_text(new, encoding="utf-8", newline="\n")
        return True
    return False


def main() -> None:
    n = 0
    for d in DIRS:
        if not d.is_dir():
            continue
        for path in sorted(d.glob("*.bentext")):
            if patch_file(path):
                n += 1
                print(path.relative_to(ROOT))
    print(f"Fichiers modifiés : {n}")


if __name__ == "__main__":
    main()
