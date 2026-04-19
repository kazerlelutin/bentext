#!/usr/bin/env python3
"""Migre le bloc Bento vers le format v2 : ligne « couvert », stains/prep_time.

- Ancien format valeur seule : 4 lignes (base) ou 5–9 lignes (avec optionnelles sans couvert).
- Nouveau format : 5–10 lignes (transport, reheat, cold, cover, eating, stains, smell, prep_time, holding, extra_notes).

Insère « couvert » par défaut selon la langue du fichier ; convertit l’ancienne ligne « veille » (prep_ahead) en échelle prep_time."""

from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent / "recipes"

COVER_DEFAULT = {"fr": "Optionnel", "en": "Optional", "ja": "任意", "zh": "可选", "ko": "선택"}

# Ancienne veille → nouvelle échelle durée (approximation)
PREP_FROM_AHEAD_FR = {
    "oui": "Moyen",
    "partielle": "Moyen",
    "non": "Rapide",
}
PREP_FROM_AHEAD_EN = {"yes": "Medium", "partial": "Medium", "no": "Quick"}
PREP_FROM_AHEAD_JA = {"はい": "中程度", "一部": "中程度", "いいえ": "短時間"}
PREP_FROM_AHEAD_ZH = {"是": "中", "部分": "中", "否": "快"}
PREP_FROM_AHEAD_KO = {"예": "보통", "부분": "보통", "아니오": "짧음"}


def split_sections(content: str) -> list[str]:
    raw = content.replace("\r\n", "\n").replace("\r", "\n")
    parts = [p.strip() for p in raw.split("---")]
    while parts and parts[0] == "":
        parts.pop(0)
    while parts and parts[-1] == "":
        parts.pop()
    return parts


def lang_from_filename(name: str) -> str:
    m = re.match(r"^.+\.(fr|en|ja|zh|ko)\.bentext$", name)
    return m.group(1) if m else "fr"


def non_empty_lines(block: str) -> list[str]:
    return [ln.strip() for ln in block.split("\n") if ln.strip()]


def is_short_cover_token(line: str) -> bool:
    s = line.strip()
    if len(s) > 24:
        return False
    low = s.lower()
    if low in ("non", "oui", "optionnel", "no", "yes", "optional"):
        return True
    return s in (
        "不要",
        "任意",
        "必要",
        "不需要",
        "可选",
        "需要",
        "불필요",
        "선택",
        "필요",
    )


def convert_prep_line(line: str, lang: str) -> str:
    raw = line.strip()
    low = raw.lower()
    if lang == "fr":
        k = low.split()[0] if low else ""
        if k in PREP_FROM_AHEAD_FR:
            return PREP_FROM_AHEAD_FR[k]
        if raw in ("Rapide", "Moyen", "Long"):
            return raw
    if lang == "en":
        if low in PREP_FROM_AHEAD_EN:
            return PREP_FROM_AHEAD_EN[low]
        if raw in ("Quick", "Medium", "Long"):
            return raw
    if lang == "ja" and raw in PREP_FROM_AHEAD_JA:
        return PREP_FROM_AHEAD_JA[raw]
    if lang == "zh" and raw in PREP_FROM_AHEAD_ZH:
        return PREP_FROM_AHEAD_ZH[raw]
    if lang == "ko" and raw in PREP_FROM_AHEAD_KO:
        return PREP_FROM_AHEAD_KO[raw]
    # Déjà une échelle ou texte libre
    if lang == "fr" and raw in ("Rapide", "Moyen", "Long"):
        return raw
    # Heuristique : mots courts anciens FR
    if lang == "fr":
        if low.startswith("oui"):
            return "Moyen"
        if low.startswith("partielle"):
            return "Moyen"
        if low.startswith("non"):
            return "Rapide"
    if lang == "en":
        if low.startswith("yes"):
            return "Medium"
        if low.startswith("partial"):
            return "Medium"
        if low.startswith("no"):
            return "Quick"
    return raw


def migrate_bento_lines(lines: list[str], lang: str) -> list[str] | None:
    n = len(lines)
    if n < 4:
        return None
    if any("|" in ln for ln in lines):
        return None
    if n == 10 and is_short_cover_token(lines[3]):
        return None  # déjà format v2 plein
    if n >= 5 and n <= 9 and is_short_cover_token(lines[3]):
        # 5–9 lignes avec couvert en ligne 4 : ne migrer que prep si besoin
        out = list(lines)
        if n >= 8:
            out[7] = convert_prep_line(out[7], lang)
        return out if out != lines else None
    if n >= 5 and is_short_cover_token(lines[3]):
        return None

    cover = COVER_DEFAULT.get(lang, "Optionnel")
    if n == 4:
        out = [lines[0], lines[1], lines[2], cover, lines[3]]
        return out

    # Ancien 5–9 : insérer cover à l’index 3
    out = lines[:3] + [cover] + lines[3:]
    # prep_time était ligne 7 en ancien 1-based = index 6 ; après insert = index 7
    if len(out) > 7:
        out[7] = convert_prep_line(out[7], lang)
    return out


def process_file(path: Path) -> bool:
    raw = path.read_text(encoding="utf-8")
    parts = split_sections(raw)
    if len(parts) < 2:
        return False
    lang = lang_from_filename(path.name)
    last = parts[-1]
    lines = non_empty_lines(last)
    new_lines = migrate_bento_lines(lines, lang)
    if new_lines is None:
        return False
    parts[-1] = "\n".join(new_lines)
    new_body = "\n---\n".join(parts)
    if not new_body.endswith("\n"):
        new_body += "\n"
    path.write_text(new_body, encoding="utf-8", newline="\n")
    return True


def main() -> None:
    n = 0
    for p in sorted(ROOT.glob("*.bentext")):
        if process_file(p):
            n += 1
            print(p.name)
    print(f"Migrés : {n} fichier(s)")


if __name__ == "__main__":
    main()
