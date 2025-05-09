package model

type Recipe struct {
    Ingredients []string
}

type Graph map[string][][]string

type Node struct {
    Element   string
    Path      []Recipe
    Visited   map[string]bool
    Depth     int
}

type SearchResult struct {
    TargetElement string   `json:"targetElement"`
    Recipes       []Recipe `json:"recipes"`
    VisitedNodes  int      `json:"visitedNodes"`
    TimeElapsed   float64  `json:"timeElapsed"`
}

type SearchParams struct {
    TargetElement  string `json:"targetElement"`
    Algorithm      string `json:"algorithm"`
    FindShortest   bool   `json:"findShortest"`
    MaxRecipes     int    `json:"maxRecipes"`
}

type ElementListResponse struct {
    Elements []string `json:"elements"`
}