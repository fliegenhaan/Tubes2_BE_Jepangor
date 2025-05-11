package service

import (
    "fmt"
    "os"
    "path/filepath"
    "testing"

    "github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func TestNewRecipeService(t *testing.T) {
    // sesuaikan path file untuk testing
    currentDir, _ := os.Getwd()
    parentDir := filepath.Dir(filepath.Dir(currentDir))
    dataFilePath := filepath.Join(parentDir, "data", "elements_with_tier.json")

    fmt.Println("Path file yang digunakan:", dataFilePath)

    service, err := NewRecipeService(dataFilePath)
    if err != nil {
        t.Fatalf("Gagal inisialisasi recipe service: %v", err)
    }

    // cek graph sudah terbentuk
    if len(service.Graph) == 0 {
        t.Errorf("Graph kosong setelah inisialisasi")
    }

    // cek tierMap sudah terbentuk
    if len(service.TierMap) == 0 {
        t.Errorf("TierMap kosong setelah inisialisasi")
    }

    // cek elemen dasar
    baseElements := []string{"Air", "Earth", "Fire", "Water"}
    for _, element := range baseElements {
        tier, exists := service.TierMap[element]
        if !exists {
            t.Errorf("Elemen dasar %s tidak ditemukan dalam TierMap", element)
            continue
        }

        if tier != 0 {
            t.Errorf("Ekspektasi tier 0 untuk elemen dasar %s, dapat %d", element, tier)
        }

        _, exists = service.Graph[element]
        if !exists {
            t.Errorf("Elemen dasar %s tidak ditemukan dalam Graph", element)
        }
    }

    // cek beberapa elemen non-dasar
    nonBaseElements := []string{"Brick", "Glass", "Mud", "Lava"}
    for _, element := range nonBaseElements {
        _, exists := service.TierMap[element]
        if !exists {
            t.Errorf("Elemen non-dasar %s tidak ditemukan dalam TierMap", element)
            continue
        }

        combinations, exists := service.Graph[element]
        if !exists {
            t.Errorf("Elemen non-dasar %s tidak ditemukan dalam Graph", element)
            continue
        }

        if len(combinations) == 0 {
            t.Errorf("Tidak ada kombinasi untuk elemen non-dasar %s", element)
        }
    }

    // output debug
    fmt.Println("Total jumlah elemen di Graph:", len(service.Graph))
    fmt.Println("Total jumlah elemen di TierMap:", len(service.TierMap))

    // tampilkan contoh kombinasi
    for _, element := range nonBaseElements {
        combinations, exists := service.Graph[element]
        if exists {
            fmt.Printf("Kombinasi untuk %s (tier %d):\n", element, service.TierMap[element])
            for i, combo := range combinations {
                fmt.Printf("  %d: %v\n", i+1, combo)
            }
        }
    }
}

func TestFindRecipes(t *testing.T) {
    // sesuaikan path file untuk testing
    currentDir, _ := os.Getwd()
    parentDir := filepath.Dir(filepath.Dir(currentDir))
    dataFilePath := filepath.Join(parentDir, "data", "elements_with_tier.json")

    service, err := NewRecipeService(dataFilePath)
    if err != nil {
        t.Fatalf("Gagal inisialisasi recipe service: %v", err)
    }

    // test case
    testCases := []struct {
        targetElement string
        algorithm     string
        findShortest  bool
        maxRecipes    int
        expectRecipes bool // apakah seharusnya ditemukan resep
    }{
        {"Mud", "bfs", true, 1, true},              // Tier 1, seharusnya ditemukan
        {"Glass", "dfs", true, 1, true},            // Tier 4, seharusnya ditemukan
        {"Brick", "bfs", false, 5, true},           // Tier 2, seharusnya beberapa resep ditemukan
        {"ElemTidakAda", "bfs", true, 1, false},    // Elemen tidak ada, seharusnya tidak ditemukan
    }

    for _, tc := range testCases {
        params := model.SearchParams{
            TargetElement: tc.targetElement,
            Algorithm:     tc.algorithm,
            FindShortest:  tc.findShortest,
            MaxRecipes:    tc.maxRecipes,
        }

        result := service.FindRecipes(params)

        fmt.Printf("\nHasil pencarian untuk %s dengan algoritma %s:\n", tc.targetElement, tc.algorithm)
        fmt.Printf("  Jumlah resep yang ditemukan: %d\n", len(result.Recipes))
        fmt.Printf("  Jumlah node yang dikunjungi: %d\n", result.VisitedNodes)
        fmt.Printf("  Waktu pencarian: %.6f detik\n", result.TimeElapsed)

        if tc.expectRecipes && len(result.Recipes) == 0 {
            t.Errorf("Ekspektasi menemukan resep untuk %s, tapi tidak ditemukan", tc.targetElement)
        } else if !tc.expectRecipes && len(result.Recipes) > 0 {
            t.Errorf("Ekspektasi tidak ada resep untuk %s, tapi ditemukan %d", tc.targetElement, len(result.Recipes))
        }

        // tampilkan detail resep
        for i, recipe := range result.Recipes {
            fmt.Printf("  Resep %d: %v\n", i+1, recipe.Ingredients)
        }
    }
}

func TestFindMultipleRecipes(t *testing.T) {
    // sesuaikan path file untuk testing
    currentDir, _ := os.Getwd()
    parentDir := filepath.Dir(filepath.Dir(currentDir))
    dataFilePath := filepath.Join(parentDir, "data", "elements_with_tier.json")

    service, err := NewRecipeService(dataFilePath)
    if err != nil {
        t.Fatalf("Gagal inisialisasi recipe service: %v", err)
    }

    // test case untuk multiple recipes
    params := model.SearchParams{
        TargetElement: "Glass", // elemen yang memiliki beberapa resep
        Algorithm:     "bfs",
        FindShortest:  false,
        MaxRecipes:    5,
    }

    result := service.FindMultipleRecipes(params)

    fmt.Printf("\nHasil pencarian multiple untuk %s:\n", params.TargetElement)
    fmt.Printf("  Jumlah resep yang ditemukan: %d (dari max %d)\n", len(result.Recipes), params.MaxRecipes)
    fmt.Printf("  Jumlah node yang dikunjungi: %d\n", result.VisitedNodes)
    fmt.Printf("  Waktu pencarian: %.6f detik\n", result.TimeElapsed)

    if len(result.Recipes) == 0 {
        t.Errorf("Ekspektasi menemukan multiple recipes untuk %s, tapi tidak ditemukan", params.TargetElement)
    }

    // tampilkan detail resep
    for i, recipe := range result.Recipes {
        fmt.Printf("  Resep %d: %v\n", i+1, recipe.Ingredients)
    }
}

func TestSwimmingPoolRecipes(t *testing.T) {
    // sesuaikan path file untuk testing
    currentDir, _ := os.Getwd()
    parentDir := filepath.Dir(filepath.Dir(currentDir))
    dataFilePath := filepath.Join(parentDir, "data", "elements_with_tier.json")

    service, err := NewRecipeService(dataFilePath)
    if err != nil {
        t.Fatalf("Gagal inisialisasi recipe service: %v", err)
    }

    // test untuk Swimming pool dengan BFS
    testBFS := model.SearchParams{
        TargetElement: "Swimming pool",
        Algorithm:     "bfs",
        FindShortest:  false,
        MaxRecipes:    5,
    }

    resultBFS := service.FindRecipes(testBFS)

    fmt.Printf("\nHasil pencarian untuk %s dengan algoritma BFS:\n", testBFS.TargetElement)
    fmt.Printf("  Jumlah resep yang ditemukan: %d\n", len(resultBFS.Recipes))
    fmt.Printf("  Jumlah node yang dikunjungi: %d\n", resultBFS.VisitedNodes)
    fmt.Printf("  Waktu pencarian: %.6f detik\n", resultBFS.TimeElapsed)

    // tampilkan detail resep
    for i, recipe := range resultBFS.Recipes {
        fmt.Printf("  Resep %d: %v\n", i+1, recipe.Ingredients)
    }

    // test untuk Swimming pool dengan DFS
    testDFS := model.SearchParams{
        TargetElement: "Swimming pool",
        Algorithm:     "dfs",
        FindShortest:  false,
        MaxRecipes:    5,
    }

    resultDFS := service.FindRecipes(testDFS)

    fmt.Printf("\nHasil pencarian untuk %s dengan algoritma DFS:\n", testDFS.TargetElement)
    fmt.Printf("  Jumlah resep yang ditemukan: %d\n", len(resultDFS.Recipes))
    fmt.Printf("  Jumlah node yang dikunjungi: %d\n", resultDFS.VisitedNodes)
    fmt.Printf("  Waktu pencarian: %.6f detik\n", resultDFS.TimeElapsed)

    // tampilkan detail resep
    for i, recipe := range resultDFS.Recipes {
        fmt.Printf("  Resep %d: %v\n", i+1, recipe.Ingredients)
    }

    // info tentang Swimming pool dari data
    fmt.Println("\nInformasi Swimming pool dari data:")
    combinations, exists := service.Graph["Swimming pool"]
    if exists {
        fmt.Printf("  Ditemukan %d kombinasi dalam graph\n", len(combinations))
        for i, combo := range combinations {
            fmt.Printf("  - Kombinasi %d: %v\n", i+1, combo)
        }
    } else {
        fmt.Println("  Swimming pool tidak ditemukan dalam graph!")
    }

    tier, exists := service.TierMap["Swimming pool"]
    if exists {
        fmt.Printf("  Tier Swimming pool: %d\n", tier)
    } else {
        fmt.Println("  Swimming pool tidak ditemukan dalam tierMap!")
    }

    if len(resultBFS.Recipes) == 0 {
        t.Errorf("BFS: Ekspektasi menemukan resep untuk Swimming pool, tapi tidak ditemukan")
    }

    if len(resultDFS.Recipes) == 0 {
        t.Errorf("DFS: Ekspektasi menemukan resep untuk Swimming pool, tapi tidak ditemukan")
    }
}