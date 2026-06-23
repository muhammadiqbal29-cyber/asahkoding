package user

import (
	"encoding/json"
	"net/http"

	"github.com/muhammad-iqbal/leetcode-backend/internal/config"
	"github.com/muhammad-iqbal/leetcode-backend/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	// Pastikan Database tidak mati
	if config.DB == nil {
		http.Error(w, "Database is offline", http.StatusInternalServerError)
		return
	}

	// Enkripsi Password menggunakan bcrypt (Sangat Aman)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal mengenkripsi password", http.StatusInternalServerError)
		return
	}

	user := User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         RoleUser, // Pendaftar baru otomatis menjadi user biasa (bukan admin)
	}

	// Simpan ke MySQL via GORM
	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Gagal membuat user. Email atau Username mungkin sudah dipakai.", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil mendaftar!"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	if config.DB == nil {
		http.Error(w, "Database is offline", http.StatusInternalServerError)
		return
	}

	var user User
	// Cari user berdasarkan email
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Email atau Password salah", http.StatusUnauthorized)
		return
	}

	// Verifikasi kecocokan password asli dengan Hash acak di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Email atau Password salah", http.StatusUnauthorized)
		return
	}

	// Jika cocok, cetak tiket JWT
	token, err := middleware.GenerateToken(user.ID, string(user.Role), user.Username)
	if err != nil {
		http.Error(w, "Gagal membuat token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Login berhasil",
		"token":   token,
	})
}
