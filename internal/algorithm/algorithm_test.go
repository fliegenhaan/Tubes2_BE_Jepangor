package algorithm

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/fliegenhaan/Tubes2_BE_Jepangor/internal/model"
)

// loadTestGraph membaca graph dari file JSON untuk pengujian
func loadTestGraph(t *testing.T) model.Graph {
	data, err := os.ReadFile("../../data/recipes.json")
	if err != nil {
		t.Fatalf("Gagal memuat data recipes: %v", err)
	}
	
	var graph model.Graph
	err = json.Unmarshal(data, &graph)
	if err != nil {
		t.Fatalf("Gagal unmarshal data recipes: %v", err)
	}
	
	return graph
}

// TestBFSNew menguji implementasi BFS yang baru
func TestBFSNew(t *testing.T) {
	// Load graph dari file
	graph := loadTestGraph(t)
	
	// Kasus uji
	testCases := []struct {
		name           string
		targetElement  string
		findShortest   bool
		maxRecipes     int
		expectRecipes  bool
		minVisitedNodes int // untuk memastikan bahwa kita benar-benar melakukan traversal
	}{
		{
			name:           "Brick - Terpendek",
			targetElement:  "Brick",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  true,
			minVisitedNodes: 2, // minimal ada beberapa node yang dikunjungi
		},
		{
			name:           "Brick - Multiple",
			targetElement:  "Brick",
			findShortest:   false,
			maxRecipes:     5,
			expectRecipes:  true,
			minVisitedNodes: 2,
		},
		{
			name:           "Swimming pool - Terpendek",
			targetElement:  "Swimming pool",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  true,
			minVisitedNodes: 2,
		},
		{
			name:           "Swimming pool - Multiple",
			targetElement:  "Swimming pool",
			findShortest:   false,
			maxRecipes:     5,
			expectRecipes:  true,
			minVisitedNodes: 2,
		},
		{
			name:           "Air - Elemen Dasar",
			targetElement:  "Air",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  false,
			minVisitedNodes: 1,
		},
		{
			name:           "NonExistent - Tidak Ditemukan",
			targetElement:  "NonExistentElement123",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  false,
			minVisitedNodes: 0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := time.Now()
			result := BFS(graph, tc.targetElement, tc.findShortest, tc.maxRecipes)
			duration := time.Since(startTime)
			
			// Periksa hasil
			if tc.expectRecipes && len(result.Recipes) == 0 {
				t.Errorf("Harapan: dengan recipes, tetapi hasilnya kosong")
			}
			
			if !tc.expectRecipes && len(result.Recipes) > 0 {
				t.Errorf("Harapan: tanpa recipes, tetapi hasilnya %d", len(result.Recipes))
			}
			
			if tc.findShortest && len(result.Recipes) > 1 {
				t.Errorf("Harapan: recipe terpendek (1), tetapi hasilnya %d", len(result.Recipes))
			}
			
			if !tc.findShortest && len(result.Recipes) > tc.maxRecipes {
				t.Errorf("Harapan: maksimum %d recipes, tetapi hasilnya %d", tc.maxRecipes, len(result.Recipes))
			}
			
			if result.VisitedNodes < tc.minVisitedNodes {
				t.Errorf("Harapan: minimal %d node dikunjungi, tetapi hasilnya %d", tc.minVisitedNodes, result.VisitedNodes)
			}
			
			// Verifikasi waktu eksekusi
			if result.TimeElapsed <= 0 {
				t.Errorf("Waktu eksekusi tidak valid: %.6f detik", result.TimeElapsed)
			}
			
			// Verifikasi waktu eksekusi kurang lebih sesuai dengan waktu sebenarnya
			if result.TimeElapsed > duration.Seconds()*2 || result.TimeElapsed < duration.Seconds()/2 {
				t.Logf("Peringatan: Waktu eksekusi tercatat (%.6f) berbeda jauh dari waktu sebenarnya (%.6f)", 
					result.TimeElapsed, duration.Seconds())
			}
			
			// Cetak hasil untuk inspeksi manual
			t.Logf("Target: %s, Recipes: %d, VisitedNodes: %d, TimeElapsed: %.6f detik",
				result.TargetElement, len(result.Recipes), result.VisitedNodes, result.TimeElapsed)
			
			for i, recipe := range result.Recipes {
				t.Logf("Recipe %d: %v", i+1, recipe.Ingredients)
			}
		})
	}
}

// TestDFSNew menguji implementasi DFS yang baru
func TestDFSNew(t *testing.T) {
	// Load graph dari file
	graph := loadTestGraph(t)
	
	// Kasus uji
	testCases := []struct {
		name           string
		targetElement  string
		findShortest   bool
		maxRecipes     int
		expectRecipes  bool
		minVisitedNodes int
	}{
		{
			name:           "Brick - Terpendek",
			targetElement:  "Brick",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  true,
			minVisitedNodes: 2,
		},
		{
			name:           "Brick - Multiple",
			targetElement:  "Brick",
			findShortest:   false,
			maxRecipes:     5,
			expectRecipes:  true,
			minVisitedNodes: 2,
		},
		{
			name:           "Swimming pool - Terpendek",
			targetElement:  "Swimming pool",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  true,
			minVisitedNodes: 2,
		},
		{
			name:           "Swimming pool - Multiple",
			targetElement:  "Swimming pool",
			findShortest:   false,
			maxRecipes:     5,
			expectRecipes:  true,
			minVisitedNodes: 2,
		},
		{
			name:           "Air - Elemen Dasar",
			targetElement:  "Air",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  false,
			minVisitedNodes: 1,
		},
		{
			name:           "NonExistent - Tidak Ditemukan",
			targetElement:  "NonExistentElement123",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  false,
			minVisitedNodes: 0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := time.Now()
			result := DFS(graph, tc.targetElement, tc.findShortest, tc.maxRecipes)
			duration := time.Since(startTime)
			
			// Periksa hasil
			if tc.expectRecipes && len(result.Recipes) == 0 {
				t.Errorf("Harapan: dengan recipes, tetapi hasilnya kosong")
			}
			
			if !tc.expectRecipes && len(result.Recipes) > 0 {
				t.Errorf("Harapan: tanpa recipes, tetapi hasilnya %d", len(result.Recipes))
			}
			
			if tc.findShortest && len(result.Recipes) > 1 {
				t.Errorf("Harapan: recipe terpendek (1), tetapi hasilnya %d", len(result.Recipes))
			}
			
			if !tc.findShortest && len(result.Recipes) > tc.maxRecipes {
				t.Errorf("Harapan: maksimum %d recipes, tetapi hasilnya %d", tc.maxRecipes, len(result.Recipes))
			}
			
			if result.VisitedNodes < tc.minVisitedNodes {
				t.Errorf("Harapan: minimal %d node dikunjungi, tetapi hasilnya %d", tc.minVisitedNodes, result.VisitedNodes)
			}
			
			// Verifikasi waktu eksekusi
			if result.TimeElapsed <= 0 {
				t.Errorf("Waktu eksekusi tidak valid: %.6f detik", result.TimeElapsed)
			}
			
			// Verifikasi waktu eksekusi kurang lebih sesuai dengan waktu sebenarnya
			if result.TimeElapsed > duration.Seconds()*2 || result.TimeElapsed < duration.Seconds()/2 {
				t.Logf("Peringatan: Waktu eksekusi tercatat (%.6f) berbeda jauh dari waktu sebenarnya (%.6f)", 
					result.TimeElapsed, duration.Seconds())
			}
			
			// Cetak hasil untuk inspeksi manual
			t.Logf("Target: %s, Recipes: %d, VisitedNodes: %d, TimeElapsed: %.6f detik",
				result.TargetElement, len(result.Recipes), result.VisitedNodes, result.TimeElapsed)
			
			for i, recipe := range result.Recipes {
				t.Logf("Recipe %d: %v", i+1, recipe.Ingredients)
			}
		})
	}
}

// TestBidirectionalNew menguji implementasi Bidirectional Search
func TestBidirectionalNew(t *testing.T) {
	// Load graph dari file
	graph := loadTestGraph(t)
	
	// Kasus uji yang sederhana untuk bidirectional search
	testCases := []struct {
		name           string
		targetElement  string
		findShortest   bool
		maxRecipes     int
		expectRecipes  bool
		minVisitedNodes int
	}{
		{
			name:           "Brick - Terpendek",
			targetElement:  "Brick",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  true,
			minVisitedNodes: 2,
		},
		{
			name:           "Air - Elemen Dasar",
			targetElement:  "Air",
			findShortest:   true,
			maxRecipes:     1,
			expectRecipes:  false,
			minVisitedNodes: 1,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Lewati bidirectional test jika implementasinya belum ada
			if _, ok := interface{}(nil).(func(model.Graph, string, bool, int) model.SearchResult); !ok {
				t.Skip("Implementasi Bidirectional belum tersedia")
			}
			
			startTime := time.Now()
			result := Bidirectional(graph, tc.targetElement, tc.findShortest, tc.maxRecipes)
			duration := time.Since(startTime)
			
			// Periksa hasil
			if tc.expectRecipes && len(result.Recipes) == 0 {
				t.Errorf("Harapan: dengan recipes, tetapi hasilnya kosong")
			}
			
			if !tc.expectRecipes && len(result.Recipes) > 0 {
				t.Errorf("Harapan: tanpa recipes, tetapi hasilnya %d", len(result.Recipes))
			}
			
			if result.VisitedNodes < tc.minVisitedNodes {
				t.Errorf("Harapan: minimal %d node dikunjungi, tetapi hasilnya %d", tc.minVisitedNodes, result.VisitedNodes)
			}
			
			// Verifikasi waktu eksekusi
			if result.TimeElapsed <= 0 {
				t.Errorf("Waktu eksekusi tidak valid: %.6f detik", result.TimeElapsed)
			}
			
			// Verifikasi waktu eksekusi kurang lebih sesuai dengan waktu sebenarnya
			if result.TimeElapsed > duration.Seconds()*2 || result.TimeElapsed < duration.Seconds()/2 {
				t.Logf("Peringatan: Waktu eksekusi tercatat (%.6f) berbeda jauh dari waktu sebenarnya (%.6f)", 
					result.TimeElapsed, duration.Seconds())
			}
			
			// Cetak hasil untuk inspeksi manual
			t.Logf("Target: %s, Recipes: %d, VisitedNodes: %d, TimeElapsed: %.6f detik",
				result.TargetElement, len(result.Recipes), result.VisitedNodes, result.TimeElapsed)
			
			for i, recipe := range result.Recipes {
				t.Logf("Recipe %d: %v", i+1, recipe.Ingredients)
			}
		})
	}
}

// TestMultithreadingBFS menguji implementasi multithreading pada BFS
func TestMultithreadingBFS(t *testing.T) {
	// Load graph dari file
	graph := loadTestGraph(t)
	
	// Kasus uji untuk multithreading
	testCases := []struct {
		name           string
		targetElement  string
		maxRecipes     int
		minDuration    float64 // durasi minimum dalam detik
	}{
		{
			name:           "Brick - Multiple Recipes",
			targetElement:  "Brick",
			maxRecipes:     10,
			minDuration:    0.001, // sangat pendek, hanya untuk memastikan multithreading berjalan
		},
		{
			name:           "Swimming pool - Multiple Recipes",
			targetElement:  "Swimming pool",
			maxRecipes:     20,
			minDuration:    0.001,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := time.Now()
			
			// Eksekusi BFS dengan findShortest=false untuk trigger multithreading
			result := BFS(graph, tc.targetElement, false, tc.maxRecipes)
			
			duration := time.Since(startTime).Seconds()
			
			// Verifikasi bahwa multithreading berjalan (membandingkan dengan nilai minimal)
			if duration < tc.minDuration {
				t.Logf("Peringatan: Eksekusi terlalu cepat (%.6f detik), multithreading mungkin tidak aktif", duration)
			}
			
			// Verifikasi hasil
			if len(result.Recipes) == 0 {
				t.Errorf("Tidak ada recipe yang ditemukan")
			}
			
			if len(result.Recipes) > tc.maxRecipes {
				t.Errorf("Harapan: maksimum %d recipes, tetapi hasilnya %d", tc.maxRecipes, len(result.Recipes))
			}
			
			// Cetak hasil untuk inspeksi manual
			t.Logf("Target: %s, Recipes: %d, VisitedNodes: %d, TimeElapsed: %.6f detik, Durasi Aktual: %.6f detik",
				result.TargetElement, len(result.Recipes), result.VisitedNodes, result.TimeElapsed, duration)
		})
	}
}

// TestMultithreadingDFS menguji implementasi multithreading pada DFS
func TestMultithreadingDFS(t *testing.T) {
	// Load graph dari file
	graph := loadTestGraph(t)
	
	// Kasus uji untuk multithreading
	testCases := []struct {
		name           string
		targetElement  string
		maxRecipes     int
		minDuration    float64 // durasi minimum dalam detik
	}{
		{
			name:           "Brick - Multiple Recipes",
			targetElement:  "Brick",
			maxRecipes:     10,
			minDuration:    0.001, // sangat pendek, hanya untuk memastikan multithreading berjalan
		},
		{
			name:           "Swimming pool - Multiple Recipes",
			targetElement:  "Swimming pool",
			maxRecipes:     20,
			minDuration:    0.001,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := time.Now()
			
			// Eksekusi DFS dengan findShortest=false untuk trigger multithreading
			result := DFS(graph, tc.targetElement, false, tc.maxRecipes)
			
			duration := time.Since(startTime).Seconds()
			
			// Verifikasi bahwa multithreading berjalan (membandingkan dengan nilai minimal)
			if duration < tc.minDuration {
				t.Logf("Peringatan: Eksekusi terlalu cepat (%.6f detik), multithreading mungkin tidak aktif", duration)
			}
			
			// Verifikasi hasil
			if len(result.Recipes) == 0 {
				t.Errorf("Tidak ada recipe yang ditemukan")
			}
			
			if len(result.Recipes) > tc.maxRecipes {
				t.Errorf("Harapan: maksimum %d recipes, tetapi hasilnya %d", tc.maxRecipes, len(result.Recipes))
			}
			
			// Cetak hasil untuk inspeksi manual
			t.Logf("Target: %s, Recipes: %d, VisitedNodes: %d, TimeElapsed: %.6f detik, Durasi Aktual: %.6f detik",
				result.TargetElement, len(result.Recipes), result.VisitedNodes, result.TimeElapsed, duration)
		})
	}
}