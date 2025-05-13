package algorithm

import (
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func Bidirectional(graph model.Graph, tierMap model.TierMap, targetElement string, findShortest bool, maxRecipes int) ([]model.Recipe, int) {
    
    baseElements := map[string]bool{
        "Air": true,
        "Earth": true,
        "Fire": true,
        "Water": true,
        "Time": true,
    }
    
    forwardQueue := []struct{
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
    
    backwardQueues := make(map[string][]struct{
        element string
        path [][]string
        visited map[string]bool
        depth int
    })
    
    for baseElement := range baseElements {
        backwardQueues[baseElement] = []struct{
            element string
            path [][]string
            visited map[string]bool
            depth int
        }{
            {
                element: baseElement,
                path: [][]string{},
                visited: make(map[string]bool),
                depth: 0,
            },
        }
    }
    
    forwardVisited := make(map[string]struct{
        path [][]string
        depth int
    })
    backwardVisited := make(map[string]map[string]struct{
        path [][]string
        depth int
    })
    
    for baseElement := range baseElements {
        backwardVisited[baseElement] = make(map[string]struct{
            path [][]string
            depth int
        })
    }
    
    recipes := []model.Recipe{}
    visitedCount := 0
    recipeIdCounter := 0
    
    foundRecipes := make(map[string]bool)
    
    for len(forwardQueue) > 0 && len(recipes) < maxRecipes {
        current := forwardQueue[0]
        forwardQueue = forwardQueue[1:]
        
        if _, exists := forwardVisited[current.element]; exists {
            continue
        }
        
        forwardVisited[current.element] = struct{
            path [][]string
            depth int
        }{
            path: current.path,
            depth: current.depth,
        }
        visitedCount++
        
        for _, bVisited := range backwardVisited {
            if bInfo, meetPoint := bVisited[current.element]; meetPoint {
                forwardPath := current.path
                backwardPath := bInfo.path
                
                completePath := make([][]string, len(forwardPath)+len(backwardPath))
                copy(completePath, forwardPath)
                for i, bPath := range backwardPath {
                    completePath[len(forwardPath)+i] = bPath
                }
                
                recipe := createRecipe(recipeIdCounter, targetElement, completePath, false)
                
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
        
        if baseElements[current.element] {
            continue
        }
        
        combinations := graph[current.element]
        for _, combination := range combinations {
            if isValidCombination(tierMap, current.element, combination) {
                newPath := make([][]string, len(current.path)+1)
                copy(newPath, current.path)
                newPath[len(current.path)] = combination
                
                for _, ingredient := range combination {
                    if _, visited := forwardVisited[ingredient]; !visited && tierMap[ingredient] < tierMap[current.element] {
                        newVisited := make(map[string]bool)
                        for k, v := range current.visited {
                            newVisited[k] = v
                        }
                        newVisited[current.element] = true
                        
                        forwardQueue = append(forwardQueue, struct{
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
        }
    }
    
    if len(recipes) < maxRecipes {
        maxBackwardSteps := 3
        
        for baseElement, queue := range backwardQueues {
            stepCount := 0
            
            for len(queue) > 0 && stepCount < maxBackwardSteps && len(recipes) < maxRecipes {
                stepCount++
                levelSize := len(queue)
                
                for i := 0; i < levelSize; i++ {
                    current := queue[0]
                    queue = queue[1:]
                    
                    if _, exists := backwardVisited[baseElement][current.element]; exists {
                        continue
                    }
                    
                    backwardVisited[baseElement][current.element] = struct{
                        path [][]string
                        depth int
                    }{
                        path: current.path,
                        depth: current.depth,
                    }
                    visitedCount++
                    
                    if fInfo, meetPoint := forwardVisited[current.element]; meetPoint {
                        forwardPath := fInfo.path
                        backwardPath := current.path
                        
                        completePath := make([][]string, len(forwardPath)+len(backwardPath))
                        copy(completePath, forwardPath)
                        for i, bPath := range backwardPath {
                            completePath[len(forwardPath)+i] = bPath
                        }
                        
                        recipe := createRecipe(recipeIdCounter, targetElement, completePath, false)
                        
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
                    
                    for element, combinations := range graph {
                        for _, combination := range combinations {
                            for _, ingredient := range combination {
                                if ingredient == current.element && tierMap[element] > tierMap[current.element] {
                                    newPath := make([][]string, len(current.path)+1)
                                    copy(newPath, current.path)
                                    newPath[len(current.path)] = combination
                                    
                                    queue = append(queue, struct{
                                        element string
                                        path [][]string
                                        visited map[string]bool
                                        depth int
                                    }{
                                        element: element,
                                        path: newPath,
                                        visited: current.visited,
                                        depth: current.depth + 1,
                                    })
                                }
                            }
                        }
                    }
                }
            }
            
            backwardQueues[baseElement] = queue
        }
    }
    
    return recipes, visitedCount
}