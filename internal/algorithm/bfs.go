package algorithm

import (
    "container/list"
    "sync"
    "time"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func BFS(graph model.Graph, targetElement string, findShortest bool, maxRecipes int) model.SearchResult {
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
    
    queue := list.New()
    visited := make(map[string]bool)
    
    for _, element := range baseElements {
        node := model.Node{
            Element:  element,
            Path:     []model.Recipe{},
            Visited:  map[string]bool{element: true},
            Depth:    0,
        }
        queue.PushBack(node)
    }
    
    var mutex sync.Mutex
    
    for queue.Len() > 0 && (len(results) < maxRecipes || findShortest) {
        current := queue.Remove(queue.Front()).(model.Node)
        visitedCount++
        
        if visited[current.Element] {
            continue
        }
        visited[current.Element] = true
        
        for nextElement, recipes := range graph {
            if visited[nextElement] {
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
                
                if canMake {
                    if current.Visited[requiredElement] {
                        newPath := make([]model.Recipe, len(current.Path))
                        copy(newPath, current.Path)
                        newRecipe := model.Recipe{Ingredients: []string{current.Element, requiredElement}}
                        newPath = append(newPath, newRecipe)
                        
                        newVisited := make(map[string]bool)
                        for k, v := range current.Visited {
                            newVisited[k] = v
                        }
                        newVisited[nextElement] = true
                        
                        if nextElement == targetElement {
                            mutex.Lock()
                            results = append(results, newPath...)
                            mutex.Unlock()
                            
                            if findShortest {
                                return model.SearchResult{
                                    TargetElement: targetElement,
                                    Recipes:       []model.Recipe{newRecipe},
                                    VisitedNodes:  visitedCount,
                                    TimeElapsed:   time.Since(startTime).Seconds(),
                                }
                            }
                            
                            if len(results) >= maxRecipes {
                                break
                            }
                        } else {
                            queue.PushBack(model.Node{
                                Element:  nextElement,
                                Path:     newPath,
                                Visited:  newVisited,
                                Depth:    current.Depth + 1,
                            })
                        }
                    }
                }
            }
            
            if len(results) >= maxRecipes && !findShortest {
                break
            }
        }
    }
    
    return model.SearchResult{
        TargetElement: targetElement,
        Recipes:       results,
        VisitedNodes:  visitedCount,
        TimeElapsed:   time.Since(startTime).Seconds(),
    }
}