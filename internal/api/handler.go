package api

import (
    "net/http"

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

    var result model.SearchResult
    if params.MaxRecipes > 1 && !params.FindShortest {
        result = h.RecipeService.FindMultipleRecipes(params)
    } else {
        result = h.RecipeService.FindRecipes(params)
    }

    c.JSON(http.StatusOK, result)
}

func (h *Handler) GetAllElements(c *gin.Context) {
    elements := h.RecipeService.GetAllElements()
    c.JSON(http.StatusOK, model.ElementListResponse{Elements: elements})
}