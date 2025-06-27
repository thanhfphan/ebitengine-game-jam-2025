package card

type IngredientConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type RecipeConfig struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Requires []string `json:"requires"`
	Icon     string   `json:"icon"`
}

type IngredientFile struct {
	Ingredients []IngredientConfig `json:"ingredients"`
}

type RecipeFile struct {
	Recipes []RecipeConfig `json:"recipes"`
}
