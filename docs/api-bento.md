# API — objet `bento` (repas emporté)

Les réponses JSON de `GET /api/recipes`, `GET /api/recipes/{lang}/{slug}` et `POST /api/convert/bentxt` peuvent inclure un objet **`bento`** lorsque le dernier bloc du fichier `.bentext` est reconnu comme section Bento (valeurs seules ou ancien format `Préfixe|valeur`).

Les chaînes sont **localisées** comme dans le fichier source (une langue par fichier). Pour badges, filtres ou validation côté client, utilisez les **tableaux de correspondance** ci-dessous ; ils sont alignés sur [`internal/bento/vocab.go`](../internal/bento/vocab.go).

## Schéma JSON

| Champ | Type | Ligne `.bentext` (valeurs seules) | Rôle |
| ----- | ---- | -------------------------------- | ---- |
| `transport` | `string` | 1 | Facilité de transport (souvent une échelle + courte précision) |
| `reheat` | `string` | 2 | Réchauffage (échelle + détail). Plusieurs modes possibles : même idée que les **alternatives d’ingrédients**, séparer avec **` ~ `** (ex. `Recommandé vapeur ~ micro-ondes`). |
| `cold` | `string` | 3 | Conservation au froid / chaîne du froid. Alternatives (ex. froid / ambiant, délai) : **` ~ `** comme pour le réchauffage. |
| `utensils` | `string` | 4 | **Ustensiles / modalité** : faut-il des couverts, lesquels (main, baguettes, couverts, combinaisons avec **` ~ `**). Remplace l’ancienne paire `cover` + `eating`. |
| `stains` | `string` | 5 (optionnel) | Risque de taches |
| `smell` | `string` | 6 (optionnel) | Odeur probable en boîte fermée |
| `prep_time` | `string` | 7 (optionnel) | Échelle de **durée de préparation** |
| `holding` | `string` | 8 (optionnel) | Tenue après plusieurs heures |
| `extra_notes` | `string` | 9 (optionnel) | Notes libres (cas limites) |

**Bloc fichier :** **4 lignes minimum** (sans `|`), jusqu’à **9 lignes**. Même nombre de lignes et même ordre entre les fichiers de langue d’une même recette.

**Réchauffage (`reheat`) :** pour lister des **alternatives** (vapeur, four, micro-ondes, etc.), utiliser le séparateur **` ~ `** comme entre ingrédients (`a ~ b`), éventuellement après un mot d’échelle (`Optionnel`, `Recommandé`, …). Le parseur n’interdit pas `~` dans le bloc Bento (seul `|` est réservé aux ingrédients).

**Froid (`cold`) :** formes **courtes** (ex. `Au frais`, `Frais si délai`, `Chilled if delay`). Pour **deux consignes équivalentes** (ex. froid / ambiant selon garniture), utiliser aussi **` ~ `** (ex. `Frais ~ ambiant`, `Fresh ~ room temp`).

**Ustensiles (`utensils`) :** une seule ligne qui combine besoin de couverts et modalité (ex. `À la main ~ Couverts`, `Baguettes ~ Couverts`, `À la main`). Pour **deux options** équivalentes, utiliser **` ~ `** (même convention que les ingrédients et le réchauffage).

## Intégration front

- **Affichage direct :** tout champ peut être rendu tel quel pour l’UI dans la langue courante.
- **Échelles / badges :** pour `transport`, `utensils`, `stains`, `smell`, `prep_time`, comparer la chaîne reçue aux lignes du tableau de correspondance (ou normaliser par `Key` si vous dupliquez la table en TypeScript).
- **Évolutions :** des codes normalisés additionnels (ex. `transport_ease`) pourraient compléter les chaînes plus tard ; le contrat actuel reste des `string` localisées.

## Correspondances de valeurs (référence)

### Transport (facilité)

| id | fr | en | ja | zh | ko |
| -- | -- | -- | -- | -- | -- |
| transport_easy | Facile | Easy | 簡単 | 容易 | 쉬움 |
| transport_medium | Moyen | Medium | やや注意 | 中等 | 보통 |
| transport_delicate | Délicat | Delicate | 壊れやすい | 易损 | 다소 까다로움 |

### Réchauffage (exemples)

| id | fr | en | ja | zh | ko |
| -- | -- | -- | -- | -- | -- |
| reheat_no | Non | No | 不要 | 不需要 | 불필요 |
| reheat_optional | Optionnel | Optional | 任意 | 可选 | 선택 |
| reheat_optional_oven_micro | Optionnel four ~ micro-ondes | Optional oven ~ microwave | 任意 オーブン ~ 電子レンジ | 可选 烤箱 ~ 微波炉 | 선택 오븐 ~ 전자레인지 |
| reheat_optional_oven_pan | Optionnel four ~ poêle | Optional oven ~ pan | 任意 オーブン ~ フライパン | 可选 烤箱 ~ 平底锅 | 선택 오븐 ~ 프라이팬 |
| reheat_recommended_steam_micro | Recommandé vapeur ~ micro-ondes | Recommended steam ~ microwave | 推奨 蒸し ~ 電子レンジ | 建议 蒸 ~ 微波 | 권장 찜 ~ 전자레인지 |

### Conservation au froid (`cold`, exemples)

| id | fr | en | ja | zh | ko |
| -- | -- | -- | -- | -- | -- |
| cold_no | Non | No | 不要 | 否 | 아니오 |
| cold_keep_short | Au frais | Keep cold | 要冷蔵 | 需冷藏 | 냉장 |
| cold_chilled_short | Au frais | Chilled | 冷蔵 | 冷藏 | 냉장 |
| cold_delay | Frais si délai | Chilled if delay | 遅延 ~ 冷蔵 | 久置 ~ 冷藏 | 지연 ~ 냉장 |
| cold_fresh_ambient | Frais ~ ambiant | Fresh ~ room temp | 冷蔵 ~ 常温 | 冷藏 ~ 常温 | 냉장 ~ 실온 |

### Ustensiles (`utensils`)

| id | fr | en | ja | zh | ko |
| -- | -- | -- | -- | -- | -- |
| utensils_hand_cutlery | À la main ~ Couverts | By hand ~ cutlery | 手づかみ ~ カトラリー | 手抓 ~ 餐具 | 손으로 ~ 수저 |
| utensils_hand_chopsticks | À la main ~ Baguettes | Hand ~ chopsticks | 手づかみ ~ 箸 | 手抓 ~ 筷子 | 손 ~ 젓가락 |
| utensils_chopsticks_cutlery | Baguettes ~ Couverts | Chopsticks ~ cutlery | 箸 ~ カトラリー | 筷子 ~ 餐具 | 젓가락 ~ 수저 |
| utensils_cutlery | Couverts | Cutlery | カトラリー | 餐具 | 수저 |
| utensils_chopsticks | Baguettes | Chopsticks | 箸 | 筷子 | 젓가락 |
| utensils_hand | À la main | By hand | 手づかみ | 手抓 | 손으로 |

### Taches (`stains`)

| id | fr | en | ja | zh | ko |
| -- | -- | -- | -- | -- | -- |
| stains_none | Non | None | なし | 无 | 없음 |
| stains_low | Faible | Low | 低い | 低 | 낮음 |
| stains_medium | Moyen | Medium | 中程度 | 中 | 보통 |
| stains_high | Élevé | High | 高い | 高 | 높음 |

### Odeur (`smell`)

| id | fr | en | ja | zh | ko |
| -- | -- | -- | -- | -- | -- |
| smell_discrete | Discrète | Mild | 控えめ | 清淡 | 약함 |
| smell_marked | Marquée | Strong | 強い | 明显 | 강함 |

### Temps de préparation (`prep_time`)

| id | fr | en | ja | zh | ko |
| -- | -- | -- | -- | -- | -- |
| prep_quick | Rapide | Quick | 短時間 | 快 | 짧음 |
| prep_medium | Moyen | Medium | 中程度 | 中 | 보통 |
| prep_long | Long | Long | 長め | 长 | 긴 편 |

## Compatibilité

Les champs JSON **`leaks`** et **`prep_ahead`** ne sont plus émis ; ils sont remplacés par **`stains`** et **`prep_time`**. Les champs **`cover`** et **`eating`** ne sont plus émis ; le fichier utilise une seule ligne / **`utensils`**.

Le parseur accepte encore l’**ancien format fichier** avec **deux lignes** (couvert + manger) aux rangs 4–5, ainsi que l’ancien format `Préfixe|valeur` avec clés `Couvert` et `Manger`, et fusionne en **`utensils`**.
