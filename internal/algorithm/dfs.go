package algorithm

import (
    "sync"
    "time"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func DFS(graph model.Graph, targetElement string, findShortest bool, maxRecipes int) model.SearchResult {
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
    
    var mutex sync.Mutex
    resultCount := 0
    
    visited := make(map[string]bool)
    
    var dfsRecursive func(currentElement string, path []model.Recipe, currentVisited map[string]bool, depth int)
    dfsRecursive = func(currentElement string, path []model.Recipe, currentVisited map[string]bool, depth int) {
        mutex.Lock()
        if resultCount >= maxRecipes && !findShortest {
            mutex.Unlock()
            return
        }
        mutex.Unlock()
        
        visitedCount++
        visited[currentElement] = true
        currentVisited[currentElement] = true
        
        if currentElement == targetElement && depth > 0 {
            mutex.Lock()
            if findShortest {
                if resultCount == 0 {
                    results = append(results, path...)
                    resultCount++
                }
            } else {
                results = append(results, path...)
                resultCount++
            }
            mutex.Unlock()
            return
        }
        
        for nextElement, recipes := range graph {
            if visited[nextElement] || (resultCount >= maxRecipes && !findShortest) {
                continue
            }
            
            for _, recipe := range recipes {
                if len(recipe) != 2 {
                    continue
                }

                canMake := false
                requiredElement := ""
                
                if recipe[0] == currentElement {
                    canMake = true
                    requiredElement = recipe[1]
                } else if recipe[1] == currentElement {
                    canMake = true
                    requiredElement = recipe[0]
                }
                
                if canMake && currentVisited[requiredElement] {
                    newPath := make([]model.Recipe, len(path))
                    copy(newPath, path)
                    newRecipe := model.Recipe{Ingredients: []string{currentElement, requiredElement}}
                    newPath = append(newPath, newRecipe)
                    
                    newVisited := make(map[string]bool)
                    for k, v := range currentVisited {
                        newVisited[k] = v
                    }
                    
                    dfsRecursive(nextElement, newPath, newVisited, depth+1)
                }
            }
        }
    }
    
    for _, element := range baseElements {
        if resultCount >= maxRecipes && !findShortest {
            break
        }
        
        currentVisited := map[string]bool{element: true}
        dfsRecursive(element, []model.Recipe{}, currentVisited, 0)
    }
    
    return model.SearchResult{
        TargetElement: targetElement,
        Recipes:       results,
        VisitedNodes:  visitedCount,
        TimeElapsed:   time.Since(startTime).Seconds(),
    }
}