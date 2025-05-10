package main

import (
    "fmt"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/api"
    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/service"
)

func main() {
    dataFilePath := "./data/recipes.json"
    if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
        log.Fatalf("Recipe data file not found: %s", dataFilePath)
    }

    recipeService, err := service.NewRecipeService(dataFilePath)
    if err != nil {
        log.Fatalf("Failed to initialize recipe service: %v", err)
    }

    handler := api.NewHandler(recipeService)

    router := gin.Default()
    router.Use(api.CORSMiddleware())
    router.Use(api.LoggerMiddleware())

    router.POST("/api/find-recipes", handler.FindRecipes)
    router.GET("/api/elements", handler.GetAllElements)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("Server running on port %s\n", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}