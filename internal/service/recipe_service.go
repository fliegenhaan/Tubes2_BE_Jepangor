package service

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"
    "time"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/algorithm"
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

type RecipeService struct {
    Graph model.Graph
}

func NewRecipeService(filePath string) (*RecipeService, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read recipes data: %v", err)
    }

    var graph model.Graph
    if err := json.Unmarshal(data, &graph); err != nil {
        return nil, fmt.Errorf("failed to unmarshal recipes data: %v", err)
    }

    return &RecipeService{
        Graph: graph,
    }, nil
}

func (s *RecipeService) FindRecipes(params model.SearchParams) model.SearchResult {
    if params.MaxRecipes <= 0 {
        params.MaxRecipes = 1
    }

    switch params.Algorithm {
    case "bfs":
        return algorithm.BFS(s.Graph, params.TargetElement, params.FindShortest, params.MaxRecipes)
    case "dfs":
        return algorithm.DFS(s.Graph, params.TargetElement, params.FindShortest, params.MaxRecipes)
    case "bidirectional":
        return algorithm.Bidirectional(s.Graph, params.TargetElement, params.FindShortest, params.MaxRecipes)
    default:
        return algorithm.BFS(s.Graph, params.TargetElement, params.FindShortest, params.MaxRecipes)
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
                result = algorithm.BFS(s.Graph, params.TargetElement, false, maxRecipes)
            case "dfs":
                result = algorithm.DFS(s.Graph, params.TargetElement, false, maxRecipes)
            case "bidirectional":
                result = algorithm.Bidirectional(s.Graph, params.TargetElement, false, maxRecipes)
            default:
                result = algorithm.BFS(s.Graph, params.TargetElement, false, maxRecipes)
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