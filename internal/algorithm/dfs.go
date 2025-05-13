package algorithm

import (
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func DFS(graph model.Graph, tierMap model.TierMap, targetElement string, findShortest bool, maxRecipes int) ([]model.Recipe, int) {
    
    stack := []struct{
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
    
    for len(stack) > 0 && len(recipes) < maxRecipes {
        lastIdx := len(stack) - 1
        current := stack[lastIdx]
        stack = stack[:lastIdx]
        
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
        for i := len(combinations) - 1; i >= 0; i-- {
            combination := combinations[i]
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
                            
                            stack = append(stack, struct{
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