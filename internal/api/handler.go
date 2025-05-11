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
    var params model.SearchParams
    if err := c.ShouldBindJSON(&params); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tidak valid"})
        return
    }

    fmt.Printf("Menerima parameter: %+v\n", params)

    directRecipes := getDirectRecipes(h.RecipeService.Graph, params.TargetElement)
    if len(directRecipes) > 0 {
        if params.FindShortest {
            result := model.SearchResult{
                TargetElement: params.TargetElement,
                Recipes:       directRecipes[:1],
                VisitedNodes:  1,
                TimeElapsed:   0.0,
            }
            c.JSON(http.StatusOK, result)
            return
        } else if len(directRecipes) >= params.MaxRecipes {
            result := model.SearchResult{
                TargetElement: params.TargetElement,
                Recipes:       directRecipes[:params.MaxRecipes],
                VisitedNodes:  len(directRecipes),
                TimeElapsed:   0.0,
            }
            c.JSON(http.StatusOK, result)
            return
        }
    }

    var result model.SearchResult
    if params.MaxRecipes > 1 && !params.FindShortest {
        result = h.RecipeService.FindMultipleRecipes(params)
    } else {
        result = h.RecipeService.FindRecipes(params)
    }

    fmt.Printf("Hasil: %+v\n", result)
    fmt.Printf("Jumlah recipe: %d\n", len(result.Recipes))
    for i, recipe := range result.Recipes {
        fmt.Printf("Recipe %d: %v\n", i+1, recipe.Ingredients)
    }

    c.JSON(http.StatusOK, result)
}

func (h *Handler) GetAllElements(c *gin.Context) {
    elements := h.RecipeService.GetAllElements()
    c.JSON(http.StatusOK, model.ElementListResponse{Elements: elements})
}

func getDirectRecipes(graph model.Graph, targetElement string) []model.Recipe {
    var directRecipes []model.Recipe
    
    if combinations, exists := graph[targetElement]; exists {
        for _, combo := range combinations {
            if len(combo) == 2 {
                directRecipes = append(directRecipes, model.Recipe{
                    Ingredients: []string{combo[0], combo[1]},
                })
            }
        }
    }
    
    return directRecipes
}