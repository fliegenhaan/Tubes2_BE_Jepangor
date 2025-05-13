# Tubes2_FE_Jepangor

# Little Alchemy 2 Recipe Finder

Little Alchemy 2 Recipe Finder adalah aplikasi berbasis web yang membantu pengguna menemukan recipe (resep) untuk menghasilkan elemen-elemen dalam permainan Little Alchemy 2. Aplikasi ini mengimplementasikan algoritma BFS (Breadth-First Search), DFS (Depth-First Search), dan Bidirectional Search untuk menemukan jalur pembentukan suatu elemen dari elemen-elemen dasar.

## Algoritma yang Diimplementasikan

### 1. Breadth-First Search (BFS)
BFS adalah algoritma pencarian graf yang mengeksplorasi semua node pada kedalaman yang sama sebelum melanjutkan ke node pada kedalaman berikutnya. Algoritma ini cocok untuk menemukan jalur terpendek dalam graf yang tidak berbobot.

### 2. Depth-First Search (DFS)
DFS adalah algoritma pencarian graf yang mengeksplorasi jalur secara mendalam terlebih dahulu sebelum melakukan backtracking. Algoritma ini efisien dalam ruang memori namun tidak menjamin menemukan jalur terpendek.

### 3. Bidirectional Search (Bonus)
Bidirectional Search adalah teknik pencarian yang menjalankan dua pencarian sekaligus: satu dari node awal dan satu dari node tujuan. Pencarian berhenti ketika kedua pencarian bertemu di tengah. Ini dapat mengurangi jumlah node yang perlu dieksplorasi.

## Fitur Utama

- Pencarian Recipe: Cari recipe untuk membuat elemen tertentu dari elemen dasar menggunakan algoritma BFS, DFS, dan Bidirectional Search.
- Visualisasi Tree: Tampilkan visualisasi tree untuk recipe yang ditemukan, menunjukkan kombinasi elemen yang diperlukan.
- Multiple Recipe: Temukan beberapa recipe berbeda untuk membuat elemen yang sama.

## Requirement Program dan Instalasi

### Frontend
- Node.js
- npm atau yarn

### Backend
- Go
- Gin framework

## Cara Menjalankan Program

### Frontend
1. Clone repository frontend
   ```
   git clone https://github.com/fliegenhaan/Tubes2_FE_Jepangor.git
   ```
2. Masuk ke direktori project
   ```
   cd jepangor
   ```
3. Install dependensi
   ```
   npm install
   ```
4. Jalankan aplikasi dalam mode development
   ```
   npm run dev
   ```

### Backend
1. Clone repository backend
   ```
   git clone https://github.com/fliegenhaan/Tubes2_BE_Jepangor.git
   ```
2. Jalankan aplikasi
   ```
   go run cmd/main.go
   ```

### Aplikasi juga dapat diakses secara online di:
```
jepangor.vercel.app
```

## Struktur Program

### Frontend
```
Tubes2_FE_Jepangor/
├── doc/
│   └── Jepangor.pdf
├── jepangor/
│   ├── public/
│   │   └── placeholder.png
│   ├── src/
│   │   ├── app/
│   │   │   ├── about/
│   │   │   │   └── page.tsx
│   │   │   ├── globals.css
│   │   │   ├── layout.tsx
│   │   │   └── page.tsx
│   │   ├── components/
│   │   │   ├── AlgorithmSelector.tsx
│   │   │   ├── ElementSelector.tsx
│   │   │   ├── Header.tsx
│   │   │   ├── NavBar.tsx
│   │   │   ├── RecipeTree.tsx
│   │   │   ├── ResultInfo.tsx
│   │   │   └── SearchForm.tsx
│   │   ├── hooks/
│   │   │   └── useRecipeSearch.ts
│   │   ├── types/
│   │   │   └── index.ts
│   │   └── utils/
│   │       └── api.ts
│   ├── .gitignore
│   ├── README.md
│   ├── eslint.config.mjs
│   ├── next.config.js
│   ├── package-lock.json
│   ├── package.json
│   ├── postcss.config.js
│   ├── postcss.config.mjs
│   ├── tailwind.config.js
│   └── tsconfig.json
└── README.md
```

### Backend
```
Tubes2_BE_Jepangor/
├── cmd/
│   └── main.go
├── data/
│   └── elements_with_tier.json
├── doc/
│   └── Jepangor.pdf
├── internal/
│   ├── algorithm/
│   │   ├── bfs.go
│   │   ├── bidirectional.go
│   │   └── dfs.go
│   ├── api/
│   │   ├── handler.go
│   │   └── middleware.go
│   ├── model/
│   │   └── recipe.go
│   └── service/
│       └── recipe_service.go
├── web-scraper-go/
│   ├── go.mod
│   └── go.sum
├── .gitignore
├── README.md
├── Dockerfile
├── go.mod
├── go.sum
└── vercel.json
```

## Author

### Kelompok Jepangor
- Muhammad Raihaan Perdana (13523124) - FullStack & Algorithm Developer
- Ardell Aghna Mahendra (13523151) - UI/UX & Frontend Developer
- Farrel Jabaar Altafataza (10122057) - Scraper & Algorithm Developer

## Bonus
- Video: https://youtu.be/pdd0t1Yrkwo
- pencarian dengan strategi Bidirectional
- Aplikasi di-deploy dan dapat diakses melalui internet
- Aplikasi di-containerize dengan Docker
