package handler

import (
	"encoding/json"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Bentext API",
		"routes": []map[string]string{
			{"method": "GET", "path": "/", "description": "Index JSON : liste de toutes les routes ci-dessous"},
			{"method": "GET", "path": "/health", "description": "Vérification de disponibilité du service"},
			{"method": "GET", "path": "/public/", "description": "Fichiers statiques (images de recettes, etc.)"},
			// Recettes — liste
			{"method": "GET", "path": "/api/recipes", "description": "Liste des recettes au format JSON (toutes langues)"},
			{"method": "GET", "path": "/api/recipes?lang=fr", "description": "Liste filtrée par langue (fr, en, ja, zh, ko)"},
			{"method": "GET", "path": "/api/recipes?bentext=true", "description": "Liste JSON + champ bentext (fichier brut) pour chaque recette ; identity déjà en JSON"},
			{"method": "GET", "path": "/api/recipes?include=bentext", "description": "Équivalent à bentext=true"},
			{"method": "GET", "path": "/api/recipes?lang=fr&bentext=true", "description": "Filtre langue + inclusion du bentext brut"},
			// Recettes — une par lang/slug (même handler que /api/recipes/…)
			{"method": "GET", "path": "/api/recipes/{lang}/{slug}", "description": "Une recette en JSON (ex. /api/recipes/fr/banana-cake)"},
			{"method": "GET", "path": "/api/recipes/{lang}/{slug}?format=json", "description": "Une recette en JSON (explicite, défaut si format absent)"},
			{"method": "GET", "path": "/api/recipes/{lang}/{slug}?format=bentext", "description": "Une recette en texte brut (text/plain), fichier .bentext tel quel"},
			{"method": "GET", "path": "/api/recipes/{lang}/{slug}?bentext=true", "description": "Une recette JSON avec en plus le champ bentext (texte complet du fichier)"},
			{"method": "GET", "path": "/api/recipes/{lang}/{slug}?include=bentext", "description": "Équivalent à bentext=true sur une recette"},
			// Autres API
			{"method": "POST", "path": "/api/convert/bentxt", "description": "Convertir un bentext (corps de la requête) en JSON ; query optionnelles : lang, slug"},
			{"method": "GET", "path": "/api/ingredients/lookup", "description": "Icône d’ingrédient (?q=nom)"},
			{"method": "GET", "path": "/api/ingredients/sprite", "description": "URL de la planche de sprites + coordonnées byAlias"},
		},
	})
}
