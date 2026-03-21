# Bentext API

Simple HTTP API in Go (standard library) to parse and serve recipes in the **bentxt** text format.

[![Support on Ko-fi](https://img.shields.io/badge/Ko--fi-Support%20me-ff5f5f?logo=kofi&logoColor=white)](https://ko-fi.com/kazerlelutin)

## Run the API

```bash
go run ./cmd/api
```

Server listens on **http://localhost:8080**.

## Endpoints

| Method | Route                 | Description                                 |
| ------ | --------------------- | ------------------------------------------- |
| GET    | `/`                   | List of all available routes                |
| GET    | `/health`             | Health check (status ok)                    |
| GET    | `/api/recipes`        | List all recipes (`?lang=fr` to filter; `?bentext=true` or `?include=bentext` to add raw file as `bentext` next to parsed fields including `identity`) |
| GET    | `/api/recipes/{lang}/{slug}` | One recipe: JSON by default (`?format=json`), or raw file (`?format=bentext`); `?bentext=true` embeds raw text in JSON |
| POST   | `/api/convert/bentxt` | Convert bentxt text (request body) to JSON  |
| GET    | `/public/*`           | Static files (e.g. recipe images)           |

## Recipe images

Place images in the **`public/`** folder. Each image must use the **same base name as the recipe file, without the language suffix**:

- Recipe: `recipes/onigiri-kimchi-mozza.fr.bentext` → image: `public/onigiri-kimchi-mozza.jpg` (or `.png`, `.gif`)
- Recipe: `recipes/soupe-oignon.bentext` → image: `public/soupe-oignon.jpg`

Supported formats: `.jpg`, `.jpeg`, `.png`, `.gif`. In **GET /api/recipes**, each recipe gets an optional `image` object with full URL and dimensions:

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
curl -X POST http://localhost:8080/api/convert/bentxt -H "Content-Type: text/plain" -d "Recipe name
4
Short description.
---
flour|200|g
---
Step 1.
---
tag1"
```

## Bentxt format (`.bentext`)

Bentxt is a plain-text recipe format. The file is split into **sections** by the separator `---` (three hyphens on their own line). You need **at least 3 sections**: identity, ingredients, and steps. After that come optional **notes** (conseils), **tags**, and a final **bento** block (`Transport|…`, `Réchauffage|…`, etc.).

### Section order

| Section index | Content                              | Required |
| ------------- | ------------------------------------ | -------- |
| 0             | Identity                             | yes      |
| 1             | Ingredients                          | yes      |
| 2             | Steps                                | yes      |
| 3             | Notes (conseils)                     | no       |
| 4             | Tags                                 | no       |
| 5             | Bento (repas emporté, paires préfixe–valeur) | no       |

The parser detects the **bento** block when its first line starts with `Transport|`. **Tags** are always the section immediately before bento, or the last section if there is no bento.

So (examples):

- **3 sections:** identity, ingredients, steps.
- **4 sections:** identity, ingredients, steps, tags (no notes, no bento).
- **5 sections:** identity, ingredients, steps, tags, bento — or notes + tags without bento.
- **6 sections:** identity, ingredients, steps, notes, tags, bento.

In **JSON** (`GET /api/recipes`, `GET /api/recipes/...`, `POST /api/convert/bentxt`), parsed bento fields appear under **`bento`** (e.g. `transport`, `reheat`, `cold`, `eating`, plus optional `leaks`, `smell`, `prep_ahead`, `holding`, `extra_notes`).

### Section 0 – Identity

Three lines (order matters):

1. **Recipe name** (first line).
2. **Servings** – a number (e.g. `4`).
3. **Description** – free text (can be multiple lines; the last non-numeric line is used as description in the current parser).

### Section 1 – Ingredients

One ingredient per line. Each line can be:

- **Simple:** `name|quantity|unit|note`
  - `quantity` and `unit` are optional (defaults: `1`, `piece`).
  - `note` is optional.
  - Example: `flour|200|g` or `eggs|2|piece|room temperature`.

- **With alternatives:** use `~` to separate the main ingredient from alternatives.
  - Format: `main|qty|unit|note ~ alt1|qty|unit|note ~ alt2|...`
  - Example: `yogurt|120|ml ~ milk|120|ml` (yogurt or milk).

### Section 2 – Steps

One step per line. Each line is a single step (no special syntax).

### Section 3 & 4 – Notes and tags

- **Notes:** one note per line (free text).
- **Tags:** one tag per line (e.g. `breakfast`, `vegan`).

### Example bentxt file

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
```

### JSON output shape

Parsed recipes are returned as JSON with this structure:

- `id`, `lang`
- `identity`: `name`, `servings`, `description`
- `ingredients`: array of `{ name, quantity, unit, note?, alternatives? }`
- `steps`: array of strings
- `notes`: array of strings
- `tags`: array of strings
- `image` (optional): `{ url, width, height }` — full image URL and dimensions when a matching file exists in `public/`

## Features

- **CORS**: permissive (`Access-Control-Allow-Origin: *`) for all routes
- **Rate limiting**: 100 requests per minute per IP (429 when exceeded)

## Build

```bash
go build -o bentext.exe ./cmd/api
./bentext.exe
```

## Project layout (Go conventions)

```
bentext/
├── cmd/api/            # Application entrypoint
│   └── main.go
├── internal/handler/   # HTTP handlers
│   ├── home.go
│   ├── health.go
│   ├── recipes.go
│   ├── convert.go
│   └── image.go        # Recipe image lookup & dimensions
├── internal/recipe/    # Bentxt parsing
│   ├── recipe.go
│   └── parse.go
├── public/             # Static files (recipe images, same base name as .bentext)
├── recipes/            # Recipe files (.bentext, .fr.bentext, etc.)
├── go.mod
└── README.md
```
