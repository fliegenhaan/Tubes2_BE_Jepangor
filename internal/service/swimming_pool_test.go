package service

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

func TestFindSwimmingPoolRecipes(t *testing.T) {
	// Ubah path file ke posisi relatif terhadap direktori test
	currentDir, _ := os.Getwd()
	parentDir := filepath.Dir(filepath.Dir(currentDir))
	dataFilePath := filepath.Join(parentDir, "data", "elements_with_tier.json")

	service, err := NewRecipeService(dataFilePath)
	if err != nil {
		t.Fatalf("Failed to initialize recipe service: %v", err)
	}

	// Test case untuk Swimming pool dengan BFS
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

	// Tampilkan detail resep jika ditemukan
	for i, recipe := range resultBFS.Recipes {
		fmt.Printf("  Resep %d: %v\n", i+1, recipe.Ingredients)
	}

	// Test case untuk Swimming pool dengan DFS
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

	// Tampilkan detail resep jika ditemukan
	for i, recipe := range resultDFS.Recipes {
		fmt.Printf("  Resep %d: %v\n", i+1, recipe.Ingredients)
	}

	// Debug: Tampilkan informasi tentang Swimming pool dari graph dan tier
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

	// Cek upstream dan downstream di graph
	fmt.Println("\nMencari dalam graph:")
	foundAsResult := false
	for element, combos := range service.Graph {
		for _, combo := range combos {
			if len(combo) == 2 && element == "Swimming pool" {
				fmt.Printf("  Swimming pool dapat dibuat dari: %s + %s\n", combo[0], combo[1])
				foundAsResult = true
			}
		}
	}
	if !foundAsResult {
		fmt.Println("  Swimming pool tidak ditemukan sebagai hasil kombinasi!")
	}

	// Verifikasi assertions
	if len(resultBFS.Recipes) == 0 {
		t.Errorf("BFS: Expected to find recipes for Swimming pool, but none found")
	}

	if len(resultDFS.Recipes) == 0 {
		t.Errorf("DFS: Expected to find recipes for Swimming pool, but none found")
	}
}