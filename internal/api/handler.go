package api

import (
    "net/http"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/service"
)

type Handler struct {
    RecipeService *service.RecipeService
}

func NewHandler(recipeService *service.RecipeService) *Handler {
    return &Handler{
        RecipeService: recipeService,
    }
}

func (h *Handler) FindRecipes(c *gin.Context) {
    var params model.SearchRequest
    if err := c.ShouldBindJSON(&params); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tidak valid"})
        return
    }

    fmt.Printf("Menerima parameter: %+v\n", params)

    directRecipes := getDirectRecipesWithFullPath(h.RecipeService.Graph, h.RecipeService.TierMap, params.TargetElement)
    if len(directRecipes) > 0 {
        if !params.MultipleRecipe {
            treeData := h.RecipeService.GenerateTreeData(params.TargetElement, directRecipes[:1])
            result := model.SearchResponse{
                Target: params.TargetElement,
                Algorithm: params.Algorithm,
                Time: 0.0,
                VisitedNodes: 1,
                Recipes: directRecipes[:1],
                TreeData: treeData,
            }
            c.JSON(http.StatusOK, result)
            return
        } else if len(directRecipes) >= params.MaxRecipes {
            treeData := h.RecipeService.GenerateTreeData(params.TargetElement, directRecipes[:params.MaxRecipes])
            result := model.SearchResponse{
                Target: params.TargetElement,
                Algorithm: params.Algorithm,
                Time: 0.0,
                VisitedNodes: len(directRecipes),
                Recipes: directRecipes[:params.MaxRecipes],
                TreeData: treeData,
            }
            c.JSON(http.StatusOK, result)
            return
        }
    }

    var result model.SearchResponse
    if params.MaxRecipes > 1 && params.MultipleRecipe {
        result = h.RecipeService.FindMultipleRecipes(params)
    } else {
        result = h.RecipeService.FindRecipes(params)
    }

    c.JSON(http.StatusOK, result)
}

func (h *Handler) GetAllElements(c *gin.Context) {
    elements := h.RecipeService.GetAllElements()
    c.JSON(http.StatusOK, gin.H{"elements": elements})
}

func getDirectRecipesWithFullPath(graph model.Graph, tierMap model.TierMap, targetElement string) []model.Recipe {
    var directRecipes []model.Recipe
    
    foundRecipes := make(map[string]bool)
    
    baseElements := map[string]bool{
        "Air": true,
        "Earth": true,
        "Fire": true,
        "Water": true,
        "Time": true,
    }
    
    targetTier, exists := tierMap[targetElement]
    if !exists {
        return directRecipes
    }
    
    if combinations, exists := graph[targetElement]; exists {
        for i, combo := range combinations {
            if len(combo) == 2 {
                isValidTier := true
                for _, ingredient := range combo {
                    ingredientTier, exists := tierMap[ingredient]
                    if !exists || ingredientTier >= targetTier {
                        isValidTier = false
                        break
                    }
                }
                
                if !isValidTier {
                    continue 
                }
                
                key := combo[0] + "," + combo[1]
                
                if !foundRecipes[key] {
                    recipe := model.Recipe{
                        ID: i,
                        Nodes: []model.RecipeNode{
                            {ID: fmt.Sprintf("target-%d", i), Label: targetElement, Level: 0},
                            {ID: fmt.Sprintf("ing-%d-0", i), Label: combo[0], Level: 1},
                            {ID: fmt.Sprintf("ing-%d-1", i), Label: combo[1], Level: 1},
                        },
                        Links: []model.RecipeLink{
                            {Source: fmt.Sprintf("target-%d", i), Target: fmt.Sprintf("ing-%d-0", i)},
                            {Source: fmt.Sprintf("target-%d", i), Target: fmt.Sprintf("ing-%d-1", i)},
                        },
                    }
                    
                    for j, ingredient := range combo {
                        if !baseElements[ingredient] {
                            expandIngredient(graph, tierMap, &recipe, ingredient, fmt.Sprintf("ing-%d-%d", i, j), 2, i, baseElements)
                        }
                    }
                    
                    directRecipes = append(directRecipes, recipe)
                    foundRecipes[key] = true
                }
            }
        }
    }
    
    return directRecipes
}

func expandIngredient(graph model.Graph, tierMap model.TierMap, recipe *model.Recipe, ingredient string, parentNodeID string, level int, recipeIndex int, baseElements map[string]bool) {
    // Dapatkan kombinasi untuk ingredient ini
    combinations, exists := graph[ingredient]
    if !exists || len(combinations) == 0 {
        return
    }
    
    ingredientTier, exists := tierMap[ingredient]
    if !exists {
        return
    }
    
    var validCombination []string
    for _, combination := range combinations {
        if len(combination) != 2 {
            continue
        }
        
        valid := true
        for _, ing := range combination {
            ingTier, exists := tierMap[ing]
            if !exists || ingTier >= ingredientTier {
                valid = false
                break
            }
        }
        
        if valid {
            validCombination = combination
            break
        }
    }
    
    if len(validCombination) != 2 {
        return
    }
    
    for i, subIngredient := range validCombination {
        nodeID := fmt.Sprintf("ingredient-%d-%d-%d", recipeIndex, level, i)
        
        node := model.RecipeNode{
            ID: nodeID,
            Label: subIngredient,
            Level: level,
        }
        recipe.Nodes = append(recipe.Nodes, node)
        
        link := model.RecipeLink{
            Source: parentNodeID,
            Target: nodeID,
        }
        recipe.Links = append(recipe.Links, link)
        
        if !baseElements[subIngredient] {
            expandIngredient(graph, tierMap, recipe, subIngredient, nodeID, level+1, recipeIndex, baseElements)
        }
    }
}