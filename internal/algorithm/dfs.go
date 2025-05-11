package algorithm

import (
    "time"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func DFS(graph model.Graph, tierMap model.TierMap, targetElement string, findShortest bool, maxRecipes int) model.SearchResult {
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
    
    var stack []model.Node
    visited := make(map[string]bool)
    var recipes []model.Recipe
    visitedCount := 0

    for _, element := range baseElements {
        node := model.Node{
            Element:   element,
            Path:      []model.Recipe{},
            Visited:   make(map[string]bool),
            Depth:     0,
        }
        node.Visited[element] = true
        stack = append(stack, node)
    }

    for len(stack) > 0 && len(recipes) < maxRecipes {
        lastIndex := len(stack) - 1
        currNode := stack[lastIndex]
        stack = stack[:lastIndex]
        
        visitedCount++
        
        if currNode.Depth > maxDepth {
            continue
        }
        
        if currNode.Element == targetElement {
            recipes = append(recipes, currNode.Path...)
            
            if findShortest {
                break
            }
            continue
        }

        if !visited[currNode.Element] {
            visited[currNode.Element] = true
            
            type candidate struct {
                Element, OtherIngredient, NextElement string
            }
            
            var candidates []candidate
            
            for nextElement, combinations := range graph {
                nextTier, exists := tierMap[nextElement]
                if !exists || (nextTier >= targetTier && nextElement != targetElement) {
                    continue
                }
                
                for _, combination := range combinations {
                    if len(combination) == 2 {
                        ingredient1, ingredient2 := combination[0], combination[1]
                        
                        if ingredient1 == currNode.Element || ingredient2 == currNode.Element {
                            otherIngredient := ingredient2
                            if ingredient1 != currNode.Element {
                                otherIngredient = ingredient1
                            }
                            
                            candidates = append(candidates, candidate{
                                Element:         currNode.Element,
                                OtherIngredient: otherIngredient,
                                NextElement:     nextElement,
                            })
                        }
                    }
                }
            }
            
            for i := len(candidates) - 1; i >= 0; i-- {
                cand := candidates[i]
                
                if currNode.Visited[cand.OtherIngredient] {
                    continue
                }
                
                newNode := model.Node{
                    Element:   cand.NextElement,
                    Path:      make([]model.Recipe, len(currNode.Path)+1),
                    Visited:   make(map[string]bool),
                    Depth:     currNode.Depth + 1,
                }
                
                copy(newNode.Path, currNode.Path)
                for k, v := range currNode.Visited {
                    newNode.Visited[k] = v
                }
                
                newNode.Visited[cand.NextElement] = true
                newNode.Visited[cand.OtherIngredient] = true
                
                newNode.Path[len(currNode.Path)] = model.Recipe{
                    Ingredients: []string{cand.Element, cand.OtherIngredient},
                }
                
                stack = append(stack, newNode)
            }
        }
    }

    return model.SearchResult{
        TargetElement: targetElement,
        Recipes:       recipes,
        VisitedNodes:  visitedCount,
        TimeElapsed:   time.Since(startTime).Seconds(),
    }
}