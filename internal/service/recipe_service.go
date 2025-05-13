package service

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "sync"
    "time"
    "strconv"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/algorithm"
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

type RecipeService struct {
    Graph   model.Graph
    TierMap model.TierMap
    Elements []model.Element
}

func NewRecipeService(filePath string) (*RecipeService, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("gagal membaca file data: %v", err)
    }

    var elements []model.Element
    if err := json.Unmarshal(data, &elements); err != nil {
        return nil, fmt.Errorf("gagal unmarshal data: %v", err)
    }

    graph := make(model.Graph)
    tierMap := make(model.TierMap)

    for _, element := range elements {
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
        Elements: elements,
    }, nil
}

func (s *RecipeService) FindRecipes(params model.SearchRequest) model.SearchResponse {
    if params.MaxRecipes <= 0 {
        params.MaxRecipes = 1
    }

    startTime := time.Now()
    var recipes []model.Recipe
    var visitedNodes int

    switch params.Algorithm {
    case "bfs":
        recipes, visitedNodes = algorithm.BFS(s.Graph, s.TierMap, params.TargetElement, !params.MultipleRecipe, params.MaxRecipes)
    case "dfs":
        recipes, visitedNodes = algorithm.DFS(s.Graph, s.TierMap, params.TargetElement, !params.MultipleRecipe, params.MaxRecipes)
    case "bidirectional":
        recipes, visitedNodes = algorithm.Bidirectional(s.Graph, s.TierMap, params.TargetElement, !params.MultipleRecipe, params.MaxRecipes)
    default:
        recipes, visitedNodes = algorithm.BFS(s.Graph, s.TierMap, params.TargetElement, !params.MultipleRecipe, params.MaxRecipes)
    }

    treeData := s.GenerateTreeData(params.TargetElement, recipes)

    return model.SearchResponse{
        Target: params.TargetElement,
        Algorithm: params.Algorithm,
        Time: time.Since(startTime).Seconds(),
        VisitedNodes: visitedNodes,
        Recipes: recipes,
        TreeData: treeData,
    }
}

func (s *RecipeService) FindMultipleRecipes(params model.SearchRequest) model.SearchResponse {
    var wg sync.WaitGroup
    var mutex sync.Mutex
    var totalVisitedNodes int
    var combinedRecipes []model.Recipe
    foundRecipes := make(map[string]bool)

    if !params.MultipleRecipe {
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

            var recipes []model.Recipe
            var visitedNodes int

            switch params.Algorithm {
            case "bfs":
                recipes, visitedNodes = algorithm.BFS(s.Graph, s.TierMap, params.TargetElement, false, maxRecipes)
            case "dfs":
                recipes, visitedNodes = algorithm.DFS(s.Graph, s.TierMap, params.TargetElement, false, maxRecipes)
            case "bidirectional":
                recipes, visitedNodes = algorithm.Bidirectional(s.Graph, s.TierMap, params.TargetElement, false, maxRecipes)
            default:
                recipes, visitedNodes = algorithm.BFS(s.Graph, s.TierMap, params.TargetElement, false, maxRecipes)
            }

            mutex.Lock()
            totalVisitedNodes += visitedNodes
            
            for _, recipe := range recipes {
                recipeKey := s.getRecipeKey(recipe)
                
                if !foundRecipes[recipeKey] {
                    recipe.ID = len(combinedRecipes)
                    combinedRecipes = append(combinedRecipes, recipe)
                    foundRecipes[recipeKey] = true
                    
                    if len(combinedRecipes) >= params.MaxRecipes {
                        break
                    }
                }
            }
            mutex.Unlock()
        }(i)
    }

    wg.Wait()

    if len(combinedRecipes) > params.MaxRecipes {
        combinedRecipes = combinedRecipes[:params.MaxRecipes]
    }

    treeData := s.GenerateTreeData(params.TargetElement, combinedRecipes)

    return model.SearchResponse{
        Target: params.TargetElement,
        Algorithm: params.Algorithm,
        Time: time.Since(startTime).Seconds(),
        VisitedNodes: totalVisitedNodes,
        Recipes: combinedRecipes,
        TreeData: treeData,
    }
}

func (s *RecipeService) getRecipeKey(recipe model.Recipe) string {
    key := ""
    for _, node := range recipe.Nodes {
        if node.Level == 1 {
            key += node.Label + ","
        }
    }
    return key
}

func (s *RecipeService) GenerateTreeData(targetElement string, recipes []model.Recipe) model.TreeNode {
    rootNode := model.TreeNode{
        ID: "root",
        Name: targetElement,
        Combine: []model.TreeNode{},
    }

    for i, recipe := range recipes {
        combineNodes := []model.TreeNode{}
        
        for _, node := range recipe.Nodes {
            if node.Level == 1 {
                combineNode := model.TreeNode{
                    ID: node.ID,
                    Name: node.Label,
                    Combine: s.buildSubtree(recipe, node.ID),
                }
                combineNodes = append(combineNodes, combineNode)
            }
        }
        
        recipeNode := model.TreeNode{
            ID: "recipe-" + strconv.Itoa(i),
            Name: "Recipe " + strconv.Itoa(i+1),
            Combine: combineNodes,
        }
        
        rootNode.Combine = append(rootNode.Combine, recipeNode)
    }

    return rootNode
}

func (s *RecipeService) buildSubtree(recipe model.Recipe, nodeID string) []model.TreeNode {
    result := []model.TreeNode{}
    
    for _, link := range recipe.Links {
        if link.Source == nodeID {
            var targetNode model.RecipeNode
            for _, node := range recipe.Nodes {
                if node.ID == link.Target {
                    targetNode = node
                    break
                }
            }
            
            childNode := model.TreeNode{
                ID: targetNode.ID,
                Name: targetNode.Label,
                Combine: s.buildSubtree(recipe, targetNode.ID),
            }
            
            result = append(result, childNode)
        }
    }
    
    return result
}

func (s *RecipeService) GetAllElements() []string {
    elements := make([]string, 0, len(s.Elements))
    for _, element := range s.Elements {
        elements = append(elements, element.Name)
    }
    return elements
}