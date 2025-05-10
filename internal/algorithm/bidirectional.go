package algorithm

import (
    "container/list"
    "time"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func Bidirectional(graph model.Graph, targetElement string, findShortest bool, maxRecipes int) model.SearchResult {
    startTime := time.Now()
    visitedCount := 0
    results := []model.Recipe{}
    
    baseElements := []string{"Air", "Earth", "Fire", "Water"}
    
    if _, exists := graph[targetElement]; !exists {
        return model.SearchResult{
            TargetElement: targetElement,
            Recipes:       results,
            VisitedNodes:  visitedCount,
            TimeElapsed:   time.Since(startTime).Seconds(),
        }
    }
    
    for _, base := range baseElements {
        if targetElement == base {
            return model.SearchResult{
                TargetElement: targetElement,
                Recipes:       []model.Recipe{},
                VisitedNodes:  1,
                TimeElapsed:   time.Since(startTime).Seconds(),
            }
        }
    }
    
    forwardQueue := list.New()
    backwardQueue := list.New()
    
    forwardVisited := make(map[string]model.Node)
    backwardVisited := make(map[string]model.Node)
    
    for _, element := range baseElements {
        node := model.Node{
            Element:  element,
            Path:     []model.Recipe{},
            Visited:  map[string]bool{element: true},
            Depth:    0,
        }
        forwardQueue.PushBack(node)
        forwardVisited[element] = node
    }
    
    targetNode := model.Node{
        Element:  targetElement,
        Path:     []model.Recipe{},
        Visited:  map[string]bool{targetElement: true},
        Depth:    0,
    }
    backwardQueue.PushBack(targetNode)
    backwardVisited[targetElement] = targetNode
    
    for forwardQueue.Len() > 0 && backwardQueue.Len() > 0 && len(results) < maxRecipes {
        currentLevelSize := forwardQueue.Len()
        for i := 0; i < currentLevelSize; i++ {
            current := forwardQueue.Remove(forwardQueue.Front()).(model.Node)
            visitedCount++
            
            if backNode, exists := backwardVisited[current.Element]; exists {
                combinedPath := combinePaths(current.Path, backNode.Path)
                results = append(results, combinedPath...)
                
                if findShortest || len(results) >= maxRecipes {
                    return model.SearchResult{
                        TargetElement: targetElement,
                        Recipes:       results,
                        VisitedNodes:  visitedCount,
                        TimeElapsed:   time.Since(startTime).Seconds(),
                    }
                }
                continue
            }
            
            exploreNeighbors(graph, current, forwardQueue, forwardVisited, true)
        }
        
        currentLevelSize = backwardQueue.Len()
        for i := 0; i < currentLevelSize; i++ {
            current := backwardQueue.Remove(backwardQueue.Front()).(model.Node)
            visitedCount++

            if forwardNode, exists := forwardVisited[current.Element]; exists {
                combinedPath := combinePaths(forwardNode.Path, current.Path)
                results = append(results, combinedPath...)
                
                if findShortest || len(results) >= maxRecipes {
                    return model.SearchResult{
                        TargetElement: targetElement,
                        Recipes:       results,
                        VisitedNodes:  visitedCount,
                        TimeElapsed:   time.Since(startTime).Seconds(),
                    }
                }
                continue
            }
            
            exploreNeighbors(graph, current, backwardQueue, backwardVisited, false)
        }
    }
    
    return model.SearchResult{
        TargetElement: targetElement,
        Recipes:       results,
        VisitedNodes:  visitedCount,
        TimeElapsed:   time.Since(startTime).Seconds(),
    }
}

func exploreNeighbors(
    graph model.Graph, 
    current model.Node, 
    queue *list.List, 
    visited map[string]model.Node,
    isForward bool,
) {
    if isForward {
        for nextElement, recipes := range graph {
            if _, exists := visited[nextElement]; exists {
                continue
            }
            
            for _, recipe := range recipes {
                if len(recipe) != 2 {
                    continue
                }
                
                canMake := false
                requiredElement := ""
                
                if recipe[0] == current.Element {
                    canMake = true
                    requiredElement = recipe[1]
                } else if recipe[1] == current.Element {
                    canMake = true
                    requiredElement = recipe[0]
                }
                
                if canMake && current.Visited[requiredElement] {
                    newPath := make([]model.Recipe, len(current.Path))
                    copy(newPath, current.Path)
                    newRecipe := model.Recipe{Ingredients: []string{current.Element, requiredElement}}
                    newPath = append(newPath, newRecipe)
                    
                    newVisited := make(map[string]bool)
                    for k, v := range current.Visited {
                        newVisited[k] = v
                    }
                    newVisited[nextElement] = true
                    
                    newNode := model.Node{
                        Element:  nextElement,
                        Path:     newPath,
                        Visited:  newVisited,
                        Depth:    current.Depth + 1,
                    }
                    
                    queue.PushBack(newNode)
                    visited[nextElement] = newNode
                }
            }
        }
    } else {
        for ingredient, recipes := range graph {
            if _, exists := visited[ingredient]; exists {
                continue
            }
            
            for _, recipe := range recipes {
                if len(recipe) != 2 {
                    continue
                }
                
                for _, elem := range recipe {
                    if elem == current.Element {
                        newPath := make([]model.Recipe, len(current.Path))
                        copy(newPath, current.Path)
                        newRecipe := model.Recipe{Ingredients: []string{ingredient, current.Element}}
                        newPath = append(newPath, newRecipe)
                        
                        newVisited := make(map[string]bool)
                        for k, v := range current.Visited {
                            newVisited[k] = v
                        }
                        newVisited[ingredient] = true
                        
                        newNode := model.Node{
                            Element:  ingredient,
                            Path:     newPath,
                            Visited:  newVisited,
                            Depth:    current.Depth + 1,
                        }
                        
                        queue.PushBack(newNode)
                        visited[ingredient] = newNode
                        break
                    }
                }
            }
        }
    }
}

func combinePaths(forwardPath, backwardPath []model.Recipe) []model.Recipe {
    resultPath := make([]model.Recipe, len(forwardPath))
    copy(resultPath, forwardPath)
    
    for i := len(backwardPath) - 1; i >= 0; i-- {
        resultPath = append(resultPath, backwardPath[i])
    }
    
    return resultPath
}