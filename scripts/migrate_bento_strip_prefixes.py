#!/usr/bin/env python3
"""Convertit les lignes Bento Préfixe|valeur en valeurs seules (ordre fixe)."""

from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent / "recipes"

KNOWN_PREFIXES = frozenset(
    {
        "Transport",
        "Réchauffage",
        "Froid",
        "Couvert",
        "Manger",
        "Tâches",
        "Fuites",
        "Odeur",
        "TempsPréparation",
        "Veille",
        "Tenue",
        "Notes",
    }
)


def split_sections(content: str) -> list[str]:
    raw = content.replace("\r\n", "\n").replace("\r", "\n")
    parts = [p.strip() for p in raw.split("---")]
    while parts and parts[0] == "":
        parts.pop(0)
    while parts and parts[-1] == "":
        parts.pop()
    return parts


def strip_bento_block(section: str) -> str | None:
    lines = [ln.rstrip() for ln in section.split("\n")]
    nonempty = [ln.strip() for ln in lines if ln.strip()]
    if not nonempty:
        return None
    out: list[str] = []
    changed = False
    for s in nonempty:
        if "|" not in s:
            out.append(s)
            continue
        key, _, rest = s.partition("|")
        k = key.strip()
        if k in KNOWN_PREFIXES:
            out.append(rest.strip())
            changed = True
        else:
            out.append(s)
    if not changed:
        return None
    return "\n".join(out)


def process_file(path: Path) -> bool:
    raw = path.read_text(encoding="utf-8")
    parts = split_sections(raw)
    if len(parts) < 2:
        return False
    new_last = strip_bento_block(parts[-1])
    if new_last is None:
        return False
    parts[-1] = new_last
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
    print(f"Mis à jour : {n} fichier(s)")


if __name__ == "__main__":
    main()
