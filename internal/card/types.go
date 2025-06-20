package card

type IngredientConfig struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Rarity  string   `json:"rarity"`
	Type    string   `json:"type"`
	Effects []string `json:"effects"`
	Icon    string   `json:"icon"`
	Copies  int      `json:"copies"`
}

type RecipeConfig struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Requires   []string `json:"requires"`
	Score      int      `json:"score"`
	Difficulty int      `json:"difficulty"`
	Icon       string   `json:"icon"`
}

type IngredientFile struct {
	Ingredients []IngredientConfig `json:"ingredients"`
}

type RecipeFile struct {
	Recipes []RecipeConfig `json:"recipes"`
}
