package model

type Element struct {
    Name string `json:"name"`
    Recipe string `json:"recipe"`
    Tier int `json:"tier"`
}

type RecipeNode struct {
    ID string `json:"id"`
    Label string `json:"label"`
    Level int `json:"level"`
}

type RecipeLink struct {
    Source string `json:"source"`
    Target string `json:"target"`
}

type Recipe struct {
    ID int `json:"id"`
    Nodes []RecipeNode `json:"nodes"`
    Links []RecipeLink `json:"links"`
}

type TreeNode struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Combine []TreeNode `json:"combine,omitempty"`
}

type SearchRequest struct {
    TargetElement string `json:"targetElement"`
    Algorithm string `json:"algorithm"`
    MultipleRecipe bool `json:"multipleRecipe"`
    MaxRecipes int `json:"maxRecipes"`
}

type SearchResponse struct {
    Target string `json:"target"`
    Algorithm string `json:"algorithm"`
    Time float64 `json:"time"`
    VisitedNodes int `json:"visitedNodes"`
    Recipes []Recipe `json:"recipes"`
    TreeData TreeNode `json:"treeData,omitempty"`
}

type Graph map[string][][]string

type TierMap map[string]int

type Node struct {
    Element string
    Path []Recipe
    Visited map[string]bool
    Depth int
}