package service

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/algorithm"
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

type RecipeService struct {
    Graph   model.Graph
    TierMap model.TierMap
}

func NewRecipeService(filePath string) (*RecipeService, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("gagal membaca file data: %v", err)
    }

    var elementsData []model.ElementData
    if err := json.Unmarshal(data, &elementsData); err != nil {
        return nil, fmt.Errorf("gagal unmarshal data: %v", err)
    }

    graph := make(model.Graph)
    tierMap := make(model.TierMap)

    for _, element := range elementsData {
        tierMap[element.Name] = element.Tier

        if element.Recipe == "Available from the start" {
            graph[element.Name] = [][]string{}
            continue
        }

        var combinations [][]string
        recipeVariants := strings.Split(element.Recipe, ";")
        
        for _, variant := range recipeVariants {
            variant = strings.TrimSpace(variant)
            if variant == "" {
                continue
            }
            
            ingredients := strings.Split(variant, "+")
            if len(ingredients) == 2 {
                combination := []string{
                    strings.TrimSpace(ingredients[0]),
                    strings.TrimSpace(ingredients[1]),
                }
                combinations = append(combinations, combination)
            }
        }
        
        graph[element.Name] = combinations
    }

    return &RecipeService{
        Graph:   graph,
        TierMap: tierMap,
    }, nil
}

func (s *RecipeService) FindRecipes(params model.SearchParams) model.SearchResult {
    if params.MaxRecipes <= 0 {
        params.MaxRecipes = 1
    }

    switch params.Algorithm {
    case "bfs":
        return algorithm.BFS(s.Graph, s.TierMap, params.TargetElement, params.FindShortest, params.MaxRecipes)
    case "dfs":
        return algorithm.DFS(s.Graph, s.TierMap, params.TargetElement, params.FindShortest, params.MaxRecipes)
    case "bidirectional":
        return algorithm.Bidirectional(s.Graph, s.TierMap, params.TargetElement, params.FindShortest, params.MaxRecipes)
    default:
        return algorithm.BFS(s.Graph, s.TierMap, params.TargetElement, params.FindShortest, params.MaxRecipes)
    }
}

func (s *RecipeService) FindMultipleRecipes(params model.SearchParams) model.SearchResult {
    var wg sync.WaitGroup
    var mutex sync.Mutex
    var totalVisitedNodes int
    var combinedResults []model.Recipe

    if params.FindShortest {
        return s.FindRecipes(params)
    }

    numGoroutines := 4
    recipesPerGoroutine := params.MaxRecipes / numGoroutines
    if recipesPerGoroutine < 1 {
        recipesPerGoroutine = 1
    }

    startTime := time.Now()

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()

            maxRecipes := recipesPerGoroutine
            if idx == numGoroutines-1 {
                maxRecipes = params.MaxRecipes - (idx * recipesPerGoroutine)
            }

            if maxRecipes <= 0 {
                return
            }

            var result model.SearchResult
            switch params.Algorithm {
            case "bfs":
                result = algorithm.BFS(s.Graph, s.TierMap, params.TargetElement, false, maxRecipes)
            case "dfs":
                result = algorithm.DFS(s.Graph, s.TierMap, params.TargetElement, false, maxRecipes)
            case "bidirectional":
                result = algorithm.Bidirectional(s.Graph, s.TierMap, params.TargetElement, false, maxRecipes)
            default:
                result = algorithm.BFS(s.Graph, s.TierMap, params.TargetElement, false, maxRecipes)
            }

            mutex.Lock()
            totalVisitedNodes += result.VisitedNodes
            combinedResults = append(combinedResults, result.Recipes...)
            mutex.Unlock()
        }(i)
    }

    wg.Wait()

    if len(combinedResults) > params.MaxRecipes {
        combinedResults = combinedResults[:params.MaxRecipes]
    }

    return model.SearchResult{
        TargetElement: params.TargetElement,
        Recipes:       combinedResults,
        VisitedNodes:  totalVisitedNodes,
        TimeElapsed:   time.Since(startTime).Seconds(),
    }
}

func (s *RecipeService) GetAllElements() []string {
    elements := make([]string, 0, len(s.Graph))
    for element := range s.Graph {
        elements = append(elements, element)
    }
    return elements
}