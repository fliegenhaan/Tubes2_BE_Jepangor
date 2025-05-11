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
    // path file data
    dataFilePath := "./data/elements_with_tier.json"
    if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
        log.Fatalf("File data recipe tidak ditemukan: %s", dataFilePath)
    }

    // inisialisasi service
    recipeService, err := service.NewRecipeService(dataFilePath)
    if err != nil {
        log.Fatalf("Gagal inisialisasi recipe service: %v", err)
    }

    // buat handler
    handler := api.NewHandler(recipeService)

    // setup router
    router := gin.Default()
    router.Use(api.CORSMiddleware())
    router.Use(api.LoggerMiddleware())

    // definisikan endpoint
    router.POST("/api/find-recipes", handler.FindRecipes)
    router.GET("/api/elements", handler.GetAllElements)

    // ambil port dari environment variable
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("Server berjalan pada port %s\n", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Gagal menjalankan server: %v", err)
    }
}