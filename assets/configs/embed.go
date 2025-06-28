package configs

import (
    _ "embed"
)

var (
    //go:embed decks/default/ingredients.json
    DefaultIngredientsJSON []byte

    //go:embed decks/default/recipes.json
    DefaultRecipesJSON []byte
)