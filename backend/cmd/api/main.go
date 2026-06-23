package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/muhammad-iqbal/leetcode-backend/internal/config"
	"github.com/muhammad-iqbal/leetcode-backend/internal/modules/problem"
	"github.com/muhammad-iqbal/leetcode-backend/internal/modules/submission"
	"github.com/muhammad-iqbal/leetcode-backend/internal/modules/user"
)

func main() {
	fmt.Println("LeetCode Clone API Server is starting...")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Inisialisasi Database
	config.InitMySQL()
	config.InitRedis()

	// Menjalankan GORM AutoMigrate (Membuat Tabel di MySQL secara otomatis)
	if config.DB != nil {
		log.Println("Running Database AutoMigrations...")
		_ = config.DB.AutoMigrate(
			&user.User{},
			&problem.Problem{},
			&problem.TestCase{},
			&submission.Submission{},
		)
	}

	// Inisialisasi Chi Router
	r := chi.NewRouter()
	
	// Konfigurasi CORS Berbasis Environment (Lokal vs Prod)
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, 
	}))

	// Middleware Standar bawaan Chi
	r.Use(chimiddleware.Logger)    // Mencatat semua log HTTP ke terminal
	r.Use(chimiddleware.Recoverer) // Mencegah API crash jika terjadi kepanikan (Panic)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK - API is running"))
	})

	// Endpoint khusus Integration Test
	r.Get("/health/redis", func(w http.ResponseWriter, r *http.Request) {
		if config.RedisClient == nil {
			http.Error(w, "Redis is not connected", http.StatusInternalServerError)
			return
		}
		
		// Test Write
		err := config.RedisClient.Set(config.Ctx, "integration_test_key", "integration_test_value", 5*time.Second).Err()
		if err != nil {
			http.Error(w, "Failed to write to Redis: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Test Read
		val, err := config.RedisClient.Get(config.Ctx, "integration_test_key").Result()
		if err != nil || val != "integration_test_value" {
			http.Error(w, "Failed to read from Redis", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK - Redis Read/Write Successful"))
	})

	// Mendaftarkan Grup Rute
	r.Mount("/api/auth", user.Routes())
	r.Mount("/api/problems", problem.Routes())
	r.Mount("/api/submissions", submission.Routes())

	log.Println("Server listening on port 8080")
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
