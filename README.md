# Bentext API

HTTP API in Go (standard library) to parse and serve recipes in the **bentxt** (`.bentext`) plain-text format.

[Support on Ko-fi](https://ko-fi.com/kazerlelutin)

## Run the API

```bash
joutée : après échec du matc
```

Server listens on **[http://localhost:8080](http://localhost:8080)**.joutée : après échec du matc

## Endpoints

`GET /` returns a JSON **index** with the full list of routes and query variants (same information as summarized below).


| Method | Route                                             | Description                                                                          |
| ------ | ------------------------------------------------- | ------------------------------------------------------------------------------------ |
| GET    | `/`                                               | Route catalog (JSON)                                                                 |
| GET    | `/health`                                         | Health check                                                                         |
| GET    | `/public/`                                        | Static files (recipe images, etc.)                                                   |
| GET    | `/api/recipes`                                    | All recipes as JSON                                                                  |
| GET    | `/api/recipes?lang=fr`                            | Filter by language (`fr`, `en`, `ja`, `zh`, `ko`)                                    |
| GET    | `/api/recipes?bentext=true` or `?include=bentext` | Same JSON plus raw file in `bentext` (alongside parsed fields, including `identity`) |
| GET    | `/api/recipes/{lang}/{slug}`                      | One recipe, e.g. `/api/recipes/fr/banana-cake`                                       |
| GET    | `/api/recipes/{lang}/{slug}?format=json`          | Explicit JSON (default if `format` is omitted)                                       |
| GET    | `/api/recipes/{lang}/{slug}?format=bentext`       | Raw `.bentext` as `text/plain; charset=utf-8`                                        |
| GET    | `/api/recipes/{lang}/{slug}?bentext=true`         | JSON with an extra `bentext` field (full source)                                     |
| POST   | `/api/convert/bentxt`                             | Body = bentxt text → JSON. Optional query: `lang`, `slug`                            |
| GET    | `/api/ingredients/lookup`                         | Ingredient icon lookup (`?q=name`)                                                   |
| GET    | `/api/ingredients/sprite`                         | Sprite sheet URL and coordinates                                                     |


Invalid `format` on a single-recipe request returns **400**; unknown `lang`/`slug` or bad path returns **404**.

## Recipe images

Place images in `**public/`**. The file base name must match the recipe **without** the language suffix:

- Recipe: `recipes/onigiri-kimchi-mozza.fr.bentext` → `public/onigiri-kimchi-mozza.jpg` (or `.png`, `.gif`)

Supported: `.jpg`, `.jpeg`, `.png`, `.gif`. When a file matches, the JSON includes an optional `image` object:

```json
"image": {
  "url": "http://localhost:8080/public/onigiri-kimchi-mozza.jpg",
  "width": 1200,
  "height": 800
}
```

## Examples

```bash
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/api/recipes
curl "http://localhost:8080/api/recipes?lang=fr"
curl "http://localhost:8080/api/recipes?bentext=true"
curl "http://localhost:8080/api/recipes/fr/banana-cake"
curl "http://localhost:8080/api/recipes/fr/banana-cake?format=bentext"
curl "http://localhost:8080/api/recipes/fr/banana-cake?bentext=true"
curl "http://localhost:8080/api/ingredients/lookup?q=flour"
curl "http://localhost:8080/api/ingredients/sprite"
curl -X POST "http://localhost:8080/api/convert/bentxt?lang=en&slug=demo" \
  -H "Content-Type: text/plain" -d "Recipe name
4
Short description.
---
flour|200|g
---
Step 1.
---
tag1"
```

You can also send a file as the body (e.g. `--data-binary @recipes/banana-cake.en.bentext` on Unix shells).

## Bentxt format (`.bentext`)

Files are split into **sections** by a line containing only `---` (three hyphens).

**Minimum:** identity, ingredients, steps.

**Optional after steps:** notes (conseils), **tags**, then a **bento** block: **one value per line**, **no labels and no `|`** in that block. Lines are **fully translated** per language file. The API maps them to JSON keys (`transport`, `reheat`, `cold`, `utensils`, …). **Order is fixed** (see below). Legacy lines `Préfixe|valeur` (e.g. `Transport|…`) are still accepted by the parser but should not be used in new content.

**API — bloc Bento :** schéma détaillé et tableaux de correspondance multilingue des valeurs standard : [docs/api-bento.md](docs/api-bento.md).

### Section order


| Index | Content                                                | Required |
| ----- | ------------------------------------------------------ | -------- |
| 0     | Identity (3 lines: name, servings number, description) | yes      |
| 1     | Ingredients                                            | yes      |
| 2     | Steps                                                  | yes      |
| 3     | Notes (conseils)                                       | no       |
| 4     | Tags                                                   | no       |
| 5     | Bento (repas emporté)                                  | no       |


The **last** section is treated as **bento** when it matches the legacy `Préfixe|valeur` format, or when it has **4–10 non-empty lines with no `|`** (value-only; **4–9** lines is the current lunch-box model). Otherwise the last section is **tags**, and earlier tail sections are **notes** (or tags + bento as in the examples). When there are only two tail sections, the first block is classified as **tags** vs **notes** using simple heuristics (length, word count, punctuation).

Examples of block counts (identity + ingredients + steps + …):

- **3 blocks:** identity, ingredients, steps only.
- **4 blocks:** … + tags (no notes, no bento).
- **5 blocks:** … + tags + bento, or … + notes + tags (no bento).
- **6 blocks:** … + notes + tags + bento.

### Section 0 – Identity

1. Recipe name
2. Servings (integer)
3. Description (one line in typical files; the parser consumes further identity lines with the existing rules)

### Section 1 – Ingredients

- **Simple:** `name|quantity|unit|note` — `quantity` / `unit` / `note` optional (defaults: `1`, `piece`).
- **Alternatives:** `main|qty|unit ~ alt|qty|unit ~ …`

Ingredient names in JSON may get an optional `icon` (`x`, `y` in the ingredient sprite) when resolved via `/api/ingredients`.

### Ingredient sprite sheet

Icons come from the `**ingredient-sprites.bentext`** file at the repository root (same format for the running API and Docker image). Each block is separated by a line `---`. The first line of a block is `row col` (grid cell in the 32×32 sprite sheet), then one or more **alias** lines.

- **Several names for one sprite:** use `|` on the same line. Example: `beurre|beurre doux` maps both names to the same cell. Spaces around `|` are ignored.
- **No exact alias:** the resolver tries **leading word sequences** from the full name (normalized: trim + lowercase), longest match first. Example: `beurre doux` uses the sprite registered for `beurre` if `beurre doux` is not listed. If both `citron` and `citron vert` exist, `citron vert` picks the more specific entry.
- **Still no match:** it picks the registered alias with smallest **Levenshtein** distance if that distance is **≤ 2**. The query and candidate must each have at least **4 runes**, so very short names skip this step (fewer false positives). If two aliases tie on distance, the **longer** alias wins.

Synonyms that are **not** close in spelling or as a leading word sequence (e.g. `margarine` vs `beurre`) must still be listed explicitly with `|`.

`GET /api/ingredients/sprite` exposes coordinates keyed by **explicit** aliases only (including every segment after splitting on `|`). Prefix and Levenshtein fallbacks apply when enriching recipe JSON and `GET /api/ingredients/lookup?q=…`.

`**public/ingredient-sprites.csv`** is a human-readable grid of the same layout (optional `|` in a cell, e.g. `beurre|beurre doux`); it is not read by the server. Edit `**ingredient-sprites.bentext**` for behavior changes.

### Section 2 – Steps

One step per line.

### Sections 3–4 – Notes and tags

- **Notes:** one line per tip.  
- **Tags:** one tag per line.

### Section 5 – Bento (optional)

**4 required lines** (in order): transport ease, reheating, cold chain / storage, **utensils** (cutlery need and how to eat in **one** line, e.g. by hand, chopsticks, `hand ~ cutlery`). **Optional lines 5–9** (same order if used): stain risk (`stains`), smell in a closed box, **prep time scale** (`prep_time`), holding quality after several hours, extra note (JSON `extra_notes`, not recipe **notes**). Do not use `|` in this section.

### Example `.bentext` file

```text
Chocolate muffins
9
Chocolate chip muffins.

---
flour|260|g
sugar|150|g
baking powder|11|g
eggs|2|piece|room temperature
yogurt|120|ml ~ milk|120|ml
oil|100|ml
chocolate|150|g
---
Mix dry ingredients in a bowl.
Mix wet ingredients in another bowl.
Combine briefly, add chocolate.
Rest 30 minutes. Bake 5 min at 220°C, then 15 min at 180°C.
---
Resting the batter is important.
---
baking
sweet
---
Easy
Optional oven ~ microwave
No
By hand ~ cutlery
```

## JSON shape (parsed recipe)

Returned by `GET /api/recipes`, `GET /api/recipes/{lang}/{slug}` (JSON mode), and `POST /api/convert/bentxt`:


| Field         | Type     | Notes                                                          |
| ------------- | -------- | -------------------------------------------------------------- |
| `slug`        | string   | Derived from filename                                          |
| `lang`        | string   | `fr`, `en`, `ja`, `zh`, `ko`, …                                |
| `identity`    | object   | `name`, `servings`, `description`                              |
| `ingredients` | array    | `name`, `quantity`, `unit`, `note?`, `alternatives[]`, `icon?` |
| `steps`       | string[] |                                                                |
| `notes`       | string[] | Conseils                                                       |
| `tags`        | string[] |                                                                |
| `bento`       | object?  | Omitted if absent in file                                      |
| `image`       | object?  | `url`, `width`, `height`                                       |
| `bentext`     | string?  | Only when `?bentext=true` or `?include=bentext`                |


`**bento**` (when present): fields are filled from **line order** in the file (or from legacy `Préfixe|valeur` lines if present). See [docs/api-bento.md](docs/api-bento.md) for canonical value sets per language.


| JSON field    | Source line (value-only file) |
| ------------- | ----------------------------- |
| `transport`   | line 1                        |
| `reheat`      | line 2                        |
| `cold`        | line 3                        |
| `utensils`    | line 4                        |
| `stains`      | line 5 (optional)             |
| `smell`       | line 6 (optional)             |
| `prep_time`   | line 7 (optional)             |
| `holding`     | line 8 (optional)             |
| `extra_notes` | line 9 (optional)             |


## Features

- **CORS:** `Access-Control-Allow-Origin:` * on all routes  
- **Rate limiting:** 100 requests per minute per IP (429 when exceeded)  
- **Bento (lunch box) metadata:** [docs/api-bento.md](docs/api-bento.md)

## Build

```bash
go build -o bentext-api ./cmd/api
./bentext-api
```

## Security (CI)

On GitHub, **[`.github/workflows/ci.yml`](.github/workflows/ci.yml)** runs `go vet`, tests, **[govulncheck](https://go.dev/doc/security/vuln/)** (known vulnerabilities affecting your code and the standard library), and **[Trivy](https://github.com/aquasecurity/trivy)** on the repository filesystem and on the Docker image built from [`Dockerfile`](Dockerfile). Scans use severity **HIGH** and **CRITICAL** and **`ignore-unfixed: true`** to reduce noise from distro advisories without patches.

**[Dependabot](https://docs.github.com/en/code-security/dependabot)** ([`.github/dependabot.yml`](.github/dependabot.yml)) opens weekly PRs for Go modules, base Docker images, and GitHub Actions used in workflows.

## Project layout

```
bentext/
├── .github/
│   ├── workflows/      # CI: tests, govulncheck, Trivy
│   └── dependabot.yml  # Weekly dependency updates
├── cmd/
│   ├── api/            # HTTP server entrypoint
│   └── check-recipes/  # CLI: language coverage per recipe slug
├── docs/               # API docs (e.g. Bento field map)
├── internal/
│   ├── bento/          # Canonical vocabulary rows (doc alignment)
│   ├── handler/        # HTTP handlers (recipes, convert, ingredients, …)
│   ├── recipe/         # Bentxt parsing (identity … bento)
│   └── ingredients/    # Icon lookup & sprite metadata
├── public/             # Static assets (recipe images, ingredient-sprites.csv layout reference)
├── ingredient-sprites.bentext  # Sprite grid aliases (section Ingredient sprite sheet)
├── recipes/            # *.bentext (e.g. *.fr.bentext)
├── scripts/            # Optional tooling (e.g. batch edits, migrate_bento_v2.py)
├── go.mod
└── README.md
```

