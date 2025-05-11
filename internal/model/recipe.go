package model

import (
    "encoding/json"
)

type Recipe struct {
    Ingredients []string `json:"ingredients"`
}

type Node struct {
    Element   string
    Path      []Recipe
    Visited   map[string]bool
    Depth     int
}

type Graph map[string][][]string

type TierMap map[string]int

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

type ElementData struct {
    Name   string `json:"name"`
    Recipe string `json:"recipe"`
    Tier   int    `json:"tier"`
}

func (r Recipe) MarshalJSON() ([]byte, error) {
    if r.Ingredients == nil {
        r.Ingredients = []string{}
    }
    return json.Marshal(struct {
        Ingredients []string `json:"ingredients"`
    }{
        Ingredients: r.Ingredients,
    })
}