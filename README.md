# Bentext API

HTTP API in Go (standard library) to parse and serve recipes in the **bentxt** (`.bentext`) plain-text format.

[![Support on Ko-fi](https://img.shields.io/badge/Ko--fi-Support%20me-ff5f5f?logo=kofi&logoColor=white)](https://ko-fi.com/kazerlelutin)

## Run the API

```bash
go run ./cmd/api
```

Server listens on **http://localhost:8080**.

## Endpoints

`GET /` returns a JSON **index** with the full list of routes and query variants (same information as summarized below).

| Method | Route | Description |
| ------ | ----- | ----------- |
| GET | `/` | Route catalog (JSON) |
| GET | `/health` | Health check |
| GET | `/public/` | Static files (recipe images, etc.) |
| GET | `/api/recipes` | All recipes as JSON |
| GET | `/api/recipes?lang=fr` | Filter by language (`fr`, `en`, `ja`, `zh`, `ko`) |
| GET | `/api/recipes?bentext=true` or `?include=bentext` | Same JSON plus raw file in `bentext` (alongside parsed fields, including `identity`) |
| GET | `/api/recipes/{lang}/{slug}` | One recipe, e.g. `/api/recipes/fr/banana-cake` |
| GET | `/api/recipes/{lang}/{slug}?format=json` | Explicit JSON (default if `format` is omitted) |
| GET | `/api/recipes/{lang}/{slug}?format=bentext` | Raw `.bentext` as `text/plain; charset=utf-8` |
| GET | `/api/recipes/{lang}/{slug}?bentext=true` | JSON with an extra `bentext` field (full source) |
| POST | `/api/convert/bentxt` | Body = bentxt text → JSON. Optional query: `lang`, `slug` |
| GET | `/api/ingredients/lookup` | Ingredient icon lookup (`?q=name`) |
| GET | `/api/ingredients/sprite` | Sprite sheet URL and coordinates |

Invalid `format` on a single-recipe request returns **400**; unknown `lang`/`slug` or bad path returns **404**.

## Recipe images

Place images in **`public/`**. The file base name must match the recipe **without** the language suffix:

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

**Optional after steps:** notes (conseils), **tags**, then a **bento** block (one `Prefix|value` per line for transport / reheating / cold chain / how to eat). In the source file, prefixes are fixed (e.g. `Transport`, `Réchauffage`, `Froid`, `Manger`) across all languages; only the text after `|` is localized.

### Section order

| Index | Content | Required |
| ----- | ------- | -------- |
| 0 | Identity (3 lines: name, servings number, description) | yes |
| 1 | Ingredients | yes |
| 2 | Steps | yes |
| 3 | Notes (conseils) | no |
| 4 | Tags | no |
| 5 | Bento (repas emporté) | no |

The parser treats the **last** section as **bento** if its first non-empty line starts with `Transport|`. Otherwise the last section is **tags**, and any sections before that in the tail are **notes** (or tags + bento as above).

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

### Section 2 – Steps

One step per line.

### Sections 3–4 – Notes and tags

- **Notes:** one line per tip.  
- **Tags:** one tag per line.

### Section 5 – Bento (optional)

Lines look like `Transport|Facile`, `Réchauffage|Optionnel (four ou micro-ondes)`, etc. Optional keys include `Fuites`, `Odeur`, `Veille`, `Tenue`, `Notes` (the line prefix `Notes|` maps to JSON `extra_notes`, distinct from recipe **notes**).

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
Transport|Facile
Réchauffage|Optionnel (four ou micro-ondes)
Froid|Non
Manger|À la main ou Couverts
```

## JSON shape (parsed recipe)

Returned by `GET /api/recipes`, `GET /api/recipes/{lang}/{slug}` (JSON mode), and `POST /api/convert/bentxt`:

| Field | Type | Notes |
| ----- | ---- | ----- |
| `slug` | string | Derived from filename |
| `lang` | string | `fr`, `en`, `ja`, `zh`, `ko`, … |
| `identity` | object | `name`, `servings`, `description` |
| `ingredients` | array | `name`, `quantity`, `unit`, `note?`, `alternatives[]`, `icon?` |
| `steps` | string[] | |
| `notes` | string[] | Conseils |
| `tags` | string[] | |
| `bento` | object? | Omitted if absent in file |
| `image` | object? | `url`, `width`, `height` |
| `bentext` | string? | Only when `?bentext=true` or `?include=bentext` |

**`bento`** (when present):

| JSON field | Source prefix in file |
| ---------- | --------------------- |
| `transport` | `Transport` |
| `reheat` | `Réchauffage` |
| `cold` | `Froid` |
| `eating` | `Manger` |
| `leaks` | `Fuites` |
| `smell` | `Odeur` |
| `prep_ahead` | `Veille` |
| `holding` | `Tenue` |
| `extra_notes` | `Notes` |

## Features

- **CORS:** `Access-Control-Allow-Origin: *` on all routes  
- **Rate limiting:** 100 requests per minute per IP (429 when exceeded)

## Build

```bash
go build -o bentext-api ./cmd/api
./bentext-api
```

## Project layout

```
bentext/
├── cmd/
│   ├── api/            # HTTP server entrypoint
│   └── check-recipes/  # CLI: language coverage per recipe slug
├── internal/
│   ├── handler/        # HTTP handlers (recipes, convert, ingredients, …)
│   ├── recipe/         # Bentxt parsing (identity … bento)
│   └── ingredients/    # Icon lookup & sprite metadata
├── public/             # Static assets (recipe images)
├── recipes/            # *.bentext (e.g. *.fr.bentext)
├── scripts/            # Optional tooling (e.g. batch edits)
├── go.mod
└── README.md
```
