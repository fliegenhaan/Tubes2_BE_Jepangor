package algorithm

import (
    "time"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func Bidirectional(graph model.Graph, tierMap model.TierMap, targetElement string, findShortest bool, maxRecipes int) model.SearchResult {
    startTime := time.Now()
    
    if _, exists := graph[targetElement]; !exists {
        return model.SearchResult{
            TargetElement: targetElement,
            Recipes:       []model.Recipe{},
            VisitedNodes:  0,
            TimeElapsed:   time.Since(startTime).Seconds(),
        }
    }

    baseElements := []string{"Air", "Earth", "Fire", "Water"}
    
    targetTier, exists := tierMap[targetElement]
    if !exists {
        targetTier = 999
    }
    
    for _, base := range baseElements {
        if base == targetElement {
            recipe := model.Recipe{
                Ingredients: []string{targetElement},
            }
            return model.SearchResult{
                TargetElement: targetElement,
                Recipes:       []model.Recipe{recipe},
                VisitedNodes:  1,
                TimeElapsed:   time.Since(startTime).Seconds(),
            }
        }
    }
    
    forwardFrontier := make(map[string]model.Node)
    forwardVisited := make(map[string]bool)
    
    backwardFrontier := make(map[string]model.Node)
    
    for _, element := range baseElements {
        node := model.Node{
            Element:   element,
            Path:      []model.Recipe{},
            Visited:   make(map[string]bool),
            Depth:     0,
        }
        node.Visited[element] = true
        forwardFrontier[element] = node
    }
    
    targetNode := model.Node{
        Element:   targetElement,
        Path:      []model.Recipe{},
        Visited:   make(map[string]bool),
        Depth:     0,
    }
    
    targetNode.Visited[targetElement] = true
    backwardFrontier[targetElement] = targetNode
    
    var meetingPoint string
    var forwardNode  model.Node
    visitedCount := 0
    maxIterations := 20
    
    for iter := 0; iter < maxIterations; iter++ {
        if len(forwardFrontier) == 0 || len(backwardFrontier) == 0 {
            break
        }
        
        visitedCount += len(forwardFrontier) + len(backwardFrontier)
        
        for elem := range forwardFrontier {
            if _, exists := backwardFrontier[elem]; exists {
                meetingPoint = elem
                forwardNode = forwardFrontier[elem]
                break
            } 
        }

        
        
        if meetingPoint != "" {
            break
        }
        
        newForwardFrontier := make(map[string]model.Node)
        for elem, node := range forwardFrontier {
            forwardVisited[elem] = true
            
            for nextElement, combinations := range graph {
                nextTier, exists := tierMap[nextElement]
                if !exists || (nextTier >= targetTier && nextElement != targetElement) {
                    continue
                }
                
                for _, combination := range combinations {
                    if len(combination) == 2 {
                        ingredient1, ingredient2 := combination[0], combination[1]
                        
                        if ingredient1 == elem || ingredient2 == elem {
                            otherIngredient := ingredient2
                            if ingredient1 != elem {
                                otherIngredient = ingredient1
                            }
                            
                            if node.Visited[otherIngredient] {
                                continue
                            }
                            
                            newNode := model.Node{
                                Element:   nextElement,
                                Path:      make([]model.Recipe, len(node.Path)+1),
                                Visited:   make(map[string]bool),
                                Depth:     node.Depth + 1,
                            }
                            
                            copy(newNode.Path, node.Path)
                            for k, v := range node.Visited {
                                newNode.Visited[k] = v
                            }
                            
                            newNode.Visited[nextElement] = true
                            newNode.Visited[otherIngredient] = true
                            
                            newNode.Path[len(node.Path)] = model.Recipe{
                                Ingredients: []string{elem, otherIngredient},
                            }
                            
                            if !forwardVisited[nextElement] {
                                newForwardFrontier[nextElement] = newNode
                            }
                        }
                    }
                }
            }
        }
        
        forwardFrontier = newForwardFrontier
    }
    
    if meetingPoint == "" {
        return BFS(graph, tierMap, targetElement, findShortest, maxRecipes)
    }
    
    var combinedPath []model.Recipe
    
    combinedPath = append(combinedPath, forwardNode.Path...)
    
    return model.SearchResult{
        TargetElement: targetElement,
        Recipes:       combinedPath,
        VisitedNodes:  visitedCount,
        TimeElapsed:   time.Since(startTime).Seconds(),
    }
}