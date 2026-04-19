#!/usr/bin/env python3
"""
Normalise la ligne réchauffage (bloc Bento) sur toutes les recettes :
- deux modes : « Optionnel four ~ micro-ondes » (et équivalents), comme les alternatives d’ingrédients ;
- un seul mode : « Optionnel » / « Optional » / … (échelle courte pour le tri).
"""
from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent

# Ordre : chaînes sources les plus longues d’abord.
REPL: list[tuple[str, str]] = [
    # Four + micro-ondes
    ("Optionnel (four ou micro-ondes)", "Optionnel four ~ micro-ondes"),
    ("Optional (oven or microwave)", "Optional oven ~ microwave"),
    ("任意（オーブンまたは電子レンジ）", "任意 オーブン ~ 電子レンジ"),
    ("可选（烤箱或微波炉）", "可选 烤箱 ~ 微波炉"),
    ("선택(오븐 또는 전자레인지)", "선택 오븐 ~ 전자레인지"),
    # Four + poêle (naan)
    ("Optionnel four ou poêle", "Optionnel four ~ poêle"),
    ("Optional (oven or pan)", "Optional oven ~ pan"),
    ("任意（オーブンまたはフライパン）", "任意 オーブン ~ フライパン"),
    ("可选烤箱或平底锅", "可选 烤箱 ~ 平底锅"),
    ("선택(오븐 또는 프라이팬)", "선택 오븐 ~ 프라이팬"),
    # Micro seul
    ("Optionnel micro-ondes", "Optionnel"),
    ("Optional microwave", "Optional"),
    ("任意（電子レンジ）", "任意"),
    ("可选微波炉", "可选"),
    ("선택(전자레인지)", "선택"),
    # Four seul (pas « Optionnel four » ici : après « four ~ … », ça casserait la ligne ; les fiches four-only sont passées en « Optionnel » court.)
    ("Optional (oven)", "Optional"),
    ("任意（オーブン）", "任意"),
    ("可选烤箱", "可选"),
    ("선택(오븐)", "선택"),
    # Poêle seul
    ("Optionnel poêle", "Optionnel"),
    ("Optional (pan)", "Optional"),
    ("任意（フライパン）", "任意"),
    ("可选平底锅", "可选"),
    ("선택(프라이팬)", "선택"),
]


def main() -> None:
    n = 0
    for folder in ("recipes", "draft_recipes"):
        d = ROOT / folder
        if not d.is_dir():
            continue
        for p in sorted(d.glob("*.bentext")):
            t = p.read_text(encoding="utf-8")
            orig = t
            for a, b in REPL:
                t = t.replace(a, b)
            if t != orig:
                p.write_text(t, encoding="utf-8", newline="\n")
                n += 1
                print(p.name)
    # Réparation si une ancienne version remplaçait « Optionnel four » dans « Optionnel four ~ … »
    repair = [
        ("Optionnel ~ micro-ondes", "Optionnel four ~ micro-ondes"),
        ("Optionnel ~ poêle", "Optionnel four ~ poêle"),
    ]
    for folder in ("recipes", "draft_recipes"):
        d = ROOT / folder
        if not d.is_dir():
            continue
        for p in sorted(d.glob("*.bentext")):
            t = p.read_text(encoding="utf-8")
            orig = t
            for a, b in repair:
                t = t.replace(a, b)
            if t != orig:
                p.write_text(t, encoding="utf-8", newline="\n")
                print("repair:", p.name)
    print(f"Fichiers modifiés : {n}")


if __name__ == "__main__":
    main()
