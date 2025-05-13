package algorithm

import (
    "strconv"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func BFS(graph model.Graph, tierMap model.TierMap, targetElement string, findShortest bool, maxRecipes int) ([]model.Recipe, int) {
    
    queue := []struct{
        element string
        path [][]string
        visited map[string]bool
        depth int
    }{
        {
            element: targetElement,
            path: [][]string{},
            visited: make(map[string]bool),
            depth: 0,
        },
    }
    
    visited := make(map[string]bool)
    recipes := []model.Recipe{}
    visitedCount := 0
    recipeIdCounter := 0
    
    foundRecipes := make(map[string]bool)
    
    baseElements := map[string]bool{
        "Air": true,
        "Earth": true,
        "Fire": true,
        "Water": true,
        "Time": true,
    }
    
    for len(queue) > 0 && len(recipes) < maxRecipes {
        current := queue[0]
        queue = queue[1:]
        
        if visited[current.element] {
            continue
        }
        
        visited[current.element] = true
        visitedCount++
        
        if baseElements[current.element] {
            recipe := createRecipe(recipeIdCounter, targetElement, current.path, true)
            
            recipeKey := getRecipeKey(recipe)
            
            if !foundRecipes[recipeKey] {
                recipes = append(recipes, recipe)
                foundRecipes[recipeKey] = true
                recipeIdCounter++
                
                if findShortest && len(recipes) > 0 {
                    return recipes, visitedCount
                }
            }
            
            continue
        }
        
        combinations := graph[current.element]
        for _, combination := range combinations {
            if isValidCombination(tierMap, current.element, combination) {
                newPath := make([][]string, len(current.path)+1)
                copy(newPath, current.path)
                newPath[len(current.path)] = combination
                
                reachedBase := true
                for _, ingredient := range combination {
                    if !baseElements[ingredient] {
                        reachedBase = false
                        if !visited[ingredient] && tierMap[ingredient] < tierMap[current.element] {
                            newVisited := make(map[string]bool)
                            for k, v := range current.visited {
                                newVisited[k] = v
                            }
                            newVisited[current.element] = true
                            
                            queue = append(queue, struct{
                                element string
                                path [][]string
                                visited map[string]bool
                                depth int
                            }{
                                element: ingredient,
                                path: newPath,
                                visited: newVisited,
                                depth: current.depth + 1,
                            })
                        }
                    }
                }
                
                if reachedBase {
                    recipe := createRecipe(recipeIdCounter, targetElement, newPath, true)
                    
                    recipeKey := getRecipeKey(recipe)
                    
                    if !foundRecipes[recipeKey] {
                        recipes = append(recipes, recipe)
                        foundRecipes[recipeKey] = true
                        recipeIdCounter++
                        
                        if findShortest && len(recipes) > 0 {
                            return recipes, visitedCount
                        }
                    }
                }
            }
        }
    }
    
    return recipes, visitedCount
}

func getRecipeKey(recipe model.Recipe) string {
    key := ""
    for _, node := range recipe.Nodes {
        if node.Level == 1 {
            key += node.Label + ","
        }
    }
    return key
}

func isValidCombination(tierMap model.TierMap, element string, combination []string) bool {
    for _, ingredient := range combination {
        if tierMap[ingredient] >= tierMap[element] {
            return false
        }
    }
    return true
}

func createRecipe(id int, targetElement string, path [][]string, isBase bool) model.Recipe {
    recipe := model.Recipe{
        ID: id,
        Nodes: []model.RecipeNode{},
        Links: []model.RecipeLink{},
    }
    
    rootNode := model.RecipeNode{
        ID: "target-" + strconv.Itoa(id),
        Label: targetElement,
        Level: 0,
    }
    recipe.Nodes = append(recipe.Nodes, rootNode)
    
    for i, combination := range path {
        level := i + 1
        
        for j, ingredient := range combination {
            nodeID := "ingredient-" + strconv.Itoa(id) + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(j)
            
            node := model.RecipeNode{
                ID: nodeID,
                Label: ingredient,
                Level: level,
            }
            recipe.Nodes = append(recipe.Nodes, node)
            var parentID string
            if i == 0 {
                parentID = "target-" + strconv.Itoa(id)
            } else {
                parentID = "ingredient-" + strconv.Itoa(id) + "-" + strconv.Itoa(i-1) + "-0"
            }
            
            link := model.RecipeLink{
                Source: parentID,
                Target: nodeID,
            }
            recipe.Links = append(recipe.Links, link)
        }
    }
    
    return recipe
}