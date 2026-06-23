package problem

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/muhammad-iqbal/leetcode-backend/internal/config"
)

// GetProblems mengembalikan daftar semua soal
func GetProblems(w http.ResponseWriter, r *http.Request) {
	var problems []Problem
	// Mengambil semua soal (Omit: Jangan kirim test_cases yang berat ke halaman daftar)
	if err := config.DB.Omit("TestCases").Find(&problems).Error; err != nil {
		http.Error(w, "Failed to fetch problems", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(problems)
}

// GetProblemByID mengembalikan detail soal beserta contoh Test Case publik
func GetProblemByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	var problem Problem
	// Preload: Ambil relasi TestCases, TAPI HANYA yang tidak disembunyikan (is_hidden = false)
	if err := config.DB.Preload("TestCases", "is_hidden = ?", false).First(&problem, id).Error; err != nil {
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(problem)
}

// CreateProblem menyimpan soal baru beserta Test Case-nya ke database
func CreateProblem(w http.ResponseWriter, r *http.Request) {
	var problem Problem
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	// Simpan ke MySQL (GORM akan otomatis menyimpan relasi TestCases)
	if err := config.DB.Create(&problem).Error; err != nil {
		http.Error(w, "Gagal menyimpan soal ke database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(problem)
}
